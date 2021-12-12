package index

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	indexoptimizer "github.com/actiontech/sqle/sqle/pkg/optimizer/index"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	defaultCalculateCardinalityMaxRow = 1000000
	defaultCompositeIndexMaxColumn    = 3
)

type Optimizer struct {
	*session.Context

	l *logrus.Entry

	// tables key is table name, use to match in execution plan.
	tables             map[string]*tableInSelect
	tableNameExtractor util.TableNameExtractor

	// optimizer options:
	calculateCardinalityMaxRow int
	compositeIndexMaxColumn    int
	createIndexStatement       func(string, ...string) string
}

func NewOptimizer(log *logrus.Entry, ctx *session.Context, opts ...optimizerOption) *Optimizer {
	log = log.WithField("optimizer", "index")

	optimizer := &Optimizer{
		Context:                    ctx,
		l:                          log,
		tables:                     make(map[string]*tableInSelect),
		createIndexStatement:       defaultCreateIndexStatement,
		compositeIndexMaxColumn:    defaultCompositeIndexMaxColumn,
		calculateCardinalityMaxRow: defaultCalculateCardinalityMaxRow,
	}

	for _, opt := range opts {
		opt.apply(optimizer)
	}

	return optimizer
}

type OptimizeResult struct {
	TableName      string
	IndexedColumns []string

	Reason string
}

// tableInSelect store the information of a table in select statement for later optimize.
// 1. when we find a table in single table select statement, we will store the select statement.
// 2. when we find a table in join statement, we will store the join on condition.
type tableInSelect struct {
	joinOnColumn   string
	singleTableSel *ast.SelectStmt
}

// Optimize give index advice for the select statement.
func (o *Optimizer) Optimize(ctx context.Context, selectStmt *ast.SelectStmt) ([]*OptimizeResult, error) {
	// select 1; ...
	if selectStmt.From == nil {
		return nil, nil
	}

	o.parseSelectStmt(selectStmt)

	restoredSQL, err := restoreSelectStmt(selectStmt)
	if err != nil {
		return nil, err
	}

	executionPlan, err := o.GetExecutionPlan(restoredSQL)
	if err != nil {
		return nil, errors.Wrap(err, "get execution plan when optimize")
	}

	executionPlan = removeDrivingTable(executionPlan)

	var needOptimizedTables []string
	for _, record := range executionPlan {
		if o.needOptimize(record) {
			needOptimizedTables = append(needOptimizedTables, record.Table)
		}
	}

	if len(needOptimizedTables) == 0 {
		return nil, nil
	}

	o.l.Infof("need optimize tables: %v", needOptimizedTables)

	var results []*OptimizeResult
	for _, tbl := range needOptimizedTables {
		table, ok := o.tables[tbl]
		if !ok {
			// given SQL: select * from t1 join t2, there is no join on condition,
			continue
		}

		var result *OptimizeResult
		if table.joinOnColumn == "" {
			result, err = o.optimizeSingleTable(ctx, tbl, table.singleTableSel)
			if err != nil {
				return nil, errors.Wrapf(err, "optimize single table %s", tbl)
			}
		} else {
			result = o.optimizeJoinTable(tbl)
		}
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

// SelectStmt:
//   1. single select on single table
//   2. single select on multiple tables, such join
//   3. multi select on multiple tables, such subqueries
func (o *Optimizer) parseSelectStmt(ss *ast.SelectStmt) {
	tne := util.TableNameExtractor{TableNames: map[string]*ast.TableName{}}
	ss.Accept(&tne)
	o.tableNameExtractor = tne

	visitor := util.SelectStmtExtractor{}
	ss.Accept(&visitor)

	for _, ss := range visitor.SelectStmts {
		if ss.From == nil {
			continue
		}

		left := ss.From.TableRefs.Left
		right := ss.From.TableRefs.Right

		if right == nil {
			leftTable, ok := left.(*ast.TableSource)
			if !ok {
				continue
			}

			if leftTable.AsName.L != "" {
				o.tables[leftTable.AsName.L] = &tableInSelect{singleTableSel: ss}
			} else {
				o.tables[leftTable.Source.(*ast.TableName).Name.L] = &tableInSelect{singleTableSel: ss}
			}
		} else {
			if ss.From.TableRefs.On != nil {
				boe, ok := ss.From.TableRefs.On.Expr.(*ast.BinaryOperationExpr)
				if !ok {
					continue
				}

				leftCNE := boe.L.(*ast.ColumnNameExpr)
				rightCNE := boe.R.(*ast.ColumnNameExpr)

				o.tables[leftCNE.Name.Table.L] = &tableInSelect{joinOnColumn: leftCNE.Name.Name.L}
				o.tables[rightCNE.Name.Table.L] = &tableInSelect{joinOnColumn: rightCNE.Name.Name.L}
			} else if ss.From.TableRefs.Using != nil {
				leftTableName := left.(*ast.TableSource).Source.(*ast.TableName).Name.L
				rightTableName := right.(*ast.TableSource).Source.(*ast.TableName).Name.L
				for _, col := range ss.From.TableRefs.Using {
					o.tables[leftTableName] = &tableInSelect{joinOnColumn: col.Name.L}
					o.tables[rightTableName] = &tableInSelect{joinOnColumn: col.Name.L}
				}
			}
		}
	}
}

func (o *Optimizer) optimizeSingleTable(ctx context.Context, tbl string, ss *ast.SelectStmt) (*OptimizeResult, error) {
	var (
		err            error
		optimizeResult *OptimizeResult
	)

	optimizeResult, err = o.doSpecifiedOptimization(ctx, ss)
	if err != nil {
		return nil, err
	}

	if optimizeResult == nil {
		optimizeResult, err = o.doGeneralOptimization(ctx, ss)
		if err != nil {
			return nil, err
		}
	}

	if optimizeResult == nil {
		return nil, nil
	}

	if len(optimizeResult.IndexedColumns) > o.compositeIndexMaxColumn {
		optimizeResult.IndexedColumns = optimizeResult.IndexedColumns[:o.compositeIndexMaxColumn]
	}

	needIndex, err := o.needIndex(optimizeResult.TableName, optimizeResult.IndexedColumns...)
	if err != nil {
		return nil, err
	}

	if !needIndex {
		return nil, nil
	}

	o.l.Infof("table:%s, indexed columns:%v, reason:%s", optimizeResult.TableName, optimizeResult.IndexedColumns, optimizeResult.Reason)

	if len(optimizeResult.IndexedColumns) > 1 {
		rowCount, err := o.GetTableRowCount(extractTableNameFromAST(ss, tbl))
		if err != nil {
			return nil, errors.Wrap(err, "get table row count when optimize")
		}
		if rowCount < o.calculateCardinalityMaxRow {
			optimizeResult.IndexedColumns, err = o.sortColumnsByCardinality(tbl, optimizeResult.IndexedColumns)
			if err != nil {
				return nil, err
			}
		}
	}

	return optimizeResult, nil
}

func (o *Optimizer) optimizeJoinTable(tbl string) *OptimizeResult {
	return &OptimizeResult{
		TableName:      tbl,
		IndexedColumns: []string{o.tables[tbl].joinOnColumn},
		Reason:         fmt.Sprintf("字段 %s 为被驱动表 %s 上的关联字段", o.tables[tbl].joinOnColumn, tbl),
	}
}

// doSpecifiedOptimization optimize single table select.
func (o *Optimizer) doSpecifiedOptimization(ctx context.Context, selectStmt *ast.SelectStmt) (*OptimizeResult, error) {
	if where := selectStmt.Where; where != nil {
		if boe, ok := where.(*ast.BinaryOperationExpr); ok {
			// check function in select stmt
			if fce, ok := boe.L.(*ast.FuncCallExpr); ok {
				result, err := o.optimizeOnFunctionCallExpression(getTableNameFromSingleSelect(selectStmt), fce)
				if err != nil {
					return nil, err
				}
				if result != nil {
					return result, nil
				}
			}
		}

		// check where like 'mike%'
		if ple, ok := where.(*ast.PatternLikeExpr); ok {
			if cne, ok := ple.Expr.(*ast.ColumnNameExpr); ok {
				if ve, ok := ple.Pattern.(*driver.ValueExpr); ok {
					datum := ve.Datum.GetString()
					if !strings.HasPrefix(datum, "%") &&
						!strings.HasPrefix(datum, "_") {
						return &OptimizeResult{
							TableName:      getTableNameFromSingleSelect(selectStmt),
							IndexedColumns: []string{cne.Name.Name.L},
							Reason:         "为前缀模式匹配添加前缀索引",
						}, nil
					}
				}
			}
		}
	}

	return nil, nil
}

func (o *Optimizer) optimizeOnFunctionCallExpression(tbl string, fce *ast.FuncCallExpr) (*OptimizeResult, error) {
	var cols []string
	for _, arg := range fce.Args {
		if cne, ok := arg.(*ast.ColumnNameExpr); ok {
			cols = append(cols, cne.Name.Name.L)
		}
	}
	if len(cols) == 0 {
		return nil, nil
	}

	var buf strings.Builder
	if err := fce.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)); err != nil {
		return nil, errors.Wrap(err, "restore func call expr when do specified optimization")
	}

	versionWithFlavor, err := o.GetSystemVariable("version")
	if err != nil {
		return nil, errors.Wrap(err, "get version when do specified optimization")
	}

	curVersion, err := semver.NewVersion(versionWithFlavor)
	if err != nil {
		return nil, errors.Wrap(err, "parse version when do specified optimization")
	}
	if curVersion.LessThan(semver.MustParse("5.7.0")) {
		return nil, nil
	}
	if curVersion.LessThan(semver.MustParse("8.0.13")) {
		return &OptimizeResult{
			TableName:      tbl,
			IndexedColumns: []string{buf.String()},
			Reason:         "MySQL5.7以上版本需要在虚拟列上创建索引",
		}, nil
	}

	return &OptimizeResult{
		TableName:      tbl,
		IndexedColumns: []string{buf.String()},
		Reason:         "MySQL8.0.13以上版本支持直接创建函数索引",
	}, nil
}

// doGeneralOptimization optimize single table select.
func (o *Optimizer) doGeneralOptimization(ctx context.Context, selectStmt *ast.SelectStmt) (*OptimizeResult, error) {
	generalOptimizer := indexoptimizer.NewOptimizer()

	restoredSQL, err := restoreSelectStmt(selectStmt)
	if err != nil {
		return nil, err
	}

	sa, err := newSelectAST(restoredSQL)
	if err != nil {
		return nil, err
	}

	indexedColumns, err := generalOptimizer.Optimize(sa)
	if err != nil {
		return nil, err
	}

	if len(indexedColumns) == 0 {
		return nil, nil
	}

	o.l.Infof("general optimize result: %v(index columns)", indexedColumns)

	return &OptimizeResult{
		TableName:      getTableNameFromSingleSelect(selectStmt),
		IndexedColumns: indexedColumns,
		Reason:         "三星索引建议",
	}, nil
}

type cardinality struct {
	columnName  string
	cardinality int
}

type cardinalities []cardinality

func (c cardinalities) Len() int {
	return len(c)
}

func (c cardinalities) Less(i, j int) bool {
	return c[i].cardinality > c[j].cardinality
}

func (c cardinalities) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (o *Optimizer) sortColumnsByCardinality(tbl string, indexedColumns []string) (sortedColumns []string, err error) {
	var tn *ast.TableName
	for tableName, currTN := range o.tableNameExtractor.TableNames {
		if tableName == tbl {
			tn = currTN
			break
		}
	}
	if tn == nil {
		return nil, errors.Errorf("table %s not found when sort columns by cardinality", tbl)
	}

	var cardinalities cardinalities
	for _, column := range indexedColumns {
		c, err := o.GetColumnCardinality(tn, column)
		if err != nil {
			return nil, errors.Wrap(err, "get column cardinality")
		}
		cardinalities = append(cardinalities, cardinality{
			columnName:  column,
			cardinality: c,
		})
	}

	o.l.Debugf("table %s column cardinalities(before sort): %+v", tbl, cardinalities)
	sort.Sort(cardinalities)
	o.l.Debugf("table %s column cardinalities(after sort): %+v", tbl, cardinalities)

	for _, c := range cardinalities {
		sortedColumns = append(sortedColumns, c.columnName)
	}
	return sortedColumns, nil
}

// needOptimize check table need optimize index of table or not.
//
// Optimize means that:
// 1. When SQL do not use index, we can create index for the select statement.
// 2. When SQL use index, but the index is not suitable, we should optimize the index.
//
// We do it by check MySQL execution plan's access_type field.
// ref: https://dev.mysql.com/doc/refman/5.7/en/explain-output.html#explain-join-types
func (o *Optimizer) needOptimize(record *executor.ExplainRecord) bool {

	// Full table scan: select * from t1 where common_column = 'a'
	// This SQL will scan all rows of table t1.
	if record.Type == executor.ExplainRecordAccessTypeAll {
		return true
	}

	// Index-only scan: select key_part2 from t1 where key_part3 = 'a'
	// This SQL will scan all rows of index idx_composite. It's a little better than previous case.
	if record.Type == executor.ExplainRecordAccessTypeIndex {
		return true
	}

	return false
}

// needIndex check need add index on tbl.columns or not.
func (o *Optimizer) needIndex(tbl string, columns ...string) (bool, error) {
	table, ok := o.tables[tbl]
	if !ok {
		return false, fmt.Errorf("table %s not found when check index", tbl)
	}

	if table.singleTableSel == nil {
		return false, fmt.Errorf("table %s do not have select statement when check index", tbl)
	}

	cts, exist, err := o.GetCreateTableStmt(extractTableNameFromAST(table.singleTableSel, tbl))
	if err != nil {
		return false, errors.Wrap(err, "get create table statement when check index")
	}
	if !exist {
		return false, fmt.Errorf("table %s not found on session context when check index", tbl)
	}

	for _, index := range util.ExtractIndexFromCreateTableStmt(cts) {
		if reflect.DeepEqual(index, columns) {
			return false, nil
		}
		if strings.HasPrefix(strings.Join(index, ","), strings.Join(columns, ",")) {
			return false, nil
		}
	}
	return true, nil
}

type optimizerOption func(*Optimizer)

func (oo optimizerOption) apply(o *Optimizer) {
	oo(o)
}

func WithCalculateCardinalityMaxRow(row int) optimizerOption {
	return func(o *Optimizer) {
		o.calculateCardinalityMaxRow = row
	}
}

func WithCompositeIndexMaxColumn(column int) optimizerOption {
	return func(o *Optimizer) {
		o.compositeIndexMaxColumn = column
	}
}

func WithCreateIndexStatement(f func(tableName string, columns ...string) string) optimizerOption {
	return func(o *Optimizer) {
		o.createIndexStatement = f
	}
}

func defaultCreateIndexStatement(tableName string, columns ...string) string {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, strings.Join(columns, "_"))

	return fmt.Sprintf("CREATE INDEX %s ON %s (%s)",
		indexName,
		tableName,
		strings.Join(columns, ", "))
}

func restoreSelectStmt(ss *ast.SelectStmt) (string, error) {
	var buf strings.Builder
	rc := format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)

	if err := ss.Restore(rc); err != nil {
		return "", errors.Wrap(err, "restore select statement")
	}

	return buf.String(), nil
}

func extractTableNameFromAST(ss *ast.SelectStmt, tbl string) *ast.TableName {
	v := util.TableNameExtractor{TableNames: make(map[string]*ast.TableName)}
	ss.Accept(&v)

	for _, t := range v.TableNames {
		if t.Name.L == tbl {
			return t
		}
	}
	return nil
}

func getTableNameFromSingleSelect(ss *ast.SelectStmt) string {
	if ss.From.TableRefs.Left == nil {
		return ""
	}
	return ss.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.L
}

// removeDrivingTable remove driving table from execution plan.
//
// Index is not silver bullet, we only give advice on driven table.
// Such as : select * from t1, t2 where t1.id = t2.id;
// There are two records in execution plan, the first one is driving table, the second one is driven table.
func removeDrivingTable(records []*executor.ExplainRecord) []*executor.ExplainRecord {
	var result []*executor.ExplainRecord

	if len(records) == 0 || len(records) == 1 {
		return records
	}

	i, j := 0, 1
	for j < len(records) {
		if records[i].Id == records[j].Id {
			result = append(result, records[j])
		} else {
			i = j
		}
		j++
	}

	return result
}

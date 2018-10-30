package inspector

import (
	"errors"
	"github.com/pingcap/tidb/ast"
	_model "github.com/pingcap/tidb/model"
	"sqle/executor"
	"sqle/model"
)

func CreateRollbackSql(task *model.Task, sql string) (string, error) {
	return "", nil
	conn, err := executor.OpenDbWithTask(task)
	if err != nil {
		return "", err
	}
	switch task.Instance.DbType {
	case model.DB_TYPE_MYSQL:
		return createRollbackSql(conn, sql)
	default:
		return "", errors.New("db type is invalid")
	}
}

// createRollbackSql create rollback sql for input sql; this sql is single sql.
func createRollbackSql(conn *executor.Conn, query string) (string, error) {
	stmts, err := parseSql(model.DB_TYPE_MYSQL, query)
	if err != nil {
		return "", err
	}
	stmt := stmts[0]
	switch n := stmt.(type) {
	case *ast.AlterTableStmt:
		tableName := getTableName(n.Table)
		createQuery, err := conn.ShowCreateTable(tableName)
		if err != nil {
			return "", err
		}
		t, err := parseOneSql(model.DB_TYPE_MYSQL, createQuery)
		if err != nil {
			return "", err
		}
		createStmt, ok := t.(*ast.CreateTableStmt)
		if !ok {
			return "", errors.New("")
		}
		return alterTableRollbackSql(createStmt, n)
	default:
		return "", nil
	}
}

func alterTableRollbackSql(t1 *ast.CreateTableStmt, t2 *ast.AlterTableStmt) (string, error) {
	table := t1.Table
	table.Schema = t2.Table.Schema

	t := &ast.AlterTableStmt{
		Table: table,
		Specs: []*ast.AlterTableSpec{},
	}
	// table name
	specs := t2.Specs
	for _, spec := range specs {
		switch spec.Tp {
		case ast.AlterTableRenameTable:
			t.Table = spec.NewTable
			t.Table.Schema = t2.Table.Schema
			newTable := t1.Table
			newTable.Schema = _model.CIStr{}
			t.Specs = append(t.Specs, &ast.AlterTableSpec{
				Tp:       ast.AlterTableRenameTable,
				NewTable: t1.Table,
			})
		case ast.AlterTableAddColumns:
			for _, col := range spec.NewColumns {
				t.Specs = append(t.Specs, &ast.AlterTableSpec{
					Tp:            ast.AlterTableDropColumn,
					OldColumnName: col.Name,
				})
			}
		case ast.AlterTableDropColumn:
			colName := spec.OldColumnName.String()
			for _, col := range t1.Cols {
				if col.Name.String() == colName {
					t.Specs = append(t.Specs, &ast.AlterTableSpec{
						Tp:         ast.AlterTableAddColumns,
						NewColumns: []*ast.ColumnDef{col},
					})
				}
			}
		}
	}
	return alterTableStmtFormat(t), nil
}

package inspector

import (
	"fmt"
	"strconv"
	"strings"

	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

// inspector rule code
const (
	DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST       = "ddl_check_table_without_if_not_exists"
	DDL_CHECK_OBJECT_NAME_LENGTH               = "ddl_check_object_name_length"
	DDL_CHECK_OBJECT_NAME_USING_KEYWORD        = "ddl_check_object_name_using_keyword"
	DDL_CHECK_PK_NOT_EXIST                     = "ddl_check_pk_not_exist"
	DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED       = "ddl_check_pk_without_bigint_unsigned"
	DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT        = "ddl_check_pk_without_auto_increment"
	DDL_CHECK_COLUMN_VARCHAR_MAX               = "ddl_check_column_varchar_max"
	DDL_CHECK_COLUMN_CHAR_LENGTH               = "ddl_check_column_char_length"
	DDL_DISABLE_FK                             = "ddl_disable_fk"
	DDL_CHECK_INDEX_COUNT                      = "ddl_check_index_count"
	DDL_CHECK_COMPOSITE_INDEX_MAX              = "ddl_check_composite_index_max"
	DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4     = "ddl_check_table_without_innodb_utf8mb4"
	DDL_CHECK_INDEX_COLUMN_WITH_BLOB           = "ddl_check_index_column_with_blob"
	DDL_CHECK_ALTER_TABLE_NEED_MERGE           = "ddl_check_alter_table_need_merge"
	DDL_DISABLE_DROP_STATEMENT                 = "ddl_disable_drop_statement"
	DML_CHECK_WHERE_IS_INVALID                 = "all_check_where_is_invalid"
	DML_DISABE_SELECT_ALL_COLUMN               = "dml_disable_select_all_column"
	DDL_CHECK_TABLE_WITHOUT_COMMENT            = "ddl_check_table_without_comment"
	DDL_CHECK_COLUMN_WITHOUT_COMMENT           = "ddl_check_column_without_comment"
	DDL_CHECK_INDEX_PREFIX                     = "ddl_check_index_prefix"
	DDL_CHECK_UNIQUE_INDEX_PRIFIX              = "ddl_check_unique_index_prefix"
	DDL_CHECK_COLUMN_WITHOUT_DEFAULT           = "ddl_check_column_without_default"
	DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT = "ddl_check_column_timestamp_without_default"
	DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL        = "ddl_check_column_blob_with_not_null"
	DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL  = "ddl_check_column_blob_default_is_not_null"
	DML_CHECK_WITH_LIMIT                       = "dml_check_with_limit"
	DML_CHECK_WITH_ORDER_BY                    = "dml_check_with_order_by"
)

// inspector config code
const (
	CONFIG_DML_ROLLBACK_MAX_ROWS = "dml_rollback_max_rows"
	CONFIG_DDL_OSC_MIN_SIZE      = "ddl_osc_min_size"
)

const (
	DML_CHECK_INSERT_COLUMNS_EXIST                   = "dml_check_insert_columns_exist"
	DML_CHECK_BATCH_INSERT_LISTS_MAX                 = "dml_check_batch_insert_lists_max"
	DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT             = "ddl_check_pk_prohibit_auto_increment"
	DML_CHECK_WHERE_EXIST_FUNC                       = "dml_check_where_exist_func"
	DML_CHECK_WHERE_EXIST_NOT                        = "dml_check_where_exist_not"
	DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION        = "dml_check_where_exist_implicit_conversion"
	DML_CHECK_LIMIT_MUST_EXIST                       = "dml_check_limit_must_exist"
	DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES         = "dml_check_where_exist_scalar_sub_queries"
	DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS = "ddl_check_indexes_exist_before_creat_constraints"
	DML_CHECK_SELECT_FOR_UPDATE                      = "dml_check_select_for_update"
	DDL_CHECK_COLLATION_DATABASE                     = "ddl_check_collation_database"
	DDL_CHECK_DECIMAL_TYPE_COLUMN                    = "ddl_check_decimal_type_column"
)

type RuleHandler struct {
	Rule          model.Rule
	Message       string
	Func          func(model.Rule, *Inspect, ast.Node) error
	IsDefaultRule bool
}

var (
	RuleHandlerMap       = map[string]RuleHandler{}
	DefaultTemplateRules = []model.Rule{}
	InitRules            = []model.Rule{}
)

var RuleHandlers = []RuleHandler{
	// config
	RuleHandler{
		Rule: model.Rule{
			Name:  CONFIG_DML_ROLLBACK_MAX_ROWS,
			Desc:  "在 DML 语句中预计影响行数超过指定值则不回滚",
			Value: "1000",
		},
		Func:          nil,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  CONFIG_DDL_OSC_MIN_SIZE,
			Desc:  "改表时，表空间超过指定大小(MB)审核时输出osc改写建议",
			Value: "16",
		},
		Func:          nil,
		IsDefaultRule: true,
	},

	// rule
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST,
			Desc:  "新建表必须加入if not exists create，保证重复执行不报错",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "新建表必须加入if not exists create，保证重复执行不报错",
		Func:          checkIfNotExist,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_OBJECT_NAME_LENGTH,
			Desc:  "表名、列名、索引名的长度不能大于64字节",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "表名、列名、索引名的长度不能大于64字节",
		Func:          checkNewObjectName,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_NOT_EXIST,
			Desc:  "表必须有主键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "表必须有主键",
		Func:          checkPrimaryKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
			Desc:  "主键建议使用自增",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "主键建议使用自增",
		Func:          checkPrimaryKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED,
			Desc:  "主键建议使用 bigint 无符号类型，即 bigint unsigned",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "主键建议使用 bigint 无符号类型，即 bigint unsigned",
		Func:          checkPrimaryKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_VARCHAR_MAX,
			Desc:  "禁止使用 varchar(max)",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止使用 varchar(max)",
		Func:          nil,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_CHAR_LENGTH,
			Desc:  "char长度大于20时，必须使用varchar类型",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "char长度大于20时，必须使用varchar类型",
		Func:          checkStringType,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_FK,
			Desc:  "禁止使用外键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止使用外键",
		Func:          checkForeignKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_COUNT,
			Desc:  "索引个数建议不超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "5",
		},
		Message:       "索引个数建议不超过%v个",
		Func:          checkIndex,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COMPOSITE_INDEX_MAX,
			Desc:  "复合索引的列数量不建议超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "3",
		},
		Message:       "复合索引的列数量不建议超过%v个",
		Func:          checkIndex,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
			Desc:  "数据库对象命名禁止使用关键字",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "数据库对象命名禁止使用关键字 %s",
		Func:          checkNewObjectName,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4,
			Desc:  "建议使用Innodb引擎,utf8mb4字符集",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "建议使用Innodb引擎,utf8mb4字符集",
		Func:          checkEngineAndCharacterSet,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_COLUMN_WITH_BLOB,
			Desc:  "禁止将blob类型的列加入索引",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止将blob类型的列加入索引",
		Func:          disableAddIndexForColumnsTypeBlob,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_IS_INVALID,
			Desc:  "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
		Func:          checkSelectWhere,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_ALTER_TABLE_NEED_MERGE,
			Desc:  "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "已存在对该表的修改语句，建议合并成一个ALTER语句",
		Func:          checkMergeAlterTable,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_DISABE_SELECT_ALL_COLUMN,
			Desc:  "不建议使用select *",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "不建议使用select *",
		Func:          checkSelectAll,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_DROP_STATEMENT,
			Desc:  "禁止除索引外的drop操作",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止除索引外的drop操作",
		Func:          disableDropStmt,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_COMMENT,
			Desc:  "表建议添加注释",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "表建议添加注释",
		Func:          checkTableWithoutComment,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_WITHOUT_COMMENT,
			Desc:  "列建议添加注释",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "列建议添加注释",
		Func:          checkColumnWithoutComment,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_PREFIX,
			Desc:  "普通索引必须要以\"idx_\"为前缀",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "普通索引必须要以\"idx_\"为前缀",
		Func:          checkIndexPrefix,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_UNIQUE_INDEX_PRIFIX,
			Desc:  "unique索引必须要以\"uniq_\"为前缀",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "unique索引必须要以\"uniq_\"为前缀",
		Func:          checkUniqIndexPrefix,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_WITHOUT_DEFAULT,
			Desc:  "除了自增列及大字段列之外，每个列都必须添加默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "除了自增列及大字段列之外，每个列都必须添加默认值",
		Func:          checkColumnWithoutDefault,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT,
			Desc:  "timestamp 类型的列必须添加默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "timestamp 类型的列必须添加默认值",
		Func:          checkColumnTimestampWithoutDefault,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
			Desc:  "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
		Func:          checkColumnBlobNotNull,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL,
			Desc:  "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
		Func:          checkColumnBlobDefaultNull,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WITH_LIMIT,
			Desc:  "delete/update 语句不能有limit条件",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "delete/update 语句不能有limit条件",
		Func:          checkDMLWithLimit,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WITH_ORDER_BY,
			Desc:  "delete/update 语句不能有order by",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "delete/update 语句不能有order by",
		Func:          checkDMLWithOrderBy,
		IsDefaultRule: true,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_INSERT_COLUMNS_EXIST,
			Desc:  "insert 语句必须指定column",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "insert 语句必须指定column",
		Func:    checkDMLWithInsertColumnExist,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_BATCH_INSERT_LISTS_MAX,
			Desc:  "单条insert语句，建议批量插入不超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "5000",
		},
		Message: "单条insert语句，建议批量插入不超过%v条",
		Func:    checkDMLWithBatchInsertMaxLimits,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT,
			Desc:  "主键禁止使用自增",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "主键禁止使用自增",
		Func:    checkPrimaryKey,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_FUNC,
			Desc:  "避免对条件字段使用函数操作",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "避免对条件字段使用函数操作",
		Func:    checkWhereExistFunc,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_NOT,
			Desc:  "不建议对条件字段使用负向查询",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议对条件字段使用负向查询",
		Func:    checkSelectWhere,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION,
			Desc:  "条件字段存在数值和字符的隐式转换",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "条件字段存在数值和字符的隐式转换",
		Func:    checkWhereColumnImplicitConversion,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_LIMIT_MUST_EXIST,
			Desc:  "delete/update 语句必须有limit条件",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "delete/update 语句必须有limit条件",
		Func:    checkDMLLimitExist,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES,
			Desc:  "避免使用标量子查询",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "避免使用标量子查询",
		Func:    checkSelectWhere,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS,
			Desc:  "建议创建约束前,先行创建索引",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议创建约束前,先行创建索引",
		Func:    checkIndexesExistBeforeCreatConstraints,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_SELECT_FOR_UPDATE,
			Desc:  "建议避免使用select for update",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议避免使用select for update",
		Func:    checkDMLSelectForUpdate,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLLATION_DATABASE,
			Desc:  "建议使用规定的数据库排序规则",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "utf8mb4_0900_ai_ci",
		},
		Message: "建议使用规定的数据库排序规则为%s",
		Func:    checkCollationDatabase,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_DECIMAL_TYPE_COLUMN,
			Desc:  "精确浮点数建议使用DECIMAL",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "精确浮点数建议使用DECIMAL",
		Func:    checkDecimalTypeColumn,
	},
}

func init() {
	for _, rh := range RuleHandlers {
		RuleHandlerMap[rh.Rule.Name] = rh
		InitRules = append(InitRules, rh.Rule)
		if rh.IsDefaultRule {
			DefaultTemplateRules = append(DefaultTemplateRules, rh.Rule)
		}
	}
}

func checkSelectAll(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// check select all column
		if stmt.Fields != nil && stmt.Fields.Fields != nil {
			for _, field := range stmt.Fields.Fields {
				if field.WildCard != nil {
					i.addResult(DML_DISABE_SELECT_ALL_COLUMN)
				}
			}
		}
	}
	return nil
}

func checkSelectWhere(rule model.Rule, i *Inspect, node ast.Node) error {
	var where ast.ExprNode
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		where = stmt.Where
	case *ast.UpdateStmt:
		where = stmt.Where
	case *ast.DeleteStmt:
		where = stmt.Where
	default:
		return nil
	}
	if where == nil || !whereStmtHasOneColumn(where) {
		i.addResult(DML_CHECK_WHERE_IS_INVALID)
	}
	if where != nil && whereStmtExistNot(where) {
		i.addResult(DML_CHECK_WHERE_EXIST_NOT)
	}
	if where != nil && whereStmtExistScalarSubQueries(where) {
		i.addResult(DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES)
	}

	return nil
}

func checkIndexesExistBeforeCreatConstraints(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		constraintMap := make(map[string]struct{})
		cols := []string{}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			if spec.Constraint != nil && (spec.Constraint.Tp == ast.ConstraintPrimaryKey ||
				spec.Constraint.Tp == ast.ConstraintUniq || spec.Constraint.Tp == ast.ConstraintUniqKey) {
				for _, key := range spec.Constraint.Keys {
					cols = append(cols, key.Column.Name.String())
				}
			}
		}
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		for _, constraints := range createTableStmt.Constraints {
			for _, key := range constraints.Keys {
				constraintMap[key.Column.Name.String()] = struct{}{}
			}
		}
		for _, col := range cols {
			if _, ok := constraintMap[col]; !ok {
				i.addResult(DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS)
				return nil
			}
		}
	}
	return nil
}

func checkPrimaryKey(rule model.Rule, i *Inspect, node ast.Node) error {
	var hasPk = false
	var pkColumnExist = false
	var pkIsAutoIncrement = false
	var pkIsBigIntUnsigned = false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}
		// check primary key
		// TODO: tidb parser not support keyword for SERIAL; it is a alias for "BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE"
		/*
			match sql like:
			CREATE TABLE  tb1 (
			a1.id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
			);
		*/
		for _, col := range stmt.Cols {
			if IsAllInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
				hasPk = true
				pkColumnExist = true
				if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) {
					pkIsBigIntUnsigned = true
				}
				if IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					pkIsAutoIncrement = true
				}
			}
		}
		/*
			match sql like:
			CREATE TABLE  tb1 (
			a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
			PRIMARY KEY (id)
			);
		*/
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				hasPk = true
				if len(constraint.Keys) == 1 {
					columnName := constraint.Keys[0].Column.Name.String()
					for _, col := range stmt.Cols {
						if col.Name.Name.String() == columnName {
							pkColumnExist = true
							if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) {
								pkIsBigIntUnsigned = true
							}
							if IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
								pkIsAutoIncrement = true
							}
						}
					}
				}
			}
		}
	default:
		return nil
	}

	if !hasPk {
		i.addResult(DDL_CHECK_PK_NOT_EXIST)
	}
	if hasPk && pkColumnExist && !pkIsAutoIncrement {
		i.addResult(DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT)
	}
	if hasPk && pkColumnExist && pkIsAutoIncrement {
		i.addResult(DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT)
	}
	if hasPk && pkColumnExist && !pkIsBigIntUnsigned {
		i.addResult(DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED)
	}

	return nil
}

func checkMergeAlterTable(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := i.getTableInfo(stmt.Table)
		if exist {
			if info.AlterTables != nil && len(info.AlterTables) > 0 {
				i.addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE)
			}
		}
	}
	return nil
}

func checkEngineAndCharacterSet(rule model.Rule, i *Inspect, node ast.Node) error {
	var engine string
	var characterSet string
	var err error
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}
		for _, op := range stmt.Options {
			switch op.Tp {
			case ast.TableOptionEngine:
				engine = op.StrValue
			case ast.TableOptionCharset:
				characterSet = op.StrValue
			}
		}
		if engine == "" {
			engine, err = i.getSchemaEngine(stmt.Table)
			if err != nil {
				return err
			}
		}
		if characterSet == "" {
			characterSet, err = i.getSchemaCharacter(stmt.Table)
			if err != nil {
				return err
			}
		}
	default:
		return nil
	}
	if strings.ToLower(engine) == "innodb" && strings.ToLower(characterSet) == "utf8mb4" {
		return nil
	}
	i.addResult(DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4)
	return nil
}

func disableAddIndexForColumnsTypeBlob(rule model.Rule, i *Inspect, node ast.Node) error {
	isTypeBlobCols := map[string]bool{}
	indexDataTypeIsBlob := false
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if MysqlDataTypeIsBlob(col.Tp.Tp) {
				if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
					indexDataTypeIsBlob = true
					break
				}
				isTypeBlobCols[col.Name.Name.String()] = true
			} else {
				isTypeBlobCols[col.Name.Name.String()] = false
			}
		}
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				for _, col := range constraint.Keys {
					if isTypeBlobCols[col.Column.Name.String()] {
						indexDataTypeIsBlob = true
						break
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		// collect columns type from original table
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if exist {
			for _, col := range createTableStmt.Cols {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		// collect columns type from alter table
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn) {
			if spec.NewColumns == nil {
				continue
			}
			for _, col := range spec.NewColumns {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
						indexDataTypeIsBlob = true
						break
					}
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniq:
				for _, col := range spec.Constraint.Keys {
					if isTypeBlobCols[col.Column.Name.String()] {
						indexDataTypeIsBlob = true
						break
					}
				}
			}
		}
	case *ast.CreateIndexStmt:
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil || !exist {
			return err
		}
		for _, col := range createTableStmt.Cols {
			if MysqlDataTypeIsBlob(col.Tp.Tp) {
				isTypeBlobCols[col.Name.Name.String()] = true
			} else {
				isTypeBlobCols[col.Name.Name.String()] = false
			}
		}
		for _, indexColumns := range stmt.IndexColNames {
			if isTypeBlobCols[indexColumns.Column.Name.String()] {
				indexDataTypeIsBlob = true
				break
			}
		}
	default:
		return nil
	}
	if indexDataTypeIsBlob {
		i.addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB)
	}
	return nil
}

func checkNewObjectName(rule model.Rule, i *Inspect, node ast.Node) error {
	names := []string{}
	invalidNames := []string{}

	switch stmt := node.(type) {
	case *ast.CreateDatabaseStmt:
		// schema
		names = append(names, stmt.Name)
	case *ast.CreateTableStmt:

		// table
		names = append(names, stmt.Table.Name.String())

		// column
		for _, col := range stmt.Cols {
			names = append(names, col.Name.Name.String())
		}
		// index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniqKey, ast.ConstraintKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
				names = append(names, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableRenameTable:
				// rename table
				names = append(names, spec.NewTable.Name.String())
			case ast.AlterTableAddColumns:
				// new column
				for _, col := range spec.NewColumns {
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableChangeColumn:
				// rename column
				for _, col := range spec.NewColumns {
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableAddConstraint:
				// if spec.Constraint.Name not index name, it will be null
				names = append(names, spec.Constraint.Name)
			case ast.AlterTableRenameIndex:
				names = append(names, spec.ToKey.String())
			}
		}
	case *ast.CreateIndexStmt:
		names = append(names, stmt.IndexName)
	default:
		return nil
	}

	// check length
	for _, name := range names {
		if len(name) > 64 {
			i.addResult(DDL_CHECK_OBJECT_NAME_LENGTH)
			break
		}
	}
	// check keyword
	for _, name := range names {
		if IsMysqlReservedKeyword(name) {
			invalidNames = append(invalidNames, name)
		}
	}
	if len(invalidNames) > 0 {
		i.addResult(DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
			strings.Join(RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func checkForeignKey(rule model.Rule, i *Inspect, node ast.Node) error {
	hasFk := false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintForeignKey {
				hasFk = true
				break
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint != nil && spec.Constraint.Tp == ast.ConstraintForeignKey {
				hasFk = true
				break
			}
		}
	default:
		return nil
	}
	if hasFk {
		i.addResult(DDL_DISABLE_FK)
	}
	return nil
}

func checkIndex(rule model.Rule, i *Inspect, node ast.Node) error {
	indexCounter := 0
	compositeIndexMax := 0
	value, err := strconv.Atoi(rule.Value)
	if err != nil {
		return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	}
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				indexCounter++
				if compositeIndexMax < len(constraint.Keys) {
					compositeIndexMax = len(constraint.Keys)
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint == nil {
				continue
			}
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				indexCounter++
				if compositeIndexMax < len(spec.Constraint.Keys) {
					compositeIndexMax = len(spec.Constraint.Keys)
				}
			}
		}
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
			}
		}

	case *ast.CreateIndexStmt:
		indexCounter++
		if compositeIndexMax < len(stmt.IndexColNames) {
			compositeIndexMax = len(stmt.IndexColNames)
		}
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
			}
		}
	default:
		return nil
	}
	if indexCounter > value {
		i.addResult(DDL_CHECK_INDEX_COUNT, value)
	}
	if compositeIndexMax > value {
		i.addResult(DDL_CHECK_COMPOSITE_INDEX_MAX, value)
	}
	return nil
}

func checkStringType(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				i.addResult(DDL_CHECK_COLUMN_CHAR_LENGTH)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					i.addResult(DDL_CHECK_COLUMN_CHAR_LENGTH)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkIfNotExist(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			i.addResult(DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST)
		}
	}
	return nil
}

func disableDropStmt(rule model.Rule, i *Inspect, node ast.Node) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	case *ast.DropTableStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	}
	return nil
}

func checkTableWithoutComment(rule model.Rule, i *Inspect, node ast.Node) error {
	var tableHasComment bool
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if has refer table, sql is create table ... like ...
		if stmt.ReferTable != nil {
			return nil
		}
		if stmt.Options != nil {
			for _, option := range stmt.Options {
				if option.Tp == ast.TableOptionComment {
					tableHasComment = true
					break
				}
			}
		}
		if !tableHasComment {
			i.addResult(DDL_CHECK_TABLE_WITHOUT_COMMENT)
		}
	}
	return nil
}

func checkColumnWithoutComment(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			columnHasComment := false
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionComment {
					columnHasComment = true
				}
			}
			if !columnHasComment {
				i.addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasComment := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionComment {
						columnHasComment = true
					}
				}
				if !columnHasComment {
					i.addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT)
					return nil
				}
			}
		}
	}
	return nil
}

func checkIndexPrefix(rule model.Rule, i *Inspect, node ast.Node) error {
	indexesName := []string{}
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex:
				indexesName = append(indexesName, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex:
				indexesName = append(indexesName, spec.Constraint.Name)
			}
		}
	case *ast.CreateIndexStmt:
		if !stmt.Unique {
			indexesName = append(indexesName, stmt.IndexName)
		}
	default:
		return nil
	}
	for _, name := range indexesName {
		if !strings.HasPrefix(name, "idx_") {
			i.addResult(DDL_CHECK_INDEX_PREFIX)
			return nil
		}
	}
	return nil
}

func checkUniqIndexPrefix(rule model.Rule, i *Inspect, node ast.Node) error {
	uniqueIndexesName := []string{}
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniq:
				uniqueIndexesName = append(uniqueIndexesName, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintUniq:
				uniqueIndexesName = append(uniqueIndexesName, spec.Constraint.Name)
			}
		}
	case *ast.CreateIndexStmt:
		if stmt.Unique {
			uniqueIndexesName = append(uniqueIndexesName, stmt.IndexName)
		}
	default:
		return nil
	}
	for _, name := range uniqueIndexesName {
		if !strings.HasPrefix(name, "uniq_") {
			i.addResult(DDL_CHECK_UNIQUE_INDEX_PRIFIX)
			return nil
		}
	}
	return nil
}

func checkColumnWithoutDefault(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col == nil {
				continue
			}
			isAutoIncrementColumn := false
			isBlobColumn := false
			columnHasDefault := false
			if HasOneInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
				isAutoIncrementColumn = true
			}
			if MysqlDataTypeIsBlob(col.Tp.Tp) {
				isBlobColumn = true
			}
			if HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
				columnHasDefault = true
			}
			if isAutoIncrementColumn || isBlobColumn {
				continue
			}
			if !columnHasDefault {
				i.addResult(DDL_CHECK_COLUMN_WITHOUT_DEFAULT)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				isAutoIncrementColumn := false
				isBlobColumn := false
				columnHasDefault := false

				if HasOneInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					isAutoIncrementColumn = true
				}
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					isBlobColumn = true
				}
				if HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
					columnHasDefault = true
				}

				if isAutoIncrementColumn || isBlobColumn {
					continue
				}
				if !columnHasDefault {
					i.addResult(DDL_CHECK_COLUMN_WITHOUT_DEFAULT)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnTimestampWithoutDefault(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			columnHasDefault := false
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionDefaultValue {
					columnHasDefault = true
				}
			}
			if !columnHasDefault && (col.Tp.Tp == mysql.TypeTimestamp || col.Tp.Tp == mysql.TypeDatetime) {
				i.addResult(DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasDefault := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionDefaultValue {
						columnHasDefault = true
					}
				}
				if !columnHasDefault && (col.Tp.Tp == mysql.TypeTimestamp || col.Tp.Tp == mysql.TypeDatetime) {
					i.addResult(DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnBlobNotNull(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp == nil {
				continue
			}
			switch col.Tp.Tp {
			case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
				for _, opt := range col.Options {
					if opt.Tp == ast.ColumnOptionNotNull {
						i.addResult(DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL)
						return nil
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
					for _, opt := range col.Options {
						if opt.Tp == ast.ColumnOptionNotNull {
							i.addResult(DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkColumnBlobDefaultNull(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp == nil {
				continue
			}
			switch col.Tp.Tp {
			case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
				for _, opt := range col.Options {
					if opt.Tp == ast.ColumnOptionDefaultValue && opt.Expr.GetType().Tp != mysql.TypeNull {
						i.addResult(DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL)
						return nil
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableModifyColumn, ast.AlterTableAlterColumn,
			ast.AlterTableChangeColumn, ast.AlterTableAddColumns) {
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
					for _, opt := range col.Options {
						if opt.Tp == ast.ColumnOptionDefaultValue && opt.Expr.GetType().Tp != mysql.TypeNull {
							i.addResult(DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkDMLWithLimit(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit != nil {
			i.addResult(DML_CHECK_WITH_LIMIT)
		}
	case *ast.DeleteStmt:
		if stmt.Limit != nil {
			i.addResult(DML_CHECK_WITH_LIMIT)
		}
	}
	return nil
}
func checkDMLLimitExist(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit == nil {
			i.addResult(DML_CHECK_LIMIT_MUST_EXIST)
		}
	case *ast.DeleteStmt:
		if stmt.Limit == nil {
			i.addResult(DML_CHECK_LIMIT_MUST_EXIST)
		}
	}
	return nil
}

func checkDMLWithOrderBy(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Order != nil {
			i.addResult(DML_CHECK_WITH_ORDER_BY)
		}
	case *ast.DeleteStmt:
		if stmt.Order != nil {
			i.addResult(DML_CHECK_WITH_ORDER_BY)
		}
	}
	return nil
}

func checkDMLWithInsertColumnExist(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Columns) == 0 {
			i.addResult(DML_CHECK_INSERT_COLUMNS_EXIST)
		}
	}
	return nil
}

func checkDMLWithBatchInsertMaxLimits(rule model.Rule, i *Inspect, node ast.Node) error {
	value, err := strconv.Atoi(rule.Value)
	if err != nil {
		return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	}
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Lists) > value {
			i.addResult(DML_CHECK_BATCH_INSERT_LISTS_MAX, value)
		}
	}
	return nil
}

func checkWhereExistFunc(rule model.Rule, i *Inspect, node ast.Node) error {

	var where ast.ExprNode
	tables := []*ast.TableName{}

	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			where = stmt.Where
			tableSources := getTableSources(stmt.From.TableRefs)
			// not select from table statement
			if len(tableSources) < 1 {
				break
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			where = stmt.Where
			tableSources := getTableSources(stmt.TableRefs.TableRefs)
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			where = stmt.Where
			tables = getTables(stmt.TableRefs.TableRefs)
		}

	default:
		return nil
	}
	if where == nil {
		return nil
	}

	var cols []*ast.ColumnDef
	for _, tableName := range tables {
		createTableStmt, exist, err := i.getCreateTableStmt(tableName)
		if exist && err == nil {
			cols = append(cols, createTableStmt.Cols...)
		}
	}
	colMap := make(map[string]struct{})
	for _, col := range cols {
		colMap[col.Name.String()] = struct{}{}
	}
	if isFuncUsedOnColumnInWhereStmt(colMap, where) {
		i.addResult(DML_CHECK_WHERE_EXIST_FUNC)
	}
	return nil

}
func checkWhereColumnImplicitConversion(rule model.Rule, i *Inspect, node ast.Node) error {
	var where ast.ExprNode
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			where = stmt.Where
			tableSources := getTableSources(stmt.From.TableRefs)
			// not select from table statement
			if len(tableSources) < 1 {
				break
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			where = stmt.Where
			tableSources := getTableSources(stmt.TableRefs.TableRefs)
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			where = stmt.Where
			tables = getTables(stmt.TableRefs.TableRefs)
		}
	default:
		return nil
	}
	if where == nil {
		return nil
	}
	var cols []*ast.ColumnDef
	for _, tableName := range tables {
		createTableStmt, exist, err := i.getCreateTableStmt(tableName)
		if exist && err == nil {
			cols = append(cols, createTableStmt.Cols...)
		}
	}
	colMap := make(map[string]string)
	for _, col := range cols {
		colType := ""
		if col.Tp == nil {
			continue
		}
		switch col.Tp.Tp {
		case mysql.TypeVarchar, mysql.TypeString:
			colType = "string"
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong, mysql.TypeLonglong, mysql.TypeDouble, mysql.TypeFloat, mysql.TypeNewDecimal:
			colType = "int"
		}
		if colType != "" {
			colMap[col.Name.String()] = colType
		}

	}
	if isColumnImplicitConversionInWhereStmt(colMap, where) {
		i.addResult(DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION)
	}
	return nil

}

func checkDMLSelectForUpdate(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.LockTp == ast.SelectLockForUpdate {
			i.addResult(DML_CHECK_SELECT_FOR_UPDATE)
		}
	}
	return nil
}

func checkCollationDatabase(rule model.Rule, i *Inspect, node ast.Node) error {

	var collationDatabase string
	var err error
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}
		for _, op := range stmt.Options {
			switch op.Tp {
			case ast.TableOptionCollate:
				collationDatabase = op.StrValue
			}
		}
		if collationDatabase == "" {
			collationDatabase, err = i.getCollationDatabase(stmt.Table)
			if err != nil {
				return err
			}
		}
	default:
		return nil
	}
	if strings.ToLower(collationDatabase) != strings.ToLower(rule.Value) {
		i.addResult(DDL_CHECK_COLLATION_DATABASE, rule.Value)
	}
	return nil
}
func checkDecimalTypeColumn(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
				i.addResult(DDL_CHECK_DECIMAL_TYPE_COLUMN)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
					i.addResult(DDL_CHECK_DECIMAL_TYPE_COLUMN)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

package inspector

import (
	"errors"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/parser"
	"sqle/executor"
	"sqle/storage"
	"strings"
)

func parseSql(dbType int, sql string) ([]ast.StmtNode, error) {
	switch dbType {
	case storage.DB_TYPE_MYSQL:
		p := parser.New()
		stmts, err := p.Parse(sql, "", "")
		if err != nil {
			fmt.Printf("parse error: %v\nsql: %v", err, sql)
			return nil, err
		}
		return stmts, nil
	default:
		return nil, errors.New("db type is invalid")
	}
}

func Inspect(config map[int]*storage.InspectConfig, task *storage.Task) ([]*storage.Sql, error) {
	sqls := []*storage.Sql{}
	stmts, err := parseSql(task.Db.DbType, task.ReqSql)
	if err != nil {
		return nil, err
	}
	db, err := executor.OpenDbWithMeta(&task.Db)
	if err != err {
		return nil, err
	}
	for _, stmt := range stmts {
		errMsgs := []string{}
		warnMsgs := []string{}
		fmt.Println("do rules")
		for _, rule := range Rules {
			fmt.Println("do rule")
			errMsg, warnMsg, err := rule.Check(config[rule.DfConfig.Code], db, stmt)
			if err != err {
				return nil, err
			}
			errMsgs = append(errMsgs, errMsg)
			warnMsgs = append(warnMsgs, warnMsg)
		}
		sql := &storage.Sql{}
		sql.CommitSql = stmt.Text()
		sql.InspectError = strings.Join(errMsgs, "\n")
		sql.InspectWarn = strings.Join(warnMsgs, "\n")
		sqls = append(sqls, sql)
	}
	return sqls, nil
}

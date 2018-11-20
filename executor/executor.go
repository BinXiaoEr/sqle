package executor

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sqle/errors"
	"sqle/model"
	"strconv"
)

type Db interface {
	Close()
	Ping() error
	Exec(query string) (driver.Result, error)
	Query(query string, args ...interface{}) ([]map[string]string, error)
}

type BaseConn struct {
	*gorm.DB
}

func newConn(instance *model.Instance, schema string) (*BaseConn, error) {
	var db *gorm.DB
	var err error
	switch instance.DbType {
	case model.DB_TYPE_MYSQL:
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			instance.User, instance.Password, instance.Host, instance.Port, schema))
	default:
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, fmt.Errorf("db type is not support"))
	}
	if err != nil {
		err = fmt.Errorf("connect to %s:%s failed, %s", instance.Host, instance.Port, err)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	return &BaseConn{db}, nil
}

func (c *BaseConn) Close() {
	c.DB.Close()
}

func (c *BaseConn) Ping() error {
	return errors.New(errors.CONNECT_REMOTE_DB_ERROR, c.DB.DB().Ping())
}

func (c *BaseConn) Exec(query string) (driver.Result, error) {
	result, err := c.DB.DB().Exec(query)
	return result, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
}

func (c *BaseConn) Query(query string, args ...interface{}) ([]map[string]string, error) {
	rows, err := c.DB.DB().Query(query, args...)
	if err != nil {
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		// unknown error
		return nil, err
	}
	result := make([]map[string]string, 0)
	for rows.Next() {
		buf := make([]interface{}, len(columns))
		data := make([]sql.NullString, len(columns))
		for i := range buf {
			buf[i] = &data[i]
		}
		if err := rows.Scan(buf...); err != nil {
			return nil, err
		}
		value := make(map[string]string, len(columns))
		for i := 0; i < len(columns); i++ {
			k := columns[i]
			v := data[i].String
			value[k] = v
		}
		result = append(result, value)
	}
	return result, nil
}

type Executor struct {
	Db Db
}

func NewExecutor(instance *model.Instance, schema string) (*Executor, error) {
	var executor = &Executor{}
	var conn Db
	var err error
	switch instance.DbType {
	case model.DB_TYPE_MYCAT:
		conn, err = newMycatConn(instance, schema)
	default:
		conn, err = newConn(instance, schema)
	}
	if err != nil {
		return nil, err
	}
	executor.Db = conn
	return executor, nil
}

func Ping(instance *model.Instance) error {
	conn, err := NewExecutor(instance, "")
	//conn, err := NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return err
	}
	defer conn.Db.Close()
	return conn.Db.Ping()
}

func ShowDatabases(instance *model.Instance) ([]string, error) {
	conn, err := NewExecutor(instance, "")
	//conn, err := NewConn(db.DbType, db.User, db.Password, db.Host, db.Port, "")
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	return conn.ShowDatabases()
}

func OpenDbWithTask(task *model.Task) (*Executor, error) {
	return NewExecutor(task.Instance, task.Schema)
}

func Exec(task *model.Task, sql string) (driver.Result, error) {
	conn, err := OpenDbWithTask(task)
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	return conn.Db.Exec(sql)
}

func (c *Executor) ShowCreateTable(tableName string) (string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show create table %s", tableName))
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		return "", errors.New(errors.CONNECT_REMOTE_DB_ERROR,
			fmt.Errorf("show create table error, result is %v", result))
	}
	if query, ok := result[0]["Create Table"]; !ok {
		return "", errors.New(errors.CONNECT_REMOTE_DB_ERROR,
			fmt.Errorf("show create table error, column \"Create Table\" not found"))
	} else {
		return query, nil
	}
}

func (c *Executor) ShowDatabases() ([]string, error) {
	result, err := c.Db.Query("show databases")
	if err != nil {
		return nil, err
	}
	dbs := make([]string, len(result))
	for n, v := range result {
		if db, ok := v["Database"]; !ok {
			return dbs, errors.New(errors.CONNECT_REMOTE_DB_ERROR,
				fmt.Errorf("show databases error, column \"Database\" not found"))
		} else {
			dbs[n] = db
		}
	}
	return dbs, nil
}

func (c *Executor) ShowSchemaTables(schema string) ([]string, error) {
	result, err := c.Db.Query("select table_name from information_schema.tables where table_schema = ?", schema)
	if err != nil {
		return nil, err
	}
	tables := make([]string, len(result))
	for n, v := range result {
		tables[n] = v["table_name"]
	}
	return tables, nil
}

type ExecutionPlanJson struct {
	QueryBlock struct {
		CostInfo struct {
			QueryCost string `json:"query_cost"`
		} `json:"cost_info"`
	} `json:"query_block"`
}

func (c *Executor) Explain(query string) (ExecutionPlanJson, error) {
	ep := ExecutionPlanJson{}
	result, err := c.Db.Query(fmt.Sprintf("EXPLAIN FORMAT=\"json\" %s", query))
	if err != nil {
		return ep, err
	}
	if len(result) == 1 {
		json.Unmarshal([]byte(result[0]["EXPLAIN"]), &ep)
	}
	return ep, nil
}

func (c *Executor) ShowMasterStatus() ([]map[string]string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show master status"))
	if err != nil {
		return nil, err
	}
	// result may be empty
	if len(result) != 1 && len(result) != 0 {
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR,
			fmt.Errorf("show master status error, result is %v", result))
	}
	return result, nil
}

func (c *Executor) FetchMasterBinlogPos() (string, int64, error) {
	result, err := c.ShowMasterStatus()
	if err != nil {
		return "", 0, err
	}
	if len(result) == 0 {
		return "", 0, nil
	}
	file := result[0]["File"]
	pos, err := strconv.ParseInt(result[0]["Position"], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return file, pos, nil
}

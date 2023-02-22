package v1

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/v1/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/pkg/errors"
)

const (
	SQLTypeDML = "dml"
	SQLTypeDDL = "ddl"
)

const (
	DriverTypeMySQL      = "MySQL"
	DriverTypePostgreSQL = "PostgreSQL"
	DriverTypeTiDB       = "TiDB"
	DriverTypeSQLServer  = "SQL Server"
	DriverTypeOracle     = "Oracle"
	DriverTypeDB2        = "DB2"
	DriverTypeOceanBase  = "OceanBase For MySQL"
)

type RuleLevel string

const (
	RuleLevelNull   RuleLevel = "" // used to indicate no rank
	RuleLevelNormal RuleLevel = "normal"
	RuleLevelNotice RuleLevel = "notice"
	RuleLevelWarn   RuleLevel = "warn"
	RuleLevelError  RuleLevel = "error"
)

var ruleLevelMap = map[RuleLevel]int{
	RuleLevelNull:   -1,
	RuleLevelNormal: 0,
	RuleLevelNotice: 1,
	RuleLevelWarn:   2,
	RuleLevelError:  3,
}

func (r RuleLevel) LessOrEqual(l RuleLevel) bool {
	return ruleLevelMap[r] <= ruleLevelMap[l]
}

func (r RuleLevel) More(l RuleLevel) bool {
	return ruleLevelMap[r] > ruleLevelMap[l]
}

func (r RuleLevel) MoreOrEqual(l RuleLevel) bool {
	return ruleLevelMap[r] >= ruleLevelMap[l]
}

// RuleLevelLessOrEqual return level a <= level b
func RuleLevelLessOrEqual(a, b string) bool {
	return RuleLevel(a).LessOrEqual(RuleLevel(b))
}

type Rule struct {
	Name       string
	Desc       string
	Annotation string

	// Category is the category of the rule. Such as "Naming Conventions"...
	// Rules will be displayed on the SQLE rule list page by category.
	Category string
	Level    RuleLevel
	Params   params.Params
}

//func (r *Rule) GetValueInt(defaultRule *Rule) int64 {
//	value := r.getValue(DefaultSingleParamKeyName, defaultRule)
//	i, err := strconv.ParseInt(value, 10, 64)
//	if err != nil {
//		return 0
//	}
//	return i
//}
//
//func (r *Rule) GetSingleValue() string {
//	value, _ := r.Params.GetParamValue(DefaultSingleParamKeyName)
//	return value
//}
//
//func (r *Rule) GetSingleValueInt() int {
//	value := r.GetSingleValue()
//	i, err := strconv.Atoi(value)
//	if err != nil {
//		return 0
//	}
//	return i
//}

//func (r *Rule) getValue(key string, defaultRule *Rule) string {
//	var value string
//	var exist bool
//	value, exist = r.Params.GetParamValue(key)
//	if !exist {
//		value, _ = defaultRule.Params.GetParamValue(key)
//	}
//	return value
//}

// NewConfig return a config for driver.
//
// 1. dsn is nil, rules is not nil. Use drive to do Offline Audit.
// 2. dsn is not nil, rule is nil. Use drive to communicate with database only.
// 3. dsn is not nil, rule is not nil. Most common usecase.
func NewConfig(dsn *DSN, rules []*Rule) (*Config, error) {
	if dsn == nil && rules == nil {
		fmt.Println("dsn is nil, and rules is nil, nothing can be done by driver")
	}

	return &Config{
		DSN:   dsn,
		Rules: rules,
	}, nil
}

type DriverNotSupportedError struct {
	DriverTyp string
}

func (e *DriverNotSupportedError) Error() string {
	return fmt.Sprintf("driver type %v is not supported", e.DriverTyp)
}

var ErrNodesCountExceedOne = errors.New("after parse, nodes count exceed one")

// Driver is a interface that must be implemented by a database.
//
// It's implementation maybe on the same process or over gRPC(by go-plugin).
//
// Driver is responsible for two primary things:
// 1. provides handle to communicate with database
// 2. audit SQL with rules
type Driver interface {
	Close(ctx context.Context)
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string) (driver.Result, error)
	Tx(ctx context.Context, queries ...string) ([]driver.Result, error)

	// Schemas export all supported schemas.
	//
	// For example, performance_schema/performance_schema... which in MySQL is not allowed for auditing.
	Schemas(ctx context.Context) ([]string, error)

	// Parse parse sqlText to Node array.
	//
	// sqlText may be single SQL or batch SQLs.
	Parse(ctx context.Context, sqlText string) ([]Node, error)

	// Audit sql with rules. sql is single SQL text.
	//
	// Multi Audit call may be in one context.
	// For example:
	//		driver, _ := NewDriver(..., ..., ...)
	// 		driver.Audit(..., "CREATE TABLE t1(id int)")
	// 		driver.Audit(..., "SELECT * FROM t1 WHERE id = 1")
	//      ...
	// driver should keep SQL context during it's lifecycle.
	Audit(ctx context.Context, sql string) (*AuditResults, error)

	// GenRollbackSQL generate sql's rollback SQL.
	GenRollbackSQL(ctx context.Context, sql string) (string, string, error)
}

// Registerer is the interface that all SQLe plugins must support.
type Registerer interface {
	// Name returns plugin name.
	Name() string

	// Rules returns all rules that plugin supported.
	Rules() []*Rule

	// AdditionalParams returns all additional params that plugin supported.
	AdditionalParams() params.Params
}

// Node is a interface which unify SQL ast tree. It produce by Driver.Parse.
type Node struct {
	// Text is the raw SQL text of Node.
	Text string

	// Type is type of SQL, such as DML/DDL/DCL.
	Type string

	// Fingerprint is fingerprint of Node's raw SQL.
	Fingerprint string
}

// // DSN like https://github.com/go-sql-driver/mysql/blob/master/dsn.go. type Config struct
// type DSN struct {
// 	Type string

// 	Host   string
// 	Port   string
// 	User   string
// 	Pass   string
// 	DBName string
// }

type AuditResults struct {
	Results []*AuditResult
}

type AuditResult struct {
	Level   RuleLevel
	Message string
}

func NewInspectResults() *AuditResults {
	return &AuditResults{
		Results: []*AuditResult{},
	}
}

// Level find highest Level in result
func (rs *AuditResults) Level() RuleLevel {
	level := RuleLevelNull
	for _, curr := range rs.Results {
		if ruleLevelMap[curr.Level] > ruleLevelMap[level] {
			level = curr.Level
		}
	}
	return level
}

func (rs *AuditResults) Message() string {
	repeatCheck := map[string]struct{}{}
	messages := []string{}
	for _, result := range rs.Results {
		token := result.Message + string(result.Level)
		if _, ok := repeatCheck[token]; ok {
			continue
		}
		repeatCheck[token] = struct{}{}

		var message string
		match, _ := regexp.MatchString(fmt.Sprintf(`^\[%s|%s|%s|%s|%s\]`,
			RuleLevelError, RuleLevelWarn, RuleLevelNotice, RuleLevelNormal, "osc"),
			result.Message)
		if match {
			message = result.Message
		} else {
			message = fmt.Sprintf("[%s]%s", result.Level, result.Message)
		}
		messages = append(messages, message)
	}
	return strings.Join(messages, "\n")
}

func (rs *AuditResults) Add(level RuleLevel, message string, args ...interface{}) {
	if level == "" || message == "" {
		return
	}

	rs.Results = append(rs.Results, &AuditResult{
		Level:   level,
		Message: fmt.Sprintf(message, args...),
	})
	rs.SortByLevel()
}

func (rs *AuditResults) SortByLevel() {
	sort.Slice(rs.Results, func(i, j int) bool {
		return rs.Results[i].Level.More(rs.Results[j].Level)
	})
}

func (rs *AuditResults) HasResult() bool {
	return len(rs.Results) != 0
}

// driverImpl implement Driver. It use for hide gRPC detail, just like DriverGRPCServer.
type driverImpl struct {
	plugin proto.DriverClient

	// driverQuitCh produce a singal for telling caller that it's time to Client.Kill() plugin process.
	driverQuitCh chan struct{}
}

func (s *driverImpl) Close(ctx context.Context) {
	s.plugin.Close(ctx, &proto.Empty{})
	close(s.driverQuitCh)
}

func (s *driverImpl) Ping(ctx context.Context) error {
	_, err := s.plugin.Ping(ctx, &proto.Empty{})
	return err
}

type dbDriverResult struct {
	lastInsertId    int64
	lastInsertIdErr string
	rowsAffected    int64
	rowsAffectedErr string
}

func (s *dbDriverResult) LastInsertId() (int64, error) {
	if s.lastInsertIdErr != "" {
		return s.lastInsertId, fmt.Errorf(s.lastInsertIdErr)
	}
	return s.lastInsertId, nil
}

func (s *dbDriverResult) RowsAffected() (int64, error) {
	if s.rowsAffectedErr != "" {
		return s.rowsAffected, fmt.Errorf(s.rowsAffectedErr)
	}
	return s.rowsAffected, nil
}

func (s *driverImpl) Exec(ctx context.Context, query string) (driver.Result, error) {
	resp, err := s.plugin.Exec(ctx, &proto.ExecRequest{Query: query})
	if err != nil {
		return nil, err
	}
	return &dbDriverResult{
		lastInsertId:    resp.LastInsertId,
		lastInsertIdErr: resp.LastInsertIdError,
		rowsAffected:    resp.RowsAffected,
		rowsAffectedErr: resp.RowsAffectedError,
	}, nil
}

func (s *driverImpl) Tx(ctx context.Context, queries ...string) ([]driver.Result, error) {
	resp, err := s.plugin.Tx(ctx, &proto.TxRequest{Queries: queries})
	if err != nil {
		return nil, err
	}

	ret := make([]driver.Result, len(resp.Results))
	for i, result := range resp.Results {
		ret[i] = &dbDriverResult{
			lastInsertId:    result.LastInsertId,
			lastInsertIdErr: result.LastInsertIdError,
			rowsAffected:    result.RowsAffected,
			rowsAffectedErr: result.RowsAffectedError,
		}
	}
	return ret, nil
}

func (s *driverImpl) Schemas(ctx context.Context) ([]string, error) {
	resp, err := s.plugin.Databases(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Databases, nil
}

func (s *driverImpl) Parse(ctx context.Context, sqlText string) ([]Node, error) {
	resp, err := s.plugin.Parse(ctx, &proto.ParseRequest{SqlText: sqlText})
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, len(resp.Nodes))
	for i, node := range resp.Nodes {
		nodes[i] = Node{
			Type:        node.Type,
			Text:        node.Text,
			Fingerprint: node.Fingerprint,
		}
	}
	return nodes, nil
}

func (s *driverImpl) Audit(ctx context.Context, sql string) (*AuditResults, error) {
	resp, err := s.plugin.Audit(ctx, &proto.AuditRequest{Sql: sql})
	if err != nil {
		return nil, err
	}

	ret := &AuditResults{}
	for _, result := range resp.Results {
		ret.Results = append(ret.Results, &AuditResult{
			Level:   RuleLevel(result.Level),
			Message: result.Message,
		})
	}
	return ret, nil
}

func (s *driverImpl) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	resp, err := s.plugin.GenRollbackSQL(ctx, &proto.GenRollbackSQLRequest{Sql: sql})
	if err != nil {
		return "", "", err
	}

	return resp.Sql, resp.Reason, nil
}

func convertRuleFromProtoToDriver(rule *proto.Rule) *Rule {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		ps = append(ps, &params.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  params.ParamType(p.Type),
		})
	}
	return &Rule{
		Name:       rule.Name,
		Category:   rule.Category,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      RuleLevel(rule.Level),
		Params:     ps,
	}
}

func convertRuleFromDriverToProto(rule *Rule) *proto.Rule {
	var params = make([]*proto.Param, 0, len(rule.Params))
	for _, p := range rule.Params {
		params = append(params, &proto.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  string(p.Type),
		})
	}
	return &proto.Rule{
		Name:       rule.Name,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      string(rule.Level),
		Category:   rule.Category,
		Params:     params,
	}
}

package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"

	"github.com/labstack/echo/v4"
)

type GetSQLQueryHistoryReqV1 struct {
	FilterFuzzySearch string `json:"filter_fuzzy_search" query:"filter_fuzzy_search"`
	PageIndex         uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize          uint32 `json:"page_size" query:"page_size" valid:"required"`
}

type GetSQLQueryHistoryResV1 struct {
	controller.BaseRes
	Data GetSQLQueryHistoryResDataV1 `json:"data"`
}

type GetSQLQueryHistoryResDataV1 struct {
	SQLHistories []SQLHistoryItemResV1 `json:"sql_histories"`
}

type SQLHistoryItemResV1 struct {
	SQL string `json:"sql"`
}

// GetSQLQueryHistory get current user sql query history
// @Summary 获取当前用户历史查询SQL
// @Description get sql query history
// @Id getSQLQueryHistory
// @Tags sql_query
// @Param instance_name path string true "instance name"
// @Param filter_fuzzy_search query string false "fuzzy search filter"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLQueryHistoryResV1
// @router /v1/sql_query/{instance_name}/history [get]
func GetSQLQueryHistory(c echo.Context) error {
	return nil
}

type GetSQLResultReqV1 struct {
	PageIndex uint32 `json:"page_index" query:"page_index" valid:"required"`
	PageSize  uint32 `json:"page_size" query:"page_size" valid:"required"`
	SQLId     string `json:"sql_id" from:"sql_id"`
}

type GetSQLResultResV1 struct {
	controller.BaseRes
	Data GetSQLResultResDataV1 `json:"data"`
}

type GetSQLResultResDataV1 struct {
	ExecuteResult []SQLResultItemResV1 `json:"execute_result"`
}

// multiple SQLs may be passed in, and each SQL corresponds to an Item
type SQLResultItemResV1 struct {
	SQL         string                               `json:"sql"`
	StartLine   int                                  `json:"start_line"`
	EndLine     int                                  `json:"end_line"`
	CurrentPage int                                  `json:"current_page"`
	ExecuteTime int                                  `json:"execution_time"`
	Rows        []map[string] /* head name */ string `json:"rows"`
	Head        []SQLResultItemHeadResV1             `json:"head"`
}

type SQLResultItemHeadResV1 struct {
	HeadName string `json:"head_name"`
}

// GetSQLResult get sql query result
// @Summary 获取SQL查询结果
// @Description get sql query result
// @Id getSQLQueryResult
// @Tags sql_query
// @Param instance_name path string true "instance name"
// @Param sql_id query string false "sql id"
// @Param page_index query uint32 false "page index"
// @Param page_size query uint32 false "size of per page"
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetSQLResultResV1
// @router /v1/sql_query/{instance_name}/result [get]
func GetSQLResult(c echo.Context) error {
	return nil
}

type PrepareSQLQueryReqV1 struct {
	SQL string `json:"sql" from:"sql"`
}

type PrepareSQLQueryResV1 struct {
	controller.BaseRes
	Data PrepareSQLQueryResDataV1 `json:"data"`
}

type PrepareSQLQueryResDataV1 struct {
	SQLIds   []PrepareSQLQueryResSQLV1          `json:"sql_ids"`   // When there is no wrong SQL, an SQL ID will be generated for each SQL
	ErrorSQL []PrepareSQLQueryResErrorSQLItemV1 `json:"error_sql"` // SQL not allowed
}

type PrepareSQLQueryResSQLV1 struct {
	SQL   string `json:"sql"`
	SQLId string `json:"sql_id"`
}

type PrepareSQLQueryResErrorSQLItemV1 struct {
	SQL   string `json:"sql"`
	Error string `json:"error"`
}

// PrepareSQLQuery prepare execute sql query
// @Summary 准备执行查询sql
// @Accept json
// @Description execute sql query
// @Id execSQLQuery
// @Tags sql_query
// @Param instance_name path string true "instance name"
// @Param req body v1.PrepareSQLQueryReqV1 true "exec sql"
// @Security ApiKeyAuth
// @Success 200 {object} v1.PrepareSQLQueryResV1
// @router /v1/sql_query/{instance_name}/prepare [post]
func PrepareSQLQuery(c echo.Context) error {
	return nil
}

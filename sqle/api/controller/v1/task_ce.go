//go:build !enterprise
// +build !enterprise

package v1

import (
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func getTaskAnalysisData(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.NewNotSupportGetTaskAnalysisDataErr())
}

func getSqlFileOrderMethod(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errors.NewNotSupportGetSqlFileOrderMethodErr())
}

func sortAuditFiles(auditFiles []*model.AuditFile, orderMethod string) {}

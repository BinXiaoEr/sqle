package middleware

import (
	"errors"
	"fmt"
	"github.com/actiontech/sqle/sqle/model"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTTokenAdapter is a `echo` middleware,　by rewriting the header, the jwt token support header
// "Authorization: {token}" and "Authorization: Bearer {token}".
func JWTTokenAdapter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			if auth != "" && !strings.HasPrefix(auth, middleware.DefaultJWTConfig.AuthScheme) {
				c.Request().Header.Set(echo.HeaderAuthorization,
					fmt.Sprintf("%s %s", middleware.DefaultJWTConfig.AuthScheme, auth))
			}
			return next(c)
		}
	}
}

var errAuditPlanMisMatch = errors.New("audit plan name don't match the token")
var errAuditPlanNotFound = errors.New("audit plan not found")
var errAuditPlanTokenIncorrect = errors.New("audit plan token incorrect")

// ScannerVerifier is a `echo` middleware. Every audit plan should be
// scanner-scoped which means that scanner-A should not push SQL to scanner-B.
func ScannerVerifier() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWT parser expect no 'Bearer' ahead of token, so
			// we cut the leading auth schema.
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			parts := strings.Split(auth, " ")
			token := parts[0]
			if len(parts) == 2 {
				token = parts[1]
			}

			apnInToken, err := utils.ParseAuditPlanName(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			apnInParam := c.Param("audit_plan_name")
			if apnInToken != apnInParam {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanMisMatch.Error())
			}

			apn, apnExist, err := model.GetStorage().GetAuditPlanByName(apnInParam)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if !apnExist {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanNotFound.Error())
			}
			if apn.Token != token {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanTokenIncorrect.Error())
			}

			return next(c)
		}
	}
}

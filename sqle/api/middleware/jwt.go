package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	dmsCommonJwt "github.com/actiontech/dms/pkg/dms-common/api/jwt"
	"github.com/actiontech/sqle/sqle/model"
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
			// sqle-token为空时，可能是cookie过期被清理了，希望返回的错误是http.StatusUnauthorized
			// 但sqle-token为空时jwt返回的错误是http.StatusBadRequest
			_, err := c.Cookie("dms-token")
			if err == http.ErrNoCookie && auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "can not find dms-token")
			}
			return next(c)
		}
	}
}

func JWTWithConfig(key interface{}) echo.MiddlewareFunc {
	c := middleware.DefaultJWTConfig
	c.SigningKey = key
	c.TokenLookup = "cookie:dms-token,header:Authorization" // tell the middleware where to get token: from cookie and header
	return middleware.JWTWithConfig(c)
}

var errAuditPlanMisMatch = errors.New("audit plan don't match the token or audit plan not found")

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

			// check token belong to instance audit plan
			iapidInToken, err := dmsCommonJwt.ParseAuditPlanName(token)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if iapidInToken != utils.Md5(c.Param("instance_audit_plan_id")) {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanMisMatch.Error())
			}

			iapidInParam, err := strconv.Atoi(c.Param("instance_audit_plan_id"))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			aptypParam := c.Param("audit_plan_type")

			apn, err := model.GetStorage().GetAuditPlanDetailByIDType(iapidInParam, aptypParam)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if apn.Token != token {
				return echo.NewHTTPError(http.StatusInternalServerError, errAuditPlanMisMatch.Error())
			}

			return next(c)
		}
	}
}

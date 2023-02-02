package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/labstack/echo/v4"
)

// AuthLogin is the resolver for the authLogin field.
func (r *QueryResolverImpl) AuthLogin(ctx context.Context, provider string, configuration *string, credentials interface{}, linkUser *bool) (*model.AuthInfo, error) {
	mp, ok := credentials.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("credentials format is incorrect")
	}
	cbUser := fmt.Sprintf("%v", mp["user"])

	l := log.NewEntry()

	// 同步信息
	err := service.SyncCurrentUser(cbUser)
	if err != nil {
		l.Errorf("sync cloudbeaver user %v info failed: %v", cbUser, err)
	}
	err = service.SyncUserBindInstance(cbUser)
	if err != nil {
		l.Errorf("sync cloudbeaver user %v bind instance failed: %v", cbUser, err)
	}

	data, err := r.Next(r.Ctx)
	if err != nil {
		return nil, err
	}
	resp := &struct {
		Data struct {
			AuthInfo *model.AuthInfo `json:"authInfo"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp.Data.AuthInfo, RedirectCookie(r.Ctx)
}

// login页面无法访问sql_query页面的cookie, 这将导致登录SQLE时无法判断CloudBeaver当前登陆状态, 所以需要将cookie放到根目录下, 使用时在还原回原来的位置, 还原位置: cloudbeaver_wrapper.RedirectCookie
func RedirectCookie(c echo.Context) error {
	// cookie示例: cb-session-id=yhl4lmtzm4mo19gp72pcfkswn90; Path=/sql_query/
	cookie := c.Response().Header().Get("Set-Cookie")
	// len("cb-session-id=") = 14
	start := 14
	end := strings.Index(cookie, ";")
	if cookie != "" && end > start {
		c.Response().Header().Del("Set-Cookie")
		newCookie := cookie[start:end]
		c.SetCookie(&http.Cookie{
			Name:  "cb-session-id-sqle",
			Value: newCookie,
			Path:  "/",
		})
	}
	return nil
}

// ActiveUser is the resolver for the activeUser field.
func (r *QueryResolverImpl) ActiveUser(ctx context.Context) (*model.UserInfo, error) {
	data, err := r.Next(r.Ctx)
	if err != nil {
		return nil, err
	}

	resp := &struct {
		Data struct {
			User *model.UserInfo `json:"user"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Data.User != nil && resp.Data.User.DisplayName != nil {
		*resp.Data.User.DisplayName = service.RestoreFromCloudBeaverUserName(*resp.Data.User.DisplayName)
	}

	return resp.Data.User, err
}

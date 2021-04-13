package controller

import (
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"

	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"github.com/labstack/echo/v4"
)

type BaseRes struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message" example:"ok"`
}

func NewBaseReq(err error) BaseRes {
	res := BaseRes{}
	switch e := err.(type) {
	case *errors.CodeError:
		res.Code = e.Code()
		res.Message = e.Error()
	default:
		if err == nil {
			res.Code = 0
			res.Message = "ok"
		} else {
			res.Code = -1
			res.Message = e.Error()
		}
	}
	return res
}

func BindAndValidateReq(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := Validate(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func GetUserName(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["name"].(string)
}

func GetCurrentUser(c echo.Context) (*model.User, error) {
	key := "current_user"
	currentUser := c.Get(key)
	if currentUser != nil {
		if user, ok := currentUser.(*model.User); ok {
			return user, nil
		}
	}
	s := model.GetStorage()
	user, exist, err := s.GetUserByName(GetUserName(c))
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(errors.DataNotExist,
			fmt.Errorf("current user is not exist"))
	}
	c.Set(key, user)
	return user, nil
}

func JSONBaseErrorReq(c echo.Context, err error) error {
	return c.JSON(http.StatusOK, NewBaseReq(err))
}

func ReadFileToByte(c echo.Context, name string) (fileName string, data []byte, err error) {
	file, err := c.FormFile(name)
	if err != nil {
		err = errors.New(errors.READ_UPLOAD_FILE_ERROR, err)
		return
	}
	src, err := file.Open()
	if err != nil {
		err = errors.New(errors.READ_UPLOAD_FILE_ERROR, err)
		return
	}
	defer src.Close()
	data, err = ioutil.ReadAll(src)
	if err != nil {
		err = errors.New(errors.READ_UPLOAD_FILE_ERROR, err)
		return
	}
	return
}

func unescapeParamString(params []*string) error {
	for i, p := range params {
		r, err := url.QueryUnescape(*p)
		if nil != err {
			return fmt.Errorf("unescape param [%v] failed: %v", params, err)
		}
		*params[i] = r
	}
	return nil
}

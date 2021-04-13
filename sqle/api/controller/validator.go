package controller

import (
	"fmt"
	"github.com/go-playground/validator"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

func Validate(i interface{}) error {
	if defaultCustomValidator == nil {
		return nil
	}
	cv := defaultCustomValidator
	err := cv.validate.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	errMsgs := make([]string, 0, len(errs))
	for _, err := range errs {
		errMsgs = append(errMsgs, err.Translate(cv.enTrans))
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}
	return nil
}

var defaultCustomValidator *CustomValidator

var (
	ValidNameTag = "name"
	ValidPortTag = "port"
)

func init() {
	cv := NewCustomValidator()

	cv.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	cv.validate.SetTagName("valid")

	// register custom validator
	cv.validate.RegisterValidation(ValidNameTag, ValidateName)
	cv.validate.RegisterValidation(ValidPortTag, ValidatePort)

	type registerTranslationArgs struct {
		tag    string
		enText string
		zhText string
		params []string
	}
	argsMap := []registerTranslationArgs{
		{
			tag:    "required",
			enText: "{0} is required",
			zhText: "{0}不能为空",
		},
		{
			tag:    "oneof",
			enText: "{0} must be one of [{1}]",
			zhText: "{0}必须是[{1}]其中之一",
		},
		{
			tag:    "email",
			enText: "{0} is invalid email addr",
			zhText: "{0}是无效的邮箱地址",
		},
		{
			tag:    ValidNameTag,
			enText: "{0} must match regexp `{1}`",
			zhText: "{0}必须匹配正则`{1}`",
			params: []string{ValidateNameRegexpPattern},
		},
		{
			tag:    ValidPortTag,
			enText: "{0} is invalid port",
			zhText: "{0}是无效的端口",
		},
	}
	// register custom validator error message
	for _, args := range argsMap {
		err := cv.RegisterTranslation(args.tag, args.enText, args.enText, args.params...)
		if err != nil {
			panic(err)
		}
	}
	defaultCustomValidator = cv
}

type CustomValidator struct {
	validate *validator.Validate
	enTrans  ut.Translator
	zhTrans  ut.Translator
}

func NewCustomValidator() *CustomValidator {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, zh)
	enTrans, _ := uni.GetTranslator(en.Locale())
	zhTrans, _ := uni.GetTranslator(zh.Locale())

	cv := &CustomValidator{
		validate: validator.New(),
		enTrans:  enTrans,
		zhTrans:  zhTrans,
	}
	return cv
}

func (cv *CustomValidator) RegisterTranslation(tag, enText, zhText string, params ...string) error {
	err := cv.validate.RegisterTranslation(tag, cv.enTrans, func(ut ut.Translator) error {
		return ut.Add(tag, enText, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		tParams := []string{fe.Namespace()}
		if len(params) > 0 {
			tParams = append(tParams, params...)
		}
		tParams = append(tParams, fe.Param())
		t, _ := ut.T(tag, tParams...)
		return t
	})
	if err != nil {
		return err
	}
	err = cv.validate.RegisterTranslation(tag, cv.zhTrans, func(ut ut.Translator) error {
		return ut.Add(tag, zhText, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		tParams := []string{fe.Namespace()}
		if len(params) > 0 {
			tParams = append(tParams, params...)
		}
		tParams = append(tParams, fe.Param())
		t, _ := ut.T(tag, tParams...)
		return t
	})
	return err
}

var ValidateNameRegexpPattern = `^[a-zA-Z][a-zA-Z0-9\_\-]{0,59}$`

// ValidateName implements validator.Func
func ValidateName(fl validator.FieldLevel) bool {
	return validateName(fl.Field().String())
}

func validateName(name string) bool {
	match, _ := regexp.MatchString(ValidateNameRegexpPattern, name)

	return match
}

// ValidatePort implements validator.Func
func ValidatePort(fl validator.FieldLevel) bool {
	return validatePort(fl.Field().String())
}

func validatePort(port string) bool {
	// Port must be a iny <= 65535.
	portNum, err := strconv.ParseInt(port, 10, 32)
	if err != nil || portNum > 65535 || portNum < 1 {
		return false
	}
	return true
}

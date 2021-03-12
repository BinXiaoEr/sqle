package v1

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api/controller"
	"actiontech.cloud/universe/sqle/v4/sqle/api/server"
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/executor"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"actiontech.cloud/universe/ucommon/v4/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateInstanceReqV1 struct {
	Name                 string   `json:"instance_name" form:"instance_name" example:"test" valid:"required"`
	User                 string   `json:"db_user" form:"db_user" example:"root" valid:"required"`
	Host                 string   `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"required,ipv4"`
	Port                 string   `json:"db_port" form:"db_port" example:"3306" valid:"required,range(1|65535)"`
	Password             string   `json:"db_password" form:"db_password" example:"123456" valid:"required"`
	Desc                 string   `json:"desc" example:"this is a test instance" valid:"-"`
	WorkflowTemplateName string   `json:"workflow_template_name" form:"workflow_template_name"`
	RuleTemplates        []string `json:"rule_template_name_list" form:"rule_template_name_list" valid:"-"`
	Roles                []string `json:"role_name_list" form:"role_name_list"`
}

// @Summary 添加实例
// @Description create a instance
// @Id createInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Accept json
// @Param instance body v1.CreateInstanceReqV1 true "add instance"
// @Success 200 {object} controller.BaseRes
// @router /v1/instances [post]
func CreateInstance(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateInstanceReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	_, exist, err := s.GetInstanceByName(req.Name)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_EXIST, fmt.Errorf("instance is exist")))
	}

	instance := &model.Instance{
		DbType:   model.DB_TYPE_MYSQL,
		Name:     req.Name,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
		Desc:     req.Desc,
	}
	templates, err := s.GetAndCheckRuleTemplateExist(req.RuleTemplates)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	roles, err := s.GetAndCheckRoleExist(req.Roles)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.Save(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateInstanceRoles(instance, roles...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = s.UpdateInstanceRuleTemplates(instance, templates...)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go server.GetSqled().UpdateAndGetInstanceStatus(log.NewEntry(), instance)

	return c.JSON(200, controller.NewBaseReq(nil))
}

type InstanceResV1 struct {
	Name                 string   `json:"instance_name"`
	Host                 string   `json:"db_host" gorm:"not null" example:"10.10.10.10"`
	Port                 string   `json:"db_port" gorm:"not null" example:"3306"`
	User                 string   `json:"db_user" gorm:"not null" example:"root"`
	Desc                 string   `json:"desc" example:"this is a instance"`
	RuleTemplates        []string `json:"rule_template_name_list,omitempty"`
	WorkflowTemplateName string   `json:"workflow_template_name"`
	Roles                []string `json:"role_name_list"`
}

type GetInstanceResV1 struct {
	controller.BaseRes
	Data InstanceResV1 `json:"data"`
}

// @Summary 获取实例信息
// @Description get instance db
// @Id getInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceResV1
// @router /v1/instances/{instance_name}/ [get]
func GetInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceDetailByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("instance is not exist")))
	}

	instanceResV1 := InstanceResV1{
		Name: instance.Name,
		Host: instance.Host,
		Port: instance.Port,
		User: instance.User,
		Desc: instance.Desc,
	}
	if len(instance.RuleTemplates) > 0 {
		ruleTemplateNames := make([]string, 0, len(instance.RuleTemplates))
		for _, rt := range instance.RuleTemplates {
			ruleTemplateNames = append(ruleTemplateNames, rt.Name)
		}
		instanceResV1.RuleTemplates = ruleTemplateNames
	}
	if len(instance.Roles) > 0 {
		roleNames := make([]string, 0, len(instance.Roles))
		for _, r := range instance.Roles {
			roleNames = append(roleNames, r.Name)
		}
		instanceResV1.Roles = roleNames
	}
	return c.JSON(200, &GetInstanceResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceResV1,
	})
}

// @Summary 删除实例
// @Description delete instance db
// @Id deleteInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} controller.BaseRes
// @router /v1/instances/{instance_name}/ [delete]
func DeleteInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("instance is not exist")))
	}
	err = s.Delete(instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	server.GetSqled().DeleteInstanceStatus(instance)
	return c.JSON(200, controller.NewBaseReq(nil))
}

type UpdateInstanceReqV1 struct {
	User               *string  `json:"db_user" form:"db_user" example:"root"`
	Host               *string  `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"ipv4"`
	Port               *string  `json:"db_port" form:"db_port" example:"3306" valid:"range(1|65535)"`
	Password           *string  `json:"db_password" form:"db_password" example:"123456"`
	Desc               *string  `json:"desc" example:"this is a test instance" valid:"-"`
	WorkflowTemplateId *string  `json:"workflow_template_name" form:"workflow_template_name"`
	RuleTemplates      []string `json:"rule_template_name_list" form:"rule_template_name_list" example:"1" valid:"-"`
	Roles              []string `json:"role_name_list" form:"role_name_list"`
}

// @Summary 更新实例
// @Description update instance
// @Id updateInstanceV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @param instance body v1.UpdateInstanceReqV1 true "update instance request"
// @Success 200 {object} controller.BaseRes
// @router /v1/instances/{instance_name}/ [patch]
func UpdateInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("instance is not exist")))
	}

	req := new(UpdateInstanceReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	updateMap := map[string]string{}

	if req.Desc != nil {
		updateMap["desc"] = *req.Desc
	}
	if req.Host != nil {
		updateMap["host"] = *req.Host
	}
	if req.Port != nil {
		updateMap["port"] = *req.Port
	}
	if req.User != nil {
		updateMap["user"] = *req.User
	}
	if req.Password != nil {
		password, err := util.AesEncrypt(*req.Password)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		updateMap["password"] = password
	}

	if req.RuleTemplates != nil {
		ruleTemplates, err := s.GetAndCheckRuleTemplateExist(req.RuleTemplates)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		err = s.UpdateInstanceRuleTemplates(instance, ruleTemplates...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}
	if req.Roles != nil {
		roles, err := s.GetAndCheckRoleExist(req.Roles)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		err = s.UpdateInstanceRoles(instance, roles...)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	err = s.UpdateInstanceById(instance.ID, updateMap)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	go server.GetSqled().UpdateAndGetInstanceStatus(log.NewEntry(), instance)
	return c.JSON(200, controller.NewBaseReq(nil))
}

type GetInstancesReqV1 struct {
	FilterUserName string `json:"filter_rule_template_name" query:"filter_rule_template_name"`
	FilterRoleName string `json:"filter_role_name" query:"filter_role_name"`
	FilterDBHost   string `json:"filter_db_host" query:"filter_db_host" valid:"ipv4"`
	FilterDBPort   string `json:"filter_db_port" query:"filter_db_port" valid:"range(1|65535)"`
	FilterDBUser   string `json:"filter_db_user" query:"filter_db_user"`
	PageIndex      uint32 `json:"page_index" query:"page_index" valid:"required,int"`
	PageSize       uint32 `json:"page_size" query:"page_size" valid:"required,int"`
}

type GetInstancesResV1 struct {
	controller.BaseRes
	Data      []InstanceResV1 `json:"data"`
	TotalNums uint64          `json:"total_nums"`
}

//// @Summary 获取实例信息列表
//// @Description get instance info list
//// @Id getInstanceListV1
//// @Tags instance
//// @Security ApiKeyAuth
//// @Param filter_rule_template_name query string false "filter rule template name"
//// @Param filter_role_name query string false "filter role name"
//// @Param filter_db_host query string false "filter db host"
//// @Param filter_db_port query string false "filter db port"
//// @Param filter_db_user query string false "filter db user"
//// @Success 200 {object} v1.GetInstancesResV1
//// @router /instances [get]
//func GetInstances(c echo.Context) error {
//	s := model.GetStorage()
//	util.DebugPause("pause between get storage handle and query storage")
//	databases, err := s.GetInstances()
//	if err != nil {
//		return controller.JSONBaseErrorReq(c, err)
//	}
//	return c.JSON(http.StatusOK, &GetInstancesResV1{
//		BaseRes: controller.NewBaseReq(nil),
//		Data:    nil,
//	})
//}

type GetInstanceConnectableResV1 struct {
	controller.BaseRes
	Data InstanceConnectableResV1 `json:"data"`
}

type InstanceConnectableResV1 struct {
	IsInstanceConnectable bool `json:"is_instance_connectable"`
}

// @Summary 实例连通性测试（实例提交后）
// @Description test instance db connection
// @Id checkInstanceIsConnectableByNameV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceConnectableResV1
// @router /v1/instances/{instance_name}/connection [get]
func CheckInstanceIsConnectableByName(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("instance is not exist")))
	}
	isInstanceConnectable := true
	if err := executor.Ping(log.NewEntry(), instance); err != nil {
		isInstanceConnectable = false
	}
	return c.JSON(200, GetInstanceConnectableResV1{
		BaseRes: controller.NewBaseReq(err),
		Data: InstanceConnectableResV1{
			IsInstanceConnectable: isInstanceConnectable,
		},
	})
}

type GetInstanceConnectableReqV1 struct {
	User     string `json:"user" form:"db_user" example:"root"`
	Host     string `json:"host" form:"db_host" example:"10.10.10.10"`
	Port     string `json:"port" form:"db_port" example:"3306"`
	Password string `json:"password" form:"db_password" example:"123456"`
}

// @Summary 实例连通性测试（实例提交前）
// @Description test instance db connection 注：可直接提交创建实例接口的body，该接口的json 内容是创建实例的 json 的子集
// @Accept json
// @Id checkInstanceIsConnectableV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance body v1.GetInstanceConnectableReqV1 true "instance info"
// @Success 200 {object} v1.GetInstanceConnectableResV1
// @router /v1/instance_connection [post]
func CheckInstanceIsConnectable(c echo.Context) error {
	req := new(GetInstanceConnectableReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	instance := &model.Instance{
		DbType:   model.DB_TYPE_MYSQL,
		User:     req.User,
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
	}
	isInstanceConnectable := true
	if err := executor.Ping(log.NewEntry(), instance); err != nil {
		isInstanceConnectable = false
	}
	return c.JSON(200, GetInstanceConnectableResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: InstanceConnectableResV1{
			IsInstanceConnectable: isInstanceConnectable,
		},
	})
}

type GetInstanceSchemaResV1 struct {
	controller.BaseRes
	Data InstanceSchemaResV1 `json:"data"`
}

type InstanceSchemaResV1 struct {
	Schemas []string `json:"schema_name_list"`
}

// @Summary 实例 Schema 列表
// @Description instance schema list
// @Id getInstanceSchemasV1
// @Tags instance
// @Security ApiKeyAuth
// @Param instance_name path string true "instance name"
// @Success 200 {object} v1.GetInstanceSchemaResV1
// @router /v1/instances/{instance_name}/schemas [get]
func GetInstanceSchemas(c echo.Context) error {
	s := model.GetStorage()
	instanceName := c.Param("instance_name")
	instance, exist, err := s.GetInstanceByName(instanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DATA_NOT_EXIST, fmt.Errorf("instance is not exist")))
	}
	status, err := server.GetSqled().UpdateAndGetInstanceStatus(log.NewEntry(), instance)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(200, &GetInstanceSchemaResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: InstanceSchemaResV1{
			Schemas: status.Schemas,
		},
	})
}

type InstanceTipResV1 struct {
	Name string `json:"instance_name"`
}

type GetInstanceTipsResV1 struct {
	controller.BaseRes
	Data []InstanceTipResV1 `json:"data"`
}

// @Summary 获取实例提示列表
// @Description get instance tip list
// @Tags instance
// @Id getInstanceTipListV1
// @Security ApiKeyAuth
// @Success 200 {object} v1.GetInstanceTipsResV1
// @router /v1/instance_tips [get]
func GetInstanceTips(c echo.Context) error {
	s := model.GetStorage()
	roles, err := s.GetAllInstanceTip()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	instanceTipsResV1 := make([]InstanceTipResV1, 0, len(roles))

	for _, role := range roles {
		instanceTipRes := InstanceTipResV1{
			Name: role.Name,
		}
		instanceTipsResV1 = append(instanceTipsResV1, instanceTipRes)
	}
	return c.JSON(http.StatusOK, &GetInstanceTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    instanceTipsResV1,
	})
}

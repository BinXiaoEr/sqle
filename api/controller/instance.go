package controller

import (
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"sqle/api/server"
	"sqle/errors"
	"sqle/executor"
	"sqle/model"
	"strings"
)

type CreateInstanceReq struct {
	Name string `json:"name" form:"name" example:"test"`
	// mysql, mycat, sqlserver
	DbType   string `json:"db_type" form:"type" example:"mysql"`
	User     string `json:"user" form:"user" example:"root"`
	Host     string `json:"host" form:"host" example:"10.10.10.10"`
	Port     string `json:"port" form:"port" example:"3306"`
	Password string `json:"password" form:"password" example:"123456"`
	Desc     string `json:"desc" form:"desc" example:"this is a test instance"`
	// this a list for rule template name
	RuleTemplates []string `json:"rule_template_name_list" form:"rule_template_name_list" example:"template_1"`
	// mycat_config is required if db_type is "mycat"
	MycatConfig *model.MycatConfig `json:"mycat_config" form:"mycat_config"`
}

type InstanceRes struct {
	BaseRes
	Data model.InstanceDetail `json:"data"`
}

// @Summary 添加实例
// @Description create a instance
// @Description create a instance
// @Accept json
// @Accept json
// @Param instance body controller.CreateInstanceReq true "add instance"
// @Success 200 {object} controller.InstanceRes
// @router /instances [post]
func CreateInst(c echo.Context) error {
	s := model.GetStorage()
	req := new(CreateInstanceReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	_, exist, err := s.GetInstByName(req.Name)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if exist {
		err := fmt.Errorf("inst is exist")
		return c.JSON(200, NewBaseReq(errors.New(errors.INSTANCE_EXIST, err)))
	}

	instance := &model.Instance{
		Name:        req.Name,
		DbType:      req.DbType,
		User:        req.User,
		Host:        req.Host,
		Port:        req.Port,
		Password:    req.Password,
		Desc:        req.Desc,
		MycatConfig: req.MycatConfig,
	}
	err = instance.MarshalMycatConfig()
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}

	notExistTs := []string{}
	for _, tplName := range req.RuleTemplates {
		t, exist, err := s.GetTemplateByName(tplName)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
		if !exist {
			notExistTs = append(notExistTs, tplName)
		}
		instance.RuleTemplates = append(instance.RuleTemplates, t)
	}

	if len(notExistTs) > 0 {
		err := fmt.Errorf("rule_template %s not exist", strings.Join(notExistTs, ", "))
		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST, err)))
	}

	err = s.Save(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	go server.GetSqled().UpdateAndGetInstanceStatus(instance)

	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(nil),
		Data:    instance.Detail(),
	})
}

// @Summary 获取实例信息
// @Description get instance db
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.InstanceRes
// @router /instances/{instance_id}/ [get]
func GetInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	return c.JSON(200, &InstanceRes{
		BaseRes: NewBaseReq(nil),
		Data:    instance.Detail(),
	})
}

// @Summary 删除实例
// @Description delete instance db
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.BaseRes
// @router /instances/{instance_id}/ [delete]
func DeleteInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	err = s.Delete(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	server.GetSqled().DeleteInstanceStatus(instance)
	return c.JSON(200, NewBaseReq(nil))
}

// @Summary 更新实例
// @Description update instance db
// @Param instance_id path string true "Instance ID"
// @param instance body controller.CreateInstanceReq true "update instance"
// @Success 200 {object} controller.BaseRes
// @router /instances/{instance_id}/ [put]
func UpdateInstance(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	req := new(CreateInstanceReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, INSTANCE_NOT_EXIST_ERROR)
	}
	if instance.Name != req.Name {
		_, exist, err := s.GetInstByName(req.Name)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
		if exist {
			return c.JSON(200, NewBaseReq(errors.New(errors.INSTANCE_EXIST,
				fmt.Errorf("instance name is exist"))))
		}
	}

	instance.Name = req.Name
	instance.Desc = req.Desc
	instance.DbType = req.DbType
	instance.Host = req.Host
	instance.Port = req.Port
	instance.User = req.User
	instance.Password = req.Password
	instance.RuleTemplates = nil

	notExistTs := []string{}
	ruleTemplates := []model.RuleTemplate{}
	for _, tplName := range req.RuleTemplates {
		t, exist, err := s.GetTemplateByName(tplName)
		if err != nil {
			return c.JSON(200, NewBaseReq(err))
		}
		if !exist {
			notExistTs = append(notExistTs, tplName)
		}
		ruleTemplates = append(ruleTemplates, t)
	}

	if len(notExistTs) > 0 {
		err := fmt.Errorf("rule_template %s not exist", strings.Join(notExistTs, ", "))
		return c.JSON(200, NewBaseReq(errors.New(errors.RULE_TEMPLATE_NOT_EXIST, err)))
	}

	err = s.Save(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	err = s.UpdateInstRuleTemplate(instance, ruleTemplates...)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}

	go server.GetSqled().UpdateAndGetInstanceStatus(instance)

	return c.JSON(200, NewBaseReq(nil))
}

type GetAllInstReq struct {
	BaseRes
	Data []*model.Instance `json:"data"`
}

// @Summary 实例列表
// @Description get all instances
// @Success 200 {object} controller.GetAllInstReq
// @router /instances [get]
func GetInsts(c echo.Context) error {
	s := model.GetStorage()
	databases, err := s.GetInstances()
	if err != nil {
		return c.JSON(http.StatusOK, NewBaseReq(err))
	}
	return c.JSON(http.StatusOK, &GetAllInstReq{
		BaseRes: NewBaseReq(nil),
		Data:    databases,
	})
}

type PingInstRes struct {
	BaseRes
	// true: ping success; false: ping failed
	Data bool `json:"data"`
}

// @Summary 实例连通性测试
// @Description test instance db connection
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.PingInstRes
// @router /instances/{instance_id}/connection [get]
func PingInst(c echo.Context) error {
	s := model.GetStorage()
	instId := c.Param("instance_id")
	inst, exist, err := s.GetInstById(instId)
	if err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(err),
			Data:    false,
		})
	}
	if !exist {
		return c.JSON(200, PingInstRes{
			BaseRes: INSTANCE_NOT_EXIST_ERROR,
			Data:    false,
		})
	}
	if err := executor.Ping(inst); err != nil {
		return c.JSON(200, PingInstRes{
			BaseRes: NewBaseReq(nil),
			Data:    false,
		})
	}
	return c.JSON(200, PingInstRes{
		BaseRes: NewBaseReq(nil),
		Data:    true,
	})
}

type GetSchemaRes struct {
	BaseRes
	Data []string `json:"data" example:"db1"`
}

// @Summary 实例 Schema 列表
// @Description instance schema list
// @Param instance_id path string true "Instance ID"
// @Success 200 {object} controller.GetSchemaRes
// @router /instances/{instance_id}/schemas [get]
func GetInstSchemas(c echo.Context) error {
	s := model.GetStorage()
	instanceId := c.Param("instance_id")
	instance, exist, err := s.GetInstById(instanceId)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	if !exist {
		return c.JSON(200, INSTANCE_NOT_EXIST_ERROR)
	}
	status, err := server.GetSqled().UpdateAndGetInstanceStatus(instance)
	if err != nil {
		return c.JSON(200, NewBaseReq(err))
	}
	return c.JSON(200, &GetSchemaRes{
		BaseRes: NewBaseReq(nil),
		Data:    status.Schemas,
	})
}

type GetAllSchemasRes struct {
	BaseRes
	Data []server.InstanceStatus `json:"data"`
}

// @Summary 所有实例的 Schema 列表
// @Description all schema list
// @Success 200 {object} controller.GetAllSchemasRes
// @router /schemas [get]
func GetAllSchemas(c echo.Context) error {
	return c.JSON(200, &GetAllSchemasRes{
		BaseRes: NewBaseReq(nil),
		Data:    server.GetSqled().GetAllInstanceStatus(),
	})
}

// @Summary 手动刷新 Schema 列表
// @Description update schema list
// @Success 200 {object} controller.BaseRes
// @router /schemas/manual_update [post]
func ManualUpdateAllSchemas(c echo.Context) error {
	go server.GetSqled().UpdateAllInstanceStatus()
	return c.JSON(200, NewBaseReq(nil))
}

type LoadMycatConfigRes struct {
	BaseRes
	Data CreateInstanceReq `json:"data"`
}

// @Summary 解析 mycat 配置
// @Description unmarshal mycat config file
// @Accept mpfd
// @Param server_xml formData file true "server.xml"
// @Param schema_xml formData file true "schema.xml"
// @Param rule_xml formData file true "rule.xml"
// @Success 200 {object} controller.CreateInstanceReq
// @router /instance/load_mycat_config [post]
func UploadMycatConfig(c echo.Context) error {
	_, server, err := readFileToByte(c, "server_xml")
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	_, schemas, err := readFileToByte(c, "schema_xml")
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	_, rules, err := readFileToByte(c, "rule_xml")
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	var serverXML = &model.ServerXML{}
	var schemasXML = &model.SchemasXML{}
	var rulesXML = &model.RulesXML{}
	var instance *model.Instance

	err = xml.Unmarshal(server, serverXML)
	if err != nil {
		goto ERROR
	}

	err = xml.Unmarshal(schemas, schemasXML)
	if err != nil {
		goto ERROR
	}
	err = xml.Unmarshal(rules, rulesXML)
	if err != nil {
		goto ERROR
	}
	instance, err = model.LoadMycatServerFromXML(serverXML, schemasXML, rulesXML)
	if err != nil {
		goto ERROR
	}
	return c.JSON(http.StatusOK, LoadMycatConfigRes{
		BaseRes: NewBaseReq(nil),
		Data: CreateInstanceReq{
			DbType:        instance.DbType,
			User:          instance.User,
			RuleTemplates: []string{},
			MycatConfig:   instance.MycatConfig,
		},
	})

ERROR:
	return c.JSON(http.StatusOK, errors.New(errors.READ_UPLOAD_FILE_ERROR, err))
}

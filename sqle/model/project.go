package model

import (
	"database/sql"
	"fmt"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

const ProjectIdForGlobalRuleTemplate = 0

type Project struct {
	Model
	Name string
	Desc string

	CreateUserId uint  `gorm:"not null"`
	CreateUser   *User `gorm:"foreignkey:CreateUserId"`

	Managers   []*User      `gorm:"many2many:project_manager;"`
	Members    []*User      `gorm:"many2many:project_user;"`
	UserGroups []*UserGroup `gorm:"many2many:project_user_group;"`

	Workflows     []*Workflow     `gorm:"foreignkey:ProjectId"`
	AuditPlans    []*AuditPlan    `gorm:"foreignkey:ProjectId"`
	Instances     []*Instance     `gorm:"foreignkey:ProjectId"`
	SqlWhitelist  []*SqlWhitelist `gorm:"foreignkey:ProjectId"`
	RuleTemplates []*RuleTemplate `gorm:"foreignkey:ProjectId"`

	WorkflowTemplateId uint              `gorm:"not null"`
	WorkflowTemplate   *WorkflowTemplate `gorm:"foreignkey:WorkflowTemplateId"`
}

// IsProjectExist 用于判断当前是否存在项目, 而非某个项目是否存在
func (s *Storage) IsProjectExist() (bool, error) {
	var count uint
	err := s.db.Model(&Project{}).Count(&count).Error
	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) CreateProject(name string, desc string, createUserID uint) error {
	wt := &WorkflowTemplate{
		Name:                          fmt.Sprintf("%v-WorkflowTemplate", name),
		Desc:                          fmt.Sprintf("%v 默认模板", name),
		AllowSubmitWhenLessAuditLevel: string(driver.RuleLevelWarn),
		Steps: []*WorkflowStepTemplate{
			{
				Number: 1,
				Typ:    WorkflowStepTypeSQLReview,
				ApprovedByAuthorized: sql.NullBool{
					Bool:  true,
					Valid: true,
				},
			},
			{
				Number: 2,
				Typ:    WorkflowStepTypeSQLExecute,
				Users:  []*User{{Model: Model{ID: createUserID}}},
			},
		},
	}

	return s.TxExec(func(tx *sql.Tx) error {
		templateID, err := saveWorkflowTemplate(wt, tx)
		if err != nil {
			return err
		}

		result, err := tx.Exec("INSERT INTO projects (`name`, `desc`, `create_user_id`,`workflow_template_id`) values (?, ?, ?,?)", name, desc, createUserID, templateID)
		if err != nil {
			return err
		}
		projectID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO project_manager (`project_id`, `user_id`) VALUES (?, ?)", projectID, createUserID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO project_user (`project_id`, `user_id`) VALUES (?, ?)", projectID, createUserID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) CheckUserCanUpdateProject(projectName string, userID uint) (bool, error) {
	user, exist, err := s.GetUserByID(userID)
	if err != nil || !exist {
		return false, err
	}

	if user.Name == DefaultAdminUser {
		return true, nil
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil || !exist {
		return false, err
	}

	for _, manager := range project.Managers {
		if manager.ID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) UpdateProjectInfoByID(projectName string, attr map[string]interface{}) error {
	err := s.db.Table("projects").Where("name = ?", projectName).Update(attr).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetProjectByID(projectID uint) (*Project, bool, error) {
	p := &Project{}
	err := s.db.Model(&Project{}).Where("id = ?", projectID).Find(p).Error
	if err == gorm.ErrRecordNotFound {
		return p, false, nil
	}
	return p, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsProjectManager(userName string, projectName string) (bool, error) {
	var count uint

	err := s.db.Table("project_manager").
		Joins("projects ON projects.id = project_manager.project_id").
		Joins("users ON project_manager.user_id = users.id").
		Where("users.login_name = ?", userName).
		Where("users.stats = 0").
		Where("projects.name = ?", projectName).
		Count(&count).Error

	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s Storage) GetProjectByName(projectName string) (*Project, bool, error) {
	p := &Project{}
	err := s.db.Preload("CreateUser").Preload("Managers").Where("name = ?", projectName).First(p).Error
	if err == gorm.ErrRecordNotFound {
		return p, false, nil
	}
	return p, true, errors.New(errors.ConnectStorageError, err)
}

func (s Storage) GetProjectTips(userName string) ([]*Project, error) {
	p := []*Project{}
	query := s.db.Table("projects").Select("name")

	var err error
	if userName != DefaultAdminUser {
		err = query.Joins("JOIN project_user on project_user.project_id = projects.id").
			Joins("JOIN users on users.id = project_user.user_id").
			Joins("JOIN project_user_group on project_user_group.project_id = projects.id").
			Joins("JOIN project_user_group on project_user_group.project_id = projects.id").
			Joins("JOIN user_group_users on project_user_group.user_group_id = user_group_users.user_group_id").
			Joins("RIGHT JOIN users as u on u.id = user_group_users.user_id").
			Where("users.stat = 0").Where("u.stat = 0").
			Where("users.login_name = ? OR u.login_name = ?", userName, userName).Find(&p).Error
	} else {
		err = query.Find(&p).Error
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return p, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsProjectMember(userName, projectName string) (bool, error) {
	query := `
SELECT EXISTS(
SELECT users.login_name 
FROM users
JOIN project_user on project_user.user_id = users.id
JOIN projects on project_user.project_id = projects.id
JOIN user_group_users on users.id = user_group_users.user_id 
JOIN project_user_group on user_group_users.user_group_id = project_user_group.user_group_id
JOIN projects as p on project_user_group.project_id = p.id
WHERE users.stat = 0
AND( 
	projects.name = ?
OR
	p.name = ?
)) AS exist
`
	var exist struct {
		Exist bool `json:"exist"`
	}
	err := s.db.Raw(query, userName, projectName).Find(&exist).Error
	return exist.Exist, errors.New(errors.ConnectStorageError, err)
}

type ProjectAndInstance struct {
	InstanceName string `json:"instance_name"`
	ProjectName  string `json:"project_name"`
}

func (s *Storage) GetProjectNamesByInstanceIds(instanceIds []uint) (map[uint] /*instance id*/ ProjectAndInstance, error) {
	instanceIds = utils.RemoveDuplicateUint(instanceIds)
	type record struct {
		InstanceId   uint   `json:"instance_id"`
		InstanceName string `json:"instance_name"`
		ProjectName  string `json:"project_name"`
	}
	records := []record{}
	err := s.db.Table("instances").
		Joins("LEFT JOIN projects ON projects.id = instances.project_id").
		Select("instances.id AS instance_id, instances.name AS instance_name, projects.name AS project_name").
		Where("instances.id IN (?)", instanceIds).
		Find(&records).Error

	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	res := make(map[uint]ProjectAndInstance, len(records))
	for _, r := range records {
		res[r.InstanceId] = ProjectAndInstance{
			InstanceName: r.InstanceName,
			ProjectName:  r.ProjectName,
		}
	}

	return res, nil
}

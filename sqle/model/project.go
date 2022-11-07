package model

import (
	"github.com/actiontech/sqle/sqle/errors"
)

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

func (s *Storage) IsProjectManager(userID uint, projectName string) (bool, error) {
	var count uint

	err := s.db.Table("project_manager").
		Joins("projects ON projects.id = project_manager.project_id").
		Where("project_manager.user_id = ?", userID).
		Where("projects.name = ?", projectName).
		Count(&count).Error

	return count > 0, errors.New(errors.ConnectStorageError, err)
}

package sqlmanage

import (
	"github.com/actiontech/sqle/sqle/model"
)

type SqlManager interface {
	SyncSqlManager(source string) error
	UpdateSqlManageRecord(sourceId, source string) error
}

type SyncFromSqlAuditRecord struct {
	Task             *model.Task
	SqlFpMap         map[string]string
	ProjectId        string
	SqlAuditRecordID string
}

func NewSyncFromSqlAudit(task *model.Task, fpMap map[string]string, projectID string, sqlAuditID string) SqlManager {
	return &SyncFromSqlAuditRecord{
		Task:             task,
		ProjectId:        projectID,
		SqlAuditRecordID: sqlAuditID,
		SqlFpMap:         fpMap,
	}
}

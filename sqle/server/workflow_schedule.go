package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/notification"

	"github.com/sirupsen/logrus"
)

func (s *Sqled) workflowScheduleLoop() {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	entry := log.NewEntry().WithField("type", "schedule_workflow")
	for {
		select {
		case <-s.exit:
			return
		case <-tick.C:
			s.WorkflowSchedule(entry)
		}
	}
}

func (s *Sqled) WorkflowSchedule(entry *logrus.Entry) {
	st := model.GetStorage()
	workflows, err := st.GetNeedScheduledWorkflows()
	if err != nil {
		entry.Errorf("get need scheduled workflows from storage error: %v", err)
		return
	}
	for _, workflow := range workflows {
		w, exist, err := st.GetWorkflowDetailById(strconv.Itoa(int(workflow.ID)))
		if err != nil {
			entry.Errorf("get workflow from storage error: %v", err)
			return
		}
		if !exist {
			entry.Errorf("workflow %s not found", workflow.Subject)
			return
		}

		currentStep := w.CurrentStep()
		if currentStep == nil {
			entry.Errorf("workflow %s not found", w.Subject)
			return
		}
		if currentStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
			entry.Errorf("workflow %s need to be approved first", w.Subject)
			return
		}

		entry.Infof("start to execute scheduled workflow %s", w.Subject)
		err = ExecuteWorkflow(w, w.Record.ScheduleUserId)
		if err != nil {
			entry.Errorf("execute scheduled workflow %s error: %v", w.Subject, err)
		} else {
			entry.Infof("execute scheduled workflow %s success", w.Subject)
		}
	}
}

func ExecuteWorkflow(workflow *model.Workflow, needExecTaskIds map[uint]struct{}, userId uint) error {
	s := model.GetStorage()

	// get task and check connection before to execute it.
	for taskId := range needExecTaskIds {
		taskId := fmt.Sprintf("%d", taskId)
		task, exist, err := s.GetTaskDetailById(taskId)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New(errors.DataNotExist, fmt.Errorf("task is not exist. taskID=%v", taskId))
		}
		if task.Instance == nil {
			return errors.New(errors.DataNotExist, fmt.Errorf("instance is not exist"))
		}

		// if instance is not connectable, exec sql must be failed;
		// commit action unable to retry, so don't to exec it.
		if err = common.CheckInstanceIsConnectable(task.Instance); err != nil {
			return errors.New(errors.ConnectRemoteDatabaseError, err)
		}
	}

	currentStep := workflow.CurrentStep()
	if currentStep == nil {
		return fmt.Errorf("workflow current step not found")
	}

	// update workflow
	waitForExecTasksCount, err := s.GetWaitExecInstancesCountByWorkflowId(workflow.ID)
	if err != nil {
		return fmt.Errorf("get count of tasks failed: %v", err)
	}
	for i, inst := range workflow.Record.InstanceRecords {
		if _, ok := needExecTaskIds[inst.TaskId]; ok {
			workflow.Record.InstanceRecords[i].IsSQLExecuted = true
		}
	}

	// 只有当所有数据源都上线时，current step状态才改为"approved"
	if waitForExecTasksCount == len(needExecTaskIds) {
		currentStep.State = model.WorkflowStepStateApprove
		workflow.Record.Status = model.WorkflowStatusFinish
		workflow.Record.CurrentWorkflowStepId = 0
	}

	// todo issue832 上线人暂时先保存最后一个上线人，后续可能需要保存到每个数据源上
	now := time.Now()
	currentStep.OperateAt = &now
	currentStep.OperationUserId = userId

	err = s.UpdateWorkflowStatus(workflow, currentStep,workflow.Record.InstanceRecords)
	if err != nil {
		return err
	}

	for taskId := range needExecTaskIds {
		id := taskId
		go func() {
			sqledServer := GetSqled()
			task, err := sqledServer.AddTaskWaitResult(strconv.Itoa(int(id)), ActionTypeExecute)
			if err != nil || task.Status == model.TaskStatusExecuteFailed {
				go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeExecuteFail)
			} else {
				go notification.NotifyWorkflow(fmt.Sprintf("%v", workflow.ID), notification.WorkflowNotifyTypeExecuteSuccess)
			}
		}()
	}
	return nil
}

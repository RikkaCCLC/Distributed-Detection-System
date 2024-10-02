package rpc

import (
	"github.com/toolkits/pkg/logger"
	"open-devops/src/models"
	"open-devops/src/modules/server/task"
)

// agent上报的单个任务结果
type ReportTask struct {
	Id     int64
	Clock  int64
	Status string
	Stdout string
	Stderr string
}

// agent的上报请求
type TaskReportRequest struct {
	AgentIp     string       // agent的ip用来获取新的任务
	ReportTasks []ReportTask // 上次任务的结果
}

// 下发给agent的任务
type TaskReportResponse struct {
	Message     string
	AssignTasks []*models.TaskMeta
}

func (t *Server) TaskReport(args TaskReportRequest, reply *TaskReportResponse) error {
	toMarkDoneIds := make(map[int64]struct{})
	if len(args.ReportTasks) > 0 {
		//	处理返回
		// 将task返回结果入库
		//
		// 首先要处理agent上报的任务结果，就是将ReportTasks翻译成models.TaskResult记录到库中

		for _, x := range args.ReportTasks {
			tRes := models.TaskResult{
				Id:     0,
				TaskId: x.Id,
				Host:   args.AgentIp,
				Status: x.Status,
				Stdout: x.Stdout,
				Stderr: x.Stderr,
			}
			err, added := tRes.Save()
			if err != nil {
				logger.Errorf("[TaskResult.SaveError][agent.ip:%+v][tRes:%+v][error:%+v]",
					args.AgentIp,
					tRes,
					err,
				)
			}
			if added {
				logger.Infof("[TaskResult.Save.Success][agent.ip:%+v][tRes:%+v]",
					args.AgentIp,
					tRes,
				)
				toMarkDoneIds[x.Id] = struct{}{}
			}

		}
		// 将job标记为已处理 ,这样下次就不会再分配这个任务了
		for id, _ := range toMarkDoneIds {
			err := models.MarkTaskMetaDone(id)
			if err != nil {
				logger.Errorf("[TaskMeta.MarkTaskMetaDone][agent.ip:%+v][id:%+v][error:%+v]",
					args.AgentIp,
					id,
					err,
				)
			} else {
				logger.Infof("[TaskMeta.MarkTaskMetaDone][agent.ip:%+v][id:%+v]",
					args.AgentIp,
					id,
				)
			}
		}

	}
	reply.AssignTasks = task.TaskCaches.GetTasksByIp(args.AgentIp)

	return nil
}

package rpc

import (
	"context"
	"github.com/go-kit/log/level"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
	"open-devops/src/modules/agent/taskjob"
	serverRpc "open-devops/src/modules/server/rpc"
	"time"
)

func TickerTaskReport(cli *RpcCli, ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)

	localIp := common.GetLocalIp()
	cli.DoTaskReport(localIp)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("TickerLogJobSync.receive_quit_signal_and_quit")
			return nil
		case <-ticker.C:
			cli.DoTaskReport(localIp)
		}
	}

}

func (r *RpcCli) DoTaskReport(localIp string) {
	// 构造TaskReport的rpc请求，其中ReportTasks来自本地map的任务收集任务
	req := serverRpc.TaskReportRequest{
		AgentIp:     localIp,
		ReportTasks: taskjob.Locals.ReportTasks(),
	}
	var resp serverRpc.TaskReportResponse
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return
	}
	err = r.Cli.Call("Server.TaskReport", req, &resp)
	if err != nil {
		r.CloseCli()
		level.Error(r.logger).Log("msg", "Server.TaskReport.error", "serverAddr", r.ServerAddr, "err", err)
		return
	}
	// 遍历rpc的结果，分配任务
	if resp.AssignTasks != nil {

		count := len(resp.AssignTasks)
		for i := 0; i < count; i++ {
			at := resp.AssignTasks[i]

			taskjob.Locals.AssignTask(at)
		}
	}

}

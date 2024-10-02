package rpc

import (
	"github.com/go-kit/log/level"
	"open-devops/src/models"
)

func (r *RpcCli) HostInfoReport(info models.AgentCollectInfo) {
	var msg string
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return
	}
	err = r.Cli.Call("Server.HostInfoReport", info, &msg)
	if err != nil {
		r.CloseCli()
		level.Error(r.logger).Log("msg", "Server.HostInfoReport.error", "serverAddr", r.ServerAddr, "err", err)
		return
	}

}

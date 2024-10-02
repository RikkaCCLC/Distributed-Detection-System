package rpc

import (
	"github.com/go-kit/log/level"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/models"
)

func (r *RpcCli) LogJobSync(hostname string) []*models.LogStrategy {
	var res []*models.LogStrategy
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return nil
	}
	err = r.Cli.Call("Server.LogJobSync", hostname, &res)
	if err != nil {
		r.CloseCli()
		level.Error(r.logger).Log("msg", "Server.LogJobSync.error", "serverAddr", r.ServerAddr, "err", err)
		return nil
	}

	logger.Infof("LogJobSync.res:%+v", res)
	return res
}

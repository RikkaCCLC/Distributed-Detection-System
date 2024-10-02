package rpc

import (
	"github.com/go-kit/log/level"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
)

func (r *RpcCli) ProberTargetSync(req common.ProberTargetsGetRequest) *common.ProberTargetsGetResponse {
	var res *common.ProberTargetsGetResponse
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return nil
	}
	err = r.Cli.Call("Server.GetProberTargets", req, &res)
	if err != nil {
		r.CloseCli()
		level.Error(r.logger).Log("msg", "Server.ProberTargetSync.error", "serverAddr", r.ServerAddr, "err", err)
		return nil
	}

	logger.Infof("ProberTargetSync.res:%+v", res)
	return res
}

func (r *RpcCli) PushProberResults(req common.ProberResultPushRequest) *common.ProberResultPushResponse {
	var res *common.ProberResultPushResponse
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return nil
	}
	err = r.Cli.Call("Server.PushProberResults", req, &res)
	if err != nil {
		r.CloseCli()
		level.Error(r.logger).Log("msg", "Server.PushProberResults.error", "serverAddr", r.ServerAddr, "err", err)
		return nil
	}

	logger.Infof("PushProberResults.res:%+v req:%+v", res, req.ProberResults[0])
	return res
}

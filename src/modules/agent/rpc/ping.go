package rpc

import "github.com/go-kit/log/level"

func (r *RpcCli) Ping() {
	var msg string
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return
	}
	// call 调用rpc
	err = r.Cli.Call("Server.Ping", "agent01", &msg)
	if err != nil {
		level.Error(r.logger).Log("msg", "Server.Ping.error", "serverAddr", r.ServerAddr, "err", err)
		return
	}
	level.Info(r.logger).Log("msg", "Server.Ping.success", "serverAddr", r.ServerAddr, "msg", msg)
}

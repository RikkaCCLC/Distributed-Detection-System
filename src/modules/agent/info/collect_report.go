package info

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"open-devops/src/common"
	"open-devops/src/models"
	"open-devops/src/modules/agent/rpc"
	"time"
)

func TickerInfoCollectAndReport(cli *rpc.RpcCli, ctx context.Context, logger log.Logger) error {
	ticker := time.NewTicker(5 * time.Second)

	level.Info(logger).Log("msg", "TickerInfoCollectAndReport.start")
	CollectBaseInfo(cli, logger)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			level.Info(logger).Log("msg", "receive_quit_signal_and_quit")
			return nil
		case <-ticker.C:
			CollectBaseInfo(cli, logger)
		}
	}

}

func CollectBaseInfo(cli *rpc.RpcCli, logger log.Logger) {

	var (
		err  error
		sn   string
		cpu  string
		mem  string
		disk string
	)
	snShellCloud := `curl -s http://169.254.169.254/a/meta-data/instance-id`
	snShellHost := `dmidecode -s system-serial-number |tail -n 1|tr -d "\n"`

	cpuShell := `cat /proc/cpuinfo |grep processor |wc -l| tr -d "\n"`
	memShell := `cat /proc/meminfo |grep MemTotal |awk '{printf "%d",$2/1024/1024}'`
	diskShell := `df -m  |grep '/dev/' |grep -v '/var/lib' |grep -v tmpfs |awk '{sum +=$2};END{printf "%d",sum/1024}'`

	sn, err = common.ShellCommand(snShellCloud)
	if err != nil || sn == "" {
		sn, err = common.ShellCommand(snShellHost)
	}
	level.Info(logger).Log("msg", "CollectBaseInfo", "sn", sn)

	cpu, err = common.ShellCommand(cpuShell)
	if err != nil {
		level.Error(logger).Log("msg", "cpuShell.error", "shell", cpuShell, "err", err)
	}
	level.Info(logger).Log("msg", "CollectBaseInfo", "cpu", cpu)

	mem, err = common.ShellCommand(memShell)
	if err != nil {
		level.Error(logger).Log("msg", "memShell.error", "shell", memShell, "err", err)
	}
	level.Info(logger).Log("msg", "CollectBaseInfo", "mem", mem)

	disk, err = common.ShellCommand(diskShell)
	if err != nil {
		level.Error(logger).Log("msg", "memShell.error", "diskShell", diskShell, "err", err)
	}
	level.Info(logger).Log("msg", "CollectBaseInfo", "disk", disk)

	ipAddr := common.GetLocalIp()
	hostName := common.GetHostName()

	hostObj := models.AgentCollectInfo{
		SN:       sn,
		CPU:      cpu,
		Mem:      mem,
		Disk:     disk,
		IpAddr:   ipAddr,
		HostName: hostName,
	}
	cli.HostInfoReport(hostObj)
}

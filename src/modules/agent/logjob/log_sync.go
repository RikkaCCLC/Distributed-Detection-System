package logjob

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/models"
	"open-devops/src/modules/agent/config"
	"open-devops/src/modules/agent/rpc"
	"time"
)

func TickerLogJobSync(cli *rpc.RpcCli, ctx context.Context, logJobsyncChan chan []*LogJob, localConfigJobs []*LogJob, metricsMap map[string]*prometheus.GaugeVec, hostName string) error {
	ticker := time.NewTicker(5 * time.Second)
	doLogJobSync(cli, logJobsyncChan, localConfigJobs, metricsMap, hostName)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("TickerLogJobSync.receive_quit_signal_and_quit")
			return nil
		case <-ticker.C:
			doLogJobSync(cli, logJobsyncChan, localConfigJobs, metricsMap, hostName)
		}
	}

}

func doLogJobSync(cli *rpc.RpcCli, logJobsyncChan chan []*LogJob, localConfigJobs []*LogJob, metricsMap map[string]*prometheus.GaugeVec, hostName string) {

	res := cli.LogJobSync(hostName)
	ls := []*models.LogStrategy{}
	//

	for _, i := range res {
		i := i
		m := map[string]string{}
		json.Unmarshal([]byte(i.TagJson), &m)
		i.Tags = m
		labels := []string{}
		for k := range i.Tags {
			labels = append(labels, k)
		}
		me := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: i.MetricName,
			Help: i.MetricHelp,
		}, labels)
		// 为了动态注册，防止重复注册，用map
		if _, loaded := metricsMap[i.MetricName]; !loaded {
			prometheus.MustRegister(me)
			metricsMap[i.MetricName] = me

		}

		ls = append(ls, i)

		logger.Infof("[doLogJobSync.rpc.res][num:%d][res:%+v]i.Tags:%+v", len(res), i, i.Tags)
	}

	newLs := config.SetLogRegs(ls)
	for _, i := range newLs {
		j := &LogJob{Stra: i}
		localConfigJobs = append(localConfigJobs, j)

	}

	logJobsyncChan <- localConfigJobs
}

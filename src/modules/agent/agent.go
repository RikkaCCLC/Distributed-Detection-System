package main

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	promlogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"open-devops/src/common"
	"open-devops/src/modules/agent/config"
	"open-devops/src/modules/agent/consumer"
	"open-devops/src/modules/agent/counter"
	"open-devops/src/modules/agent/info"
	"open-devops/src/modules/agent/logjob"
	"open-devops/src/modules/agent/metrics"
	"open-devops/src/modules/agent/rpc"
	"open-devops/src/modules/agent/taskjob"
	"open-devops/src/modules/agent/xprober"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	// 命令行解析
	app = kingpin.New(filepath.Base(os.Args[0]), "The open-devops-agent")
	// 指定配置文件
	configFile = app.Flag("config.file", "open-devops-agent configuration file path").Short('c').Default("open-devops-agent.yml").String()
)

func main() {
	// 版本信息
	app.Version(version.Print("open-devops-agent"))
	// 帮助信息
	app.HelpFlag.Short('h')

	promlogConfig := promlog.Config{}

	promlogflag.AddFlags(app, &promlogConfig)
	// 强制解析
	kingpin.MustParse(app.Parse(os.Args[1:]))
	// 设置logger
	var logger log.Logger
	logger = func(config *promlog.Config) log.Logger {
		var (
			l  log.Logger
			le level.Option
		)
		if config.Format.String() == "logfmt" {
			l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		} else {
			l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		}

		switch config.Level.String() {
		case "debug":
			le = level.AllowDebug()
		case "info":
			le = level.AllowInfo()
		case "warn":
			le = level.AllowWarn()
		case "error":
			le = level.AllowError()
		}
		l = level.NewFilter(l, le)
		l = log.With(l, "ts", log.TimestampFormat(
			func() time.Time { return time.Now().Local() },
			"2006-01-02 15:04:05.000 ",
		), "caller", log.DefaultCaller)
		return l
	}(&promlogConfig)

	level.Debug(logger).Log("debug.msg", "using config.file", "file.path", *configFile)

	sConfig, err := config.LoadFile(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "config.LoadFile Error,Exiting ...", "error", err)
		return
	}
	level.Info(logger).Log("msg", "load.config.success", "file.path", *configFile, "rpc_server_addr", sConfig.RpcServerAddr)
	// 初始化rpc client
	rpcCli := rpc.InitRpcCli(sConfig.RpcServerAddr, logger)

	// 这里开始 logJob
	// 创建metrics
	metricsMap := metrics.CreateMetrics(sConfig.LogStrategies)
	// 注册metrics
	for _, m := range metricsMap {
		prometheus.MustRegister(m)
	}

	// 统计指标的同步queue
	cq := make(chan *consumer.AnalysPoint, common.CounterQueueSize)
	// 统计指标的管理器
	pm := counter.NewPointCounterManager(cq, metricsMap)
	// 日志job管理器
	logJobManager := logjob.NewLogJobManager(cq)
	// 把配置文件中的logjob传入
	logJobsyncChan := make(chan []*logjob.LogJob, 1)
	localConfigJobs := make([]*logjob.LogJob, 0)
	for _, i := range sConfig.LogStrategies {
		i := i
		j := &logjob.LogJob{Stra: i}
		localConfigJobs = append(localConfigJobs, j)

	}
	logJobsyncChan <- localConfigJobs



	// 初始化任务缓存
	taskjob.InitLocals(sConfig.Job.MetaDir)

	// 初始化探测的缓存,LocalRegion 和LocalIp
	xprober.Init(logger, sConfig.Region)

	// 编排开始
	var g run.Group
	ctxAll, cancelAll := context.WithCancel(context.Background())
	{

		// 处理信号退出的handler
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancelC := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					level.Warn(logger).Log("msg", "Receive SIGTERM ,exiting gracefully....")
					cancelAll()
					return nil
				case <-cancelC:
					level.Warn(logger).Log("msg", "other cancel exiting")
					return nil
				}
			},
			func(err error) {
				close(cancelC)
			},
		)
	}
	{
		// 采集基础信息的
		if sConfig.EnableInfoCollectAndReport {

			g.Add(func() error {
				err := info.TickerInfoCollectAndReport(rpcCli, ctxAll, logger)
				if err != nil {
					level.Error(logger).Log("msg", "TickerInfoCollectAndReport.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}
	}
	if sConfig.EnableLogJob {
		{

			// logJobManager 增量同步策略，任务的函数
			g.Add(func() error {
				err := logJobManager.SyncManager(ctxAll, logJobsyncChan)
				if err != nil {
					level.Error(logger).Log("msg", "TickerInfoCollectAndReport.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}

		{
			// 统计计数的实体的管理器，接收ap 处理
			g.Add(func() error {
				err := pm.UpdateManager(ctxAll)
				if err != nil {
					level.Error(logger).Log("msg", "PointCounterManager.UpdateManager.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}
		{
			// 统计任务实体转换为prometheus的metrics的任务
			g.Add(func() error {
				err := pm.SetMetricsManager(ctxAll)
				if err != nil {
					level.Error(logger).Log("msg", "PointCounterManager.SetMetricsManager.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}

		{
			// logjob 结果的metrics http server
			g.Add(func() error {
				errChan := make(chan error, 1)
				go func() {
					errChan <- metrics.StartMetricWeb(sConfig.HttpAddr)
				}()
				select {
				case err := <-errChan:
					level.Error(logger).Log("msg", "logjob.metrics.web.server.error", "err", err)
					return err
				case <-ctxAll.Done():
					level.Info(logger).Log("msg", "receive_quit_signal_web_server_exit")
					return nil
				}

			}, func(err error) {
				cancelAll()
			},
			)
		}

		{
			// logJob和server直接同步的
			g.Add(func() error {
				err := logjob.TickerLogJobSync(rpcCli, ctxAll, logJobsyncChan, localConfigJobs, metricsMap, common.GetHostName())
				if err != nil {
					level.Error(logger).Log("msg", "PointCounterManager.SetMetricsManager.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}

	}

	{
		// 03 task用到的
		g.Add(func() error {
			err := rpc.TickerTaskReport(rpcCli, ctxAll)
			if err != nil {
				level.Error(logger).Log("msg", "taskjob.TickerTaskReport.error", "err", err)
				return err
			}
			return err

		}, func(err error) {
			cancelAll()
		},
		)
	}
	{
		//	04 xprober用到的

		{
			// 同步探测的目标
			g.Add(func() error {
				err := xprober.TickerGetProberTargets(rpcCli, ctxAll)
				if err != nil {
					level.Error(logger).Log("msg", "xprober.TickerGetProberTargets.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}

		{
			// 上报探测的结果
			g.Add(func() error {
				err := xprober.TickerPushProberResults(rpcCli, ctxAll)
				if err != nil {
					level.Error(logger).Log("msg", "xprober.TickerPushProberResults.error", "err", err)
					return err
				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}

	}

	g.Run()

}

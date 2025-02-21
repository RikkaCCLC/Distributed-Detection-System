package main

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/run"
	"github.com/prometheus/common/promlog"
	promlogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"open-devops/src/models"
	"open-devops/src/modules/server/config"
	mem_index "open-devops/src/modules/server/mem-index"
	"open-devops/src/modules/server/metric"
	"open-devops/src/modules/server/rpc"
	"open-devops/src/modules/server/statistics"
	"open-devops/src/modules/server/task"
	"open-devops/src/modules/server/web"
	"open-devops/src/modules/server/xprober"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	// 命令行解析
	app = kingpin.New(filepath.Base(os.Args[0]), "The open-devops-server")
	// 指定配置文件
	configFile = app.Flag("config.file", "open-devops-server configuration file path").Short('c').Default("open-devops-server.yml").String()
)

func main() {
	// 版本信息
	app.Version(version.Print("open-devops-server"))
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
	level.Info(logger).Log("msg", "load.config.success", "file.path", *configFile, "content.mysql.num", len(sConfig.MysqlS))

	rand.Seed(time.Now().UnixNano())
	// 初始化mysql
	models.InitMySQL(sConfig.MysqlS)

	// 初始化内存倒排索引
	mem_index.Init(logger, sConfig.IndexModules)
	level.Info(logger).Log("msg", "load.mysql.success", "db.num", len(models.DB))
	//models.AddResourceHostTest()
	//models.StreePathAddTest(logger)
	//models.StreePathQueryTest1(logger)
	//models.StreePathQueryTest2(logger)
	//models.StreePathQueryTest3(logger)
	//time.Sleep(10 * time.Second)
	//models.StreePathDelTest(logger)
	//models.StreePathQueryTest1(logger)

	// 注册stree 相关的metrics
	metric.NewMetrics()
	// 注册xprober相关的metrics
	xprober.NewMetrics()
	// 初始化task的本地cache map
	task.TaskCacheInit()

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
		// rpc server
		g.Add(func() error {
			errChan := make(chan error, 1)
			go func() {
				errChan <- rpc.Start(sConfig.RpcAddr, logger)
			}()
			select {
			case err := <-errChan:
				level.Error(logger).Log("msg", "rpc server error", "err", err)
				return err
			case <-ctxAll.Done():
				level.Info(logger).Log("msg", "receive_quit_signal_rpc_server_exit")
				return nil
			}

		}, func(err error) {
			cancelAll()
		},
		)
	}

	{
		// http server
		g.Add(func() error {
			errChan := make(chan error, 1)
			go func() {
				errChan <- web.StartGin(sConfig.HttpAddr, logger)
			}()
			select {
			case err := <-errChan:
				level.Error(logger).Log("msg", "web server error", "err", err)
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
		// 公有云同步
		if sConfig.PCC.Enable {
			cloudsync.Init(logger)

			g.Add(func() error {
				err := cloudsync.CloudSyncManager(ctxAll, logger)
				if err != nil {
					level.Error(logger).Log("msg", "cloudsync.CloudSyncManager.error", "err", err)

				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}
	}

	{
		// 刷新倒排索引
		cloudsync.Init(logger)

		g.Add(func() error {
			err := mem_index.RevertedIndexSyncManager(ctxAll, logger)
			if err != nil {
				level.Error(logger).Log("msg", "mem_index.RevertedIndexSyncManager.error", "err", err)

			}
			return err

		}, func(err error) {
			cancelAll()
		},
		)
	}
	{
		// 统计资源分布

		g.Add(func() error {
			err := statistics.TreeNodeStatisticsManager(ctxAll, logger)
			if err != nil {
				level.Error(logger).Log("msg", "statistics.TreeNodeStatisticsManager.error", "err", err)

			}
			return err

		}, func(err error) {
			cancelAll()
		},
		)
	}
	{
		// 任务执行同步任务

		g.Add(func() error {
			err := task.SyncTaskManager(ctxAll, logger)
			if err != nil {
				level.Error(logger).Log("msg", "task.SyncTaskManager.error", "err", err)

			}
			return err

		}, func(err error) {
			cancelAll()
		},
		)
	}

	{
		// 刷新探测目标池的任务
		tfm := xprober.NewTargetFlushManager(logger, *configFile)
		// target flush manager
		g.Add(func() error {
			err := tfm.Run(ctxAll)
			return err
		}, func(err error) {
			cancelAll()
		})
	}

	{

		// data proceess.
		g.Add(func() error {
			err := xprober.DataProcess(ctxAll, logger)
			return err
		}, func(err error) {
			cancelAll()
		})
	}

	g.Run()

}

package xprober

import (
	"context"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
	"open-devops/src/modules/agent/rpc"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	PbResMap           = sync.Map{}
	Probers            map[string]ProbeFn
	LTM                *LocalTargetManger
	ProberFuncInterval = 15 * time.Second
)

func NewLocalTargetManger(logger log.Logger) {
	localM := make(map[string]*LocalTarget)
	LTM = &LocalTargetManger{}
	LTM.logger = logger
	LTM.Map = localM
}

//type ProbeFn func(ctx context.Context, lt *LocalTarget, logger log.Logger) pb.ProberResultOne
type ProbeFn func(lt *LocalTarget) []*common.ProberResultOne

type LocalTargetManger struct {
	logger log.Logger
	mux    sync.RWMutex
	Map    map[string]*LocalTarget
}

func (ltm *LocalTargetManger) GetMapKeys() []string {
	//ltm.mux.RLock()
	//defer ltm.mux.RUnlock()
	count := len(ltm.Map)
	keys := make([]string, count)
	i := 0
	for hostname := range ltm.Map {
		keys[i] = hostname
		i++
	}
	return keys
}

func (ltm *LocalTargetManger) realRefreshWork(tgs *common.ProberTargetsGetResponse) {
	if len(tgs.Targets) == 0 {
		return
	}
	level.Info(ltm.logger).Log("msg", "realRefreshWork start")
	LTM.mux.Lock()
	defer LTM.mux.Unlock()
	remoteTargetIds := make(map[string]bool)

	localIds := LTM.GetMapKeys()
	for _, t := range tgs.Targets {

		t := t
		logger.Infof("realRefreshWork.Targets:%+v", t)
		pbFunc, funcExists := Probers[t.ProberType]
		if funcExists == false {
			continue
		}

		for _, addr := range t.Target {
			thisId := t.Region + addr + t.ProberType
			remoteTargetIds[thisId] = true
			if _, ok := LTM.Map[thisId]; ok {
				continue
			}

			nt := &LocalTarget{
				logger:       LTM.logger,
				Addr:         addr,
				SourceRegion: LocalRegion,
				TargetRegion: t.Region,
				ProbeType:    t.ProberType,
				Prober:       pbFunc,
				QuitChan:     make(chan struct{}),
			}
			LTM.Map[thisId] = nt
			go nt.Start()

		}

	}
	// stop old
	for _, key := range localIds {
		if _, found := remoteTargetIds[key]; !found {
			LTM.Map[key].Stop()
			delete(LTM.Map, key)
		}
	}

}

type LocalTarget struct {
	logger       log.Logger
	Addr         string
	SourceRegion string
	TargetRegion string

	ProbeType string
	Prober    ProbeFn
	QuitChan  chan struct{}
}

func (lt *LocalTarget) Uid() string {

	return lt.TargetRegion + lt.Addr + lt.ProbeType
}

func (lt *LocalTarget) Start() {
	ticker := time.NewTicker(ProberFuncInterval)
	level.Info(lt.logger).Log("msg", "LocalTarget probe start....", "uid", lt.Uid())
	defer ticker.Stop()
	for {
		select {
		case <-lt.QuitChan:
			level.Info(lt.logger).Log("msg", "receive_quit_signal", "uid", lt.Uid())
			return
		case <-ticker.C:
			res := lt.Prober(lt)
			if len(res) > 0 {
				PbResMap.Store(lt.Uid(), res)
			}

		}

	}

}

func (lt *LocalTarget) Stop() {
	close(lt.QuitChan)
}

// 获取探测任务的ticker
func TickerGetProberTargets(cli *rpc.RpcCli, ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	req := common.ProberTargetsGetRequest{
		LocalRegion: LocalRegion,
		LocalIp:     LocalIp,
	}

	for {
		select {
		case <-ctx.Done():
			logger.Infof("TickerGetProberTargets.receive_quit_signal_and_quit")
			return nil
		case <-ticker.C:
			res := cli.ProberTargetSync(req)
			if res != nil {

				LTM.realRefreshWork(res)
			}

		}
	}

}

// 推送探测结果的ticker
func TickerPushProberResults(cli *rpc.RpcCli, ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("TickerPushProberResults.receive_quit_signal_and_quit")
			return nil
		case <-ticker.C:
			pushPbResults(cli)
		}
	}

}

func pushPbResults(cli *rpc.RpcCli) {

	prs := make([]*common.ProberResultOne, 0)
	f := func(k, v interface{}) bool {

		va := v.([]*common.ProberResultOne)
		logger.Infof("PbResMap.Range:%+v", va)
		prs = append(prs, va...)
		return true
	}
	PbResMap.Range(f)
	req := common.ProberResultPushRequest{}
	req.ProberResults = prs

	if len(req.ProberResults) == 0 {
		logger.Warningf("empty_result_list_not_to_push")
		return
	}

	// TODO remove
	for index, i := range req.ProberResults {
		logger.Infof("index :%+v,res:%+v", index, i)
	}

	res := cli.PushProberResults(req)

	logger.Infof("pushPbResults:%+v req.num:%+v", res, len(req.ProberResults))
}

func Init(logger log.Logger, region string) {
	Probers = map[string]ProbeFn{
		"http": ProbeHTTP,
		"icmp": ProbeICMP,
		//"icmp": ProbeHTTP,
	}
	if region != "" {
		//	代表用户指定了region，
		LocalRegion = region

	} else {
		// 否则通过公有云接口获取
		GetLocalRegionByEc2(logger)
	}

	GetLocalIp(logger)
	NewLocalTargetManger(logger)
}

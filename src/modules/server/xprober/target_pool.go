package xprober

import (
	"context"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
	"open-devops/src/modules/server/config"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"fmt"
)

type TargetPool struct {
	ProbeType     string
	regionTargets sync.Map
}

var (
	IcmpRegionProberMap  = sync.Map{} // icmp的池子
	OtherRegionProberMap = sync.Map{} //http的池子
	AgentIpRegionMap     = sync.Map{} // agent ip的map
)

type TargetFlushManager struct {
	Logger     log.Logger
	ConfigFile string
}

func rangeIcmpMap() {
	f := func(k, v interface{}) bool {
		region := k.(string)
		data := v.(*common.ProberTargets)
		fmt.Println("rangeIcmpMap", region, data)
		return true
	}

	IcmpRegionProberMap.Range(f)
}

func (t *TargetFlushManager) flushAgentIpIntoGlobalMap() {
	level.Info(t.Logger).Log("msg", "flushAgentIpIntoGlobalMap run....")
	tmpM := make(map[string][]string)

	f := func(k, v interface{}) bool {
		ip := k.(string)
		region := v.(string)
		logger.Infof("[AgentIpRegionMap.show][ip:%+v][region:%+v]", ip, region)
		tmpM[region] = append(tmpM[region], ip)

		return true
	}
	AgentIpRegionMap.Range(f)
	for region, ips := range tmpM {
		tNew := &common.ProberTargets{}
		tNew.Region = region
		tNew.ProberType = "icmp"
		tNew.Target = ips

		preData, loaded := IcmpRegionProberMap.LoadOrStore(region, tNew)
		preDataN := preData.(*common.ProberTargets)
		if loaded {
			thisT := tNew.Target
			originT := preDataN.Target
			thisTM := make(map[string]string)
			for _, tt := range thisT {
				thisTM[tt] = tt
			}
			for _, tt := range originT {
				if _, exists := thisTM[tt]; exists == false {
					thisT = append(thisT, tt)
				}
			}
			a := &common.ProberTargets{
				Region:     region,
				ProberType: "icmp",
				Target:     thisT,
			}
			IcmpRegionProberMap.Store(region, a)
		}

	}
	//rangeIcmpMap()

}

func NewTargetFlushManager(logger log.Logger, configFile string) *TargetFlushManager {

	return &TargetFlushManager{Logger: logger, ConfigFile: configFile}
}
func (t *TargetFlushManager) Run(ctx context.Context) error {

	ticker := time.NewTicker(10 * time.Second)
	level.Info(t.Logger).Log("msg", "TargetFlushManager start....")
	t.refresh()
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t.refresh()

		case <-ctx.Done():
			level.Info(t.Logger).Log("msg", "TargetFlushManager exit....")
			return nil
		}
	}

}

func (t *TargetFlushManager) refreshFromConfigFile() {
	level.Info(t.Logger).Log("msg", "refreshFromConfigFile run....")

	config, _ := config.LoadFile(t.ConfigFile)
	otmpM := make(map[string][]*common.ProberTargets)
	icmpM := make(map[string]*common.ProberTargets)
	if len(config.ProberTargets) <= 0 {
		level.Info(t.Logger).Log("msg", "refreshFromConfigFile empty targets....")
		return
	}
	for _, t := range config.ProberTargets {
		tNew := &common.ProberTargets{}
		tNew.Region = t.Region
		tNew.ProberType = t.ProberType
		tNew.Target = t.Target
		otmpM[tNew.Region] = append(otmpM[tNew.Region], tNew)
		switch t.ProberType {
		case "icmp":
			icmpM[tNew.Region] = tNew
		default:
			otmpM[tNew.Region] = append(otmpM[tNew.Region], tNew)
		}

	}
	for k, v := range icmpM {
		preData, loaded := IcmpRegionProberMap.LoadOrStore(k, v)
		preDataN := preData.(*common.ProberTargets)
		if loaded {
			thisT := v.Target
			originT := preDataN.Target
			thisTM := make(map[string]string)
			for _, tt := range thisT {
				thisTM[tt] = tt
			}
			for _, tt := range originT {
				if _, exists := thisTM[tt]; exists == false {
					thisT = append(thisT, tt)
				}
			}
			a := &common.ProberTargets{
				Region:     k,
				ProberType: "icmp",
				Target:     thisT,
			}
			IcmpRegionProberMap.Store(k, a)
		}

	}

	for k, v := range otmpM {

		OtherRegionProberMap.Store(k, v)
	}
}

func (t *TargetFlushManager) refresh() {
	go t.refreshFromConfigFile()
	go t.flushAgentIpIntoGlobalMap()
}

func GetTargetsByRegion(sourceRegion string) (res []*common.ProberTargets) {

	f := func(k, v interface{}) bool {
		//key := k.(string)
		va := v.([]*common.ProberTargets)
		//if key != sourceRegion {
		res = append(res, va...)

		//}
		return true
	}
	fi := func(k, v interface{}) bool {
		key := k.(string)
		va := v.(*common.ProberTargets)
		if key != sourceRegion {
			res = append(res, va)

		}
		return true
	}

	IcmpRegionProberMap.Range(fi)
	OtherRegionProberMap.Range(f)
	return
}

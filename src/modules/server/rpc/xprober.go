package rpc

import (
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
	"open-devops/src/modules/server/xprober"
)

func (*Server) GetProberTargets(req common.ProberTargetsGetRequest, res *common.ProberTargetsGetResponse) error {
	region := req.LocalRegion
	xprober.AgentIpRegionMap.Store(req.LocalIp, req.LocalRegion)
	tgs := xprober.GetTargetsByRegion(region)
	res.Targets = tgs
	return nil
}

func (*Server) PushProberResults(req common.ProberResultPushRequest, res *common.ProberResultPushResponse) error {
	suNum := 0
	logger.Infof("PushProberResults.req:%+v", req)
	for _, prr := range req.ProberResults {
		prr := prr
		logger.Infof("PushProberResults.prr:%+v", prr)
		uid := GetProbeResultUid(prr)
		switch prr.ProbeType {
		case `icmp`:
			xprober.IcmpDataMap.Store(uid, prr)
		case `http`:
			xprober.HttpDataMap.Store(uid, prr)
		}
		suNum += 1

	}
	res.SuccessNum = int32(suNum)
	return nil

}

func GetProbeResultUid(prr *common.ProberResultOne) (uid string) {
	uid = prr.WorkerName + prr.MetricName + prr.SourceRegion + prr.TargetRegion + prr.ProbeType + prr.TargetAddr
	return

}

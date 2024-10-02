package counter

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/toolkits/pkg/logger"
	"math"
	"open-devops/src/common"
	"open-devops/src/modules/agent/consumer"
	"sync"
	"time"
)

type PointCounterManager struct {
	sync.RWMutex
	CounterQueue chan *consumer.AnalysPoint
	// key 是标签排序后的string
	TagStringMap map[string]*PointCounter
	MetricsMap   map[string]*prometheus.GaugeVec
}

func NewPointCounterManager(cq chan *consumer.AnalysPoint, m map[string]*prometheus.GaugeVec) *PointCounterManager {
	pm := &PointCounterManager{
		CounterQueue: cq,
		TagStringMap: make(map[string]*PointCounter),
		MetricsMap:   m,
	}
	return pm

}

func (pm *PointCounterManager) GetPcByUniqueName(seriesId string) *PointCounter {
	pm.RLock()
	defer pm.RUnlock()
	return pm.TagStringMap[seriesId]
}

func (pm *PointCounterManager) SetPc(seriesId string, pc *PointCounter) {
	pm.Lock()
	defer pm.Unlock()
	pm.TagStringMap[seriesId] = pc
}

// 打点的方法
func (pm *PointCounterManager) SetMetrics() {
	pm.Lock()
	defer pm.Unlock()
	for _, pc := range pm.TagStringMap {
		pc := pc
		metric, loaded := pm.MetricsMap[pc.MetricsName]
		if !loaded {
			logger.Errorf("[metrics.notfound][name:%v]", pc.MetricsName)
			continue
		}
		logger.Infof("[PointCounterManager.SetMetrics][pc:%+v]", pc)
		var value float64
		switch pc.LogFunc {
		case common.LogFuncCnt:
			value = float64(pc.Count)
		case common.LogFuncSum:
			value = float64(pc.Sum)
		case common.LogFuncMax:
			value = float64(pc.Max)
		case common.LogFuncMin:
			value = float64(pc.Min)
		case common.LogFuncAvg:
			value = float64(pc.Sum) / float64(pc.Count)
		}
		metric.With(prometheus.Labels(pc.LabelMap)).Set(value)

	}
}

func (pm *PointCounterManager) SetMetricsManager(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("PointCounterManager.SetMetricsManager.receive_quit_signal_and_quit")
			return nil
		case <-ticker.C:
			pm.SetMetrics()

		}

	}
}

func (pm *PointCounterManager) UpdateManager(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("PointCounterManager.UpdateManager.receive_quit_signal_and_quit")
			return nil
		case ap := <-pm.CounterQueue:
			pc := pm.GetPcByUniqueName(ap.MetricsName + ap.SortLabelString)
			if pc == nil {
				pc = NewPointCounter(ap.MetricsName, ap.SortLabelString, ap.LogFunc, ap.LabelMap)
				pm.SetPc(ap.MetricsName+ap.SortLabelString, pc)
			}
			pc.Update(ap.Value)

		}

	}

}

// 统计实体
type PointCounter struct {
	sync.RWMutex
	Count           int64   // 日志条数计数
	Sum             float64 // 正则数字的sum
	Max             float64 // 正则数字的Max
	Min             float64 // 正则数字的Min
	Ts              int64   // 时间戳
	MetricsName     string  // metrics name
	LogFunc         string  // 计算的方法 ，cnt/max/min
	SortLabelString string  //标签排序的结果
	LabelMap        map[string]string
}

func NewPointCounter(metricsName string, sortLabelString string, logFunc string, labelMap map[string]string) *PointCounter {
	pc := &PointCounter{
		MetricsName:     metricsName,
		LogFunc:         logFunc,
		SortLabelString: sortLabelString,
		LabelMap:        labelMap,
	}
	return pc
}

func (pc *PointCounter) Update(value float64) {
	pc.Lock()
	defer pc.Unlock()
	pc.Sum = pc.Sum + value
	if math.IsNaN(pc.Max) || value > pc.Max {
		pc.Max = value
	}
	if math.IsNaN(pc.Min) || value < pc.Min {
		pc.Min = value
	}
	pc.Count += 1
	pc.Ts = time.Now().Unix()
}

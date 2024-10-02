package consumer

import (
	"fmt"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
	"open-devops/src/models"
)

// 定义消费者组
type ConsumerGroup struct {
	Consumers   []*Consumer
	ConsumerNum int
}

// TODO 缺analyPoint
func NewConsumerGroup(filePath string, stream chan string, stra *models.LogStrategy, cq chan *AnalysPoint) *ConsumerGroup {
	cNum := common.ConsumerNum
	cg := &ConsumerGroup{
		Consumers:   make([]*Consumer, 0),
		ConsumerNum: cNum,
	}
	logger.Infof("[NewConsumerGroup  ][file:%s][num:%d]", filePath, cNum)
	for i := 0; i < cNum; i++ {
		mark := fmt.Sprintf("[log.consumer][file:%s][num:%d/%d]", filePath, i+1, cNum)
		c := &Consumer{
			FilePath:     filePath,
			Stream:       stream,
			Stra:         stra,
			Mark:         mark,
			CounterQueue: cq,
			Close:        make(chan struct{}),
		}
		cg.Consumers = append(cg.Consumers, c)
	}
	return cg

}

func (cg *ConsumerGroup) Start() {
	for i := 0; i < cg.ConsumerNum; i++ {
		cg.Consumers[i].Start()
	}
}

func (cg *ConsumerGroup) Stop() {
	for i := 0; i < cg.ConsumerNum; i++ {
		cg.Consumers[i].Stop()
	}
}

package logjob

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/common"
	"open-devops/src/models"
	"open-devops/src/modules/agent/consumer"
	"open-devops/src/modules/agent/reader"
)

type LogJob struct {
	r    *reader.Reader          // 读取日志
	c    *consumer.ConsumerGroup // 代表我们的消费者组
	Stra *models.LogStrategy     // 策略

}

func (lj *LogJob) hash() string {
	md5obj := md5.New()
	md5obj.Write([]byte(lj.Stra.MetricName))
	md5obj.Write([]byte(lj.Stra.FilePath))
	return hex.EncodeToString(md5obj.Sum(nil))
}

func (lj *LogJob) start(cq chan *consumer.AnalysPoint) {

	fPath := lj.Stra.FilePath
	// stream
	stream := make(chan string, common.LogQueueSize)
	// new reader
	r, err := reader.NewReader(fPath, stream)
	if err != nil {
		return
	}
	lj.r = r
	// new consumer
	cg := consumer.NewConsumerGroup(fPath, stream, lj.Stra, cq)
	lj.c = cg
	// 启动 r 和c
	// 先消费者
	lj.c.Start()
	// 后生产者
	go r.Start()

	//logger.Infof("[create.LogJob.start][stra:%+v]", lj.Stra,)
	logger.Infof("[create.LogJob.start][stra:%+v]", lj.Stra)

}

// 先停生成者，后消费者
func (lj *LogJob) stop() {
	lj.r.Stop()
	lj.c.Stop()
}

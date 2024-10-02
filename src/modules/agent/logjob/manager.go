package logjob

import (
	"context"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/modules/agent/consumer"
	"sync"
)

type LogJobManager struct {
	targetMtx     sync.Mutex
	activeTargets map[string]*LogJob
	cq            chan *consumer.AnalysPoint
}

func NewLogJobManager(cq chan *consumer.AnalysPoint) *LogJobManager {
	return &LogJobManager{
		activeTargets: make(map[string]*LogJob),
		cq:            cq,
	}
}

func (jm *LogJobManager) StopALl() {
	jm.targetMtx.Lock()
	defer jm.targetMtx.Unlock()
	for _, v := range jm.activeTargets {
		v.stop()
	}
}

func (jm *LogJobManager) SyncManager(ctx context.Context, syncChan chan []*LogJob) error {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("LogJobManager.receive_quit_signal_and_quit")
			jm.StopALl()
			return nil
		case jobs := <-syncChan:
			jm.Sync(jobs)
		}

	}
}

func (jm *LogJobManager) Sync(jobs []*LogJob) {
	logger.Infof("[LogJobManager.sync][num:%d][res:%+v]", len(jobs), jobs)
	thisNewTargets := make(map[string]*LogJob)
	thisAllTargets := make(map[string]*LogJob)

	jm.targetMtx.Lock()
	for _, t := range jobs {
		hash := t.hash()

		thisAllTargets[hash] = t
		if _, loaded := jm.activeTargets[hash]; !loaded {
			thisNewTargets[hash] = t
			jm.activeTargets[hash] = t
		}
	}

	// 停止旧的
	for hash, t := range jm.activeTargets {
		if _, loaded := thisAllTargets[hash]; !loaded {
			logger.Infof("stop %+v stra:%+v", t, t.Stra)
			t.stop()
			delete(jm.activeTargets, hash)
		}
	}

	jm.targetMtx.Unlock()
	// 开启新的
	for _, t := range thisNewTargets {
		t := t
		t.start(jm.cq)
		//t.start(jm.cq)
	}

}

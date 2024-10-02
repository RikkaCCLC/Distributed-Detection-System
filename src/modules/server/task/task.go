package task

import (
	"context"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/toolkits/pkg/logger"
	"open-devops/src/models"
	"sync"
	"time"
)

var TaskCaches *TaskCache

type TaskCache struct {
	sync.RWMutex
	M map[string][]*models.TaskMeta //这个map代表机器ip对应的任务数组，因为一个机器可能被同时配置多个任务
}

func TaskCacheInit() {
	TaskCaches = &TaskCache{
		M: make(map[string][]*models.TaskMeta),
	}
}

func (tc *TaskCache) GetTasksByIp(ip string) []*models.TaskMeta {
	tc.Lock()
	defer tc.Unlock()
	res, loaded := tc.M[ip]
	if !loaded {
		res = make([]*models.TaskMeta, 0)
	}
	return res
}

// 定义一个同步db中的任务方法，周期执行doSyncTask
func SyncTaskManager(ctx context.Context, logger log.Logger) error {

	level.Info(logger).Log("msg", "SyncTaskManager.start")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			level.Info(logger).Log("msg", "SyncTaskManager.exit.receive_quit_signal")
			return nil
		case <-ticker.C:

			TaskCaches.doSyncTask()
		}
	}
}

// doSyncTask从db中查询出所有为完成的任务，对应就是done=0的task
// 然后遍历 task ， 按照机器ip塞入map中
func (tc *TaskCache) doSyncTask() {
	// 获取为完成的任务
	tasks, err := models.TaskMetaGetUnDo()
	if err != nil {
		logger.Errorf("[doSyncTask][error: %+v]", err)
		return
	}
	m := make(map[string][]*models.TaskMeta)
	for _, t := range tasks {

		err := json.Unmarshal([]byte(t.HostsRaw), &t.Hosts)
		if err != nil {
			logger.Errorf("[json.Unmarshal_host][t:%+v][error: %+v]", t, err)
			continue
		}
		//logger.Debugf("[get_task_from_mysql][t:%+v]", t)
		if len(t.Hosts) == 0 {
			continue
		}

		for _, host := range t.Hosts {
			// 按照机器ip分发
			tasks, loaded := m[host]
			if !loaded {

				tasks = make([]*models.TaskMeta, 0)
			}
			t.Action = "start"

			tasks = append(tasks, t)

			m[host] = tasks
		}
	}

	tc.Lock()
	defer tc.Unlock()
	tc.M = m
	//logger.Debugf("[tc.M:%+v]", tc.M)
}

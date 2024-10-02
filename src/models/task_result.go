package models

import "fmt"

type TaskResult struct {
	Id     int64  `json:"id"`
	TaskId int64  `json:"task_id"` //归属于哪个任务
	Host   string `json:"host"`    //哪个机器
	Status string `json:"status"`  //任务执行的结果eg:success failed ...
	Stdout string `json:"stdout"`  //标准输出
	Stderr string `json:"stderr"`  // 标准错误
}

func (ts *TaskResult) Save() (error, bool) {
	var obj TaskResult
	sql := fmt.Sprintf(" task_id=%d and host='%s' ", ts.TaskId, ts.Host)
	has, err := DB["stree"].Where(sql).Get(&obj)
	if err != nil {
		return err, false
	}
	if has {
		return nil, false
	}
	_, err = DB["stree"].Insert(ts)
	return nil, true
}

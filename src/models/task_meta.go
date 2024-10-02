package models

import (
	"fmt"
	"time"
)

type TaskMeta struct {
	Id       int64     `json:"id"`
	Title    string    `json:"title"`                  // 标题
	Account  string    `json:"account"`                //脚本执行账号
	Timeout  int       `json:"timeout"`                // 执行超时
	Script   string    `json:"script"`                 //执行的脚本
	Args     string    `json:"args"`                   //执行的脚本的参数
	Creator  string    `json:"creator"`                //创建者
	Created  time.Time `xorm:"created" json:"created"` //创建时间
	HostsRaw string    `json:"hosts"`                  // 执行机器的ip列表json
	Hosts    []string  `xorm:"-" json:"-"`
	Done     int       `xorm:"done" json:"done"` //任务结束与否的标志位=0未结束，=1结束
	Clock    int64     `xorm:"-" json:"clock"`
	Action   string    `xorm:"-" json:"action"`
}

// 查询一条
func TaskMetaGet(where string, args ...interface{}) (*TaskMeta, error) {
	var obj TaskMeta
	has, err := DB["stree"].Where(where, args...).Get(&obj)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return &obj, nil
}

func TaskMetaGetUnDo() ([]*TaskMeta, error) {

	var objs []*TaskMeta
	session := DB["stree"].Where("done=0 ")
	err := session.Find(&objs)
	return objs, err

}

// 将任务标记为已完成
func MarkTaskMetaDone(id int64) error {
	session := DB["stree"].NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	sql := fmt.Sprintf("UPDATE task_meta SET done=1  WHERE id=%d ", id)

	if _, err := session.Exec(sql); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

// 插入一条记录
func (tm *TaskMeta) AddOne() (int64, error) {
	_, err := DB["stree"].InsertOne(tm)
	return tm.Id, err
}

// 查询多条
func TaskMetaGets(where string, args ...interface{}) ([]*TaskMeta, error) {
	var obj []*TaskMeta
	err := DB["stree"].Table("task_meta").Where(where, args...).Find(&obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
// 将任务的action设置为kill
func SetTaskKill(id int64) error {
	session := DB["stree"].NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	sql := fmt.Sprintf("UPDATE task_meta SET action='kill'  WHERE id=%d ", id)

	if _, err := session.Exec(sql); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

package models

import (
	"fmt"
	"open-devops/src/modules/server/config"
	"time"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

var DB = map[string]*xorm.Engine{}

func InitMySQL(mysqlS []*config.MySQLConf) {
	for _, conf := range mysqlS {
		db, err := xorm.NewEngine("mysql", conf.Addr)
		if err != nil {
			fmt.Printf("[init.mysql.error][cannot connect to mysql][addr:%v][err:%v]\n", conf.Addr, err)
			continue
		}
		db.SetMaxIdleConns(conf.Idle)
		db.SetMaxOpenConns(conf.Max)
		db.SetConnMaxLifetime(time.Hour)
		db.ShowSQL(conf.Debug)
		db.Logger().SetLevel(xlog.LOG_INFO)
		DB[conf.Name] = db
	}
}

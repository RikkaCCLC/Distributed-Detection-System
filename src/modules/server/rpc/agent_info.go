package rpc

import (
	"encoding/json"
	"fmt"
	"open-devops/src/models"
)

func (*Server) HostInfoReport(input models.AgentCollectInfo, output *string) error {

	// 统一字段

	ips := []string{input.IpAddr}
	ipJ, _ := json.Marshal(ips)
	if input.SN == "" {
		input.SN = input.HostName
	}
	if input.SN == "" {
		*output = "sn.empty"
		return nil
	}
	// 先获取对象的uid
	rh := models.ResourceHost{
		Uid:        input.SN,
		Name:       input.HostName,
		PrivateIps: ipJ,
		CPU:        input.CPU,
		Mem:        input.Mem,
		Disk:       input.Disk,
	}

	hash := rh.GenHash()
	// 用uid去db中获取之前的结果，再根据两者的hash是否一致决定 更改
	rhUid := models.ResourceHost{Uid: input.SN}

	rhUidDb, err := rhUid.GetOne()
	if err != nil {
		*output = "db_error"
		return nil
	}
	if rhUidDb == nil {
		//	说明指定udi不存在，插入
		rh.Hash = hash
		err = rh.AddOne()
		if err != nil {
			*output = fmt.Sprintf("db_error_%v", err)

		} else {
			*output = "insert_success"
		}
		return nil
	}

	// uid存在需要判断hash
	if rhUidDb.Hash != hash {
		rh.Hash = hash
		updated, err := rh.Update()
		if err != nil {
			*output = "update_error"
			return nil
		}
		if updated {
			*output = "update_success"
			return nil
		}
	}
	// uid存在并且hash相等 啥都不需要做
	//log.Printf("[host.info.same][HostInfoReport][input:%+v]", input)

	return nil
}

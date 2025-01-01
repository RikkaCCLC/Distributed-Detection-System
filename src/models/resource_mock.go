package models

import (
	"encoding/json"
	"fmt"
)

type ResourceHostTest struct {
	Id         int64           `json:"id"`
	Name       string          `json:"name"`
	PrivateIps json.RawMessage `json:"private_ips"`
	Tags       json.RawMessage `json:"tags"`
}

func (rh *ResourceHostTest) AddOne() error {
	_, err := DB["stree"].InsertOne(rh)
	return err
}

func (rh *ResourceHostTest) GetOne() (*ResourceHostTest, error) {

	has, err := DB["stree"].Get(rh)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return rh, nil

}

func AddResourceHostTest() {
	//m := map[string]string{"region": "bj", "app": "live"}
	//ips := []string{"1.1.1.1", "2.2.2.2"}
	////mRaw, _ := json.Marshal(m)
	////ipRaw, _ := json.Marshal(ips)

	//rh := ResourceHostTest{
	//	Name:       "abc",
	//	PrivateIps: ipRaw,
	//	Tags:       mRaw,
	//}
	//err := rh.AddOne()
	//fmt.Println(err)
	rhNew := ResourceHostTest{
		Name: "abc",
	}
	rhDb, err := rhNew.GetOne()
	mTag := make(map[string]string)
	err = json.Unmarshal(rhDb.Tags, &mTag)

	ipsN := make([]string, 3)
	err = json.Unmarshal(rhDb.PrivateIps, &ipsN)

	fmt.Println(mTag, err)
	fmt.Println(ipsN, err)

}

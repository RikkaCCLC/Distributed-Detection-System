package cloudsync

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"math/rand"
	"open-devops/src/common"
	"open-devops/src/models"
	"time"
)

type HostSync struct {
	CloudAlibaba
	CloudTencent
	TableName string
	Logger    log.Logger
}

func (this *HostSync) sync() {
	// 去调用公有云的sdk 取数据，我们使用一个mock的方法
	start := time.Now()
	hs := genMockResourceHost()

	// 获取本地的uid对应的hashM
	uuidHashM, err := models.GetHostUidAndHash()
	if err != nil {
		level.Error(this.Logger).Log("msg", "models.GetHostUidAndHash.error", "err", err)
		return
	}
	// - toAddSet , toModSet存放的都是 host对象
	//- toDelIds 放的是待删的uids
	toAddSet := make([]models.ResourceHost, 0)
	toModSet := make([]models.ResourceHost, 0)
	toDelIds := make([]string, 0)

	localUidSet := make(map[string]struct{})
	var toAddNum, toModNum, toDelNum int
	var suAddNum, suModNum, suDelNum int
	for _, h := range hs {
		localUidSet[h.Uid] = struct{}{}
		hash, ok := uuidHashM[h.Uid]
		if !ok {
			// 说明本地没有，公有云有，要新增
			toAddSet = append(toAddSet, h)
			toAddNum++
			continue
		}
		// 存在说明还要判断hash
		if hash == h.Hash {
			continue
		}
		// 说明uid相同当时hash不同，某些字段变更了
		toModSet = append(toModSet, h)
		toModNum++
	}

	for uid := range uuidHashM {
		// 说明db中有这个uid，但是远端公有云没有，要删除
		if _, ok := localUidSet[uid]; !ok {
			toDelIds = append(toDelIds, uid)
			toDelNum++
		}
	}

	// 以上是我们的判断流程
	// 下面是执行
	// 新增
	for _, h := range toAddSet {
		err := h.AddOne()
		if err != nil {
			level.Error(this.Logger).Log("msg", "ResourceHost.AddOne.error", "err", err, "name", h.Name)
			continue
		}
		suAddNum++
	}
	// 修改
	for _, h := range toModSet {
		isUpdate, err := h.UpdateByUid(h.Uid)
		if err != nil {
			level.Error(this.Logger).Log("msg", "ResourceHost.HostSync.error", "err", err, "name", h.Name)
			continue
		}
		if isUpdate {
			suModNum++
		}
	}
	// 删除
	if len(toDelIds) > 0 {
		num, _ := models.BatchDeleteResource(common.RESOURCE_HOST, "uid", toDelIds)
		suDelNum = int(num)
	}
	timeTook := time.Since(start)
	level.Info(this.Logger).Log("msg", "ResourceHost.HostSync.res.print",
		"public.cloud.num", len(hs),
		"db.num", len(uuidHashM),
		"toAddNum", toAddNum,
		"toModNum", toModNum,
		"toDelNum", toDelNum,
		"suAddNum", suAddNum,
		"suModNum", suModNum,
		"suDelNum", suDelNum,
		"timeTook", timeTook.Seconds(),
	)

}

func genMockResourceHost() []models.ResourceHost {
	rand.Seed(time.Now().UnixNano())
	// g.p.a标签
	randGs := []string{"inf", "ads", "web", "sys"}
	randPs := []string{"monitor", "cicd", "k8s", "mq"}
	randAs := []string{"kafka", "prometheus", "zookeeper", "es"}

	// cpu 等资源随机
	randCpus := []string{"4", "8", "16", "32", "64", "128"}
	randMems := []string{"8", "16", "32", "64", "128", "256", "512"}
	randDisks := []string{"128", "256", "512", "1024", "2048", "4096", "8192"}

	// 标签tags
	randMapKeys := []string{"arch", "idc", "os", "job"}
	randMapvalues := []string{"linux", "beijing", "windows", "shanghai", "arm64", "amd64", "darwin", "shijihulian"}

	// 公有云标签
	randRegions := []string{"beijing", "shanghai", "guangzhou", "tianjin", "shandong"}
	randCloudProviders := []string{"alibaba", "tencent", "aws", "huawei", "azure"}
	randClusters := []string{"bigdata", "inf", "middleware", "business"}
	randInts := []string{"4c8g", "4c16g", "8c32g", "16c64g"}

	// 目的是4选1 返回0-3的数组
	// 比如8-15 15-8 +8
	frn := func(n int) int {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(n)
	}

	// 每次host数量 5-20个
	frNum := func() int {
		rand.Seed(time.Now().UnixNano())
		return int(rand.Int63n(180-65) + 65)
	}
	hs := make([]models.ResourceHost, 0)
	for i := 0; i < frNum(); i++ {
		randN := i
		name := fmt.Sprintf("genMockResourceHost_host_%d", randN)
		ips := []string{fmt.Sprintf("8.8.8.%d", randN)}
		ipJ, _ := json.Marshal(ips)
		h := models.ResourceHost{
			Name:       name,
			PrivateIps: ipJ,
			CPU:        randCpus[frn(len(randCpus)-1)],
			Mem:        randMems[frn(len(randMems)-1)],
			Disk:       randDisks[frn(len(randDisks)-1)],

			StreeGroup:   randGs[frn(len(randGs)-1)],
			StreeProduct: randPs[frn(len(randPs)-1)],
			StreeApp:     randAs[frn(len(randAs)-1)],

			Region:        randRegions[frn(len(randRegions)-1)],
			CloudProvider: randCloudProviders[frn(len(randCloudProviders)-1)],
			InstanceType:  randInts[frn(len(randInts)-1)],
		}
		tagM := make(map[string]string)
		for _, i := range randMapKeys {
			tagM[i] = randMapvalues[frn(len(randMapvalues)-1)]
		}
		// cluster
		tagM["cluster"] = randClusters[frn(len(randClusters)-1)]
		tagMJ, _ := json.Marshal(tagM)
		h.Tags = tagMJ

		hash := h.GenHash()
		h.Hash = hash
		md5o := md5.New()
		md5o.Write([]byte(name))
		h.Uid = hex.EncodeToString(md5o.Sum(nil))
		hs = append(hs, h)
	}
	return hs
}

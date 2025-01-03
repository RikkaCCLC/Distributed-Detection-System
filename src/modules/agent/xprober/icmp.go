package xprober

import (
	"bytes"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"open-devops/src/common"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func execCmd(cmdStr string, logger log.Logger) (success bool, outStr string) {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		level.Error(logger).Log("execCmdMsg", err, "cmd", cmdStr)

		if strings.Contains(err.Error(), "killed") {
			return false, "killed"
		}

		return false, string(stderr.Bytes())
	}
	outStr = string(stdout.Bytes())
	return true, outStr

}

func ProbeICMP(lt *LocalTarget) []*common.ProberResultOne {

	defer func() {
		if r := recover(); r != nil {
			resultErr, _ := r.(error)
			level.Error(lt.logger).Log("msg", "ProbeICMP panic ...", "resultErr", resultErr)

		}
	}()

	pingCmd := fmt.Sprintf("/usr/bin/timeout --signal=KILL 15s  /usr/bin/ping -q -A -f -s 100 -W 1000 -c 50 %s", lt.Addr)
	level.Info(lt.logger).Log("msg", "LocalTarget  ProbeICMP start ...", "uid", lt.Uid(), "pingcmd", pingCmd)
	success, outPutStr := execCmd(pingCmd, lt.logger)
	prs := make([]*common.ProberResultOne, 3)

	var (
		pkgdLine    string
		latenLinke  string
		pkgRateNum  float64
		pingEwmaNum float64
		pingSuccess float64
	)

	pkgRateNum = -1
	pingEwmaNum = -1
	pingSuccess = 0
	prSu := common.ProberResultOne{
		MetricName:   common.MetricsNamePingTargetSuccess,
		WorkerName:   LocalIp,
		TargetAddr:   lt.Addr,
		SourceRegion: LocalRegion,
		TargetRegion: lt.TargetRegion,
		ProbeType:    lt.ProbeType,
		TimeStamp:    time.Now().Unix(),
		Value:        float32(pingSuccess),
	}
	if success == false {
		level.Error(lt.logger).Log("msg", "ProbeICMP failed ...", "uid", lt.Uid(), "err_str", outPutStr)

		if strings.Contains(outPutStr, "killed") {
			level.Error(lt.logger).Log("msg", "ProbeICMP killed ...", "uid", lt.Uid(), "err_str", outPutStr)
			prSu.Value = -1
			prs = append(prs, &prSu)
			return prs

		}
		return prs
	}

	for _, line := range strings.Split(outPutStr, "\n") {
		if strings.Contains(line, "packets transmitted") {
			pkgdLine = line
			continue
		}
		if strings.Contains(line, "min/avg/max/mdev") {
			latenLinke = line
			continue
		}

	}
	/*
		PING 10.21.45.237 (10.21.45.237) 100(128) bytes of data.

		--- 10.21.45.237 ping statistics ---
		50 packets transmitted, 0 received, 100% packet loss, time 499ms
	*/

	if len(pkgdLine) > 0 {

		pkgRate := strings.Split(pkgdLine, " ")[5]
		pkgRate = strings.Replace(pkgRate, "%", "", -1)
		pkgRateNum, _ = strconv.ParseFloat(pkgRate, 64)
	}

	if len(latenLinke) > 0 {
		pingEwmas := strings.Split(latenLinke, " ")

		pingEwma := pingEwmas[len(pingEwmas)-2]
		pingEwma = strings.Split(pingEwma, "/")[1]
		pingEwmaNum, _ = strconv.ParseFloat(pingEwma, 64)
	}

	level.Info(lt.logger).Log("msg", "ProbeICMP_one_res", "pingcmd", pingCmd, "outPutStr", outPutStr, "pkgRateNum", float32(pkgRateNum), "pingEwmaNum", float32(pingEwmaNum))
	prDr := common.ProberResultOne{
		MetricName:   common.MetricsNamePingPackageDrop,
		WorkerName:   LocalIp,
		TargetAddr:   lt.Addr,
		SourceRegion: LocalRegion,
		TargetRegion: lt.TargetRegion,
		ProbeType:    lt.ProbeType,
		TimeStamp:    time.Now().Unix(),
		Value:        float32(pkgRateNum),
	}

	prLaten := common.ProberResultOne{
		MetricName:   common.MetricsNamePingLatency,
		WorkerName:   LocalIp,
		TargetAddr:   lt.Addr,
		SourceRegion: LocalRegion,
		TargetRegion: lt.TargetRegion,
		ProbeType:    lt.ProbeType,
		TimeStamp:    time.Now().Unix(),
		Value:        float32(pingEwmaNum),
	}
	if pkgRateNum == 100 {
		prSu.Value = -1
	} else {
		prSu.Value = 1
	}

	prs = append(prs, &prSu)
	prs = append(prs, &prDr)
	prs = append(prs, &prLaten)
	//level.Info(lt.logger).Log("msg", "ping_res_prDr", "ts", prDr.TimeStamp, "value", prDr.Value)
	//level.Info(lt.logger).Log("msg", "ping_res_prLaten", "ts", prLaten.TimeStamp, "value", prLaten.Value)
	return prs
}

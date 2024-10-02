package common

import (
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetLocalIp() string {

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Printf("get local addr err:%v", err)
		return ""
	}
	localIp := strings.Split(conn.LocalAddr().String(), ":")[0]
	conn.Close()
	return localIp

}

func GetHostName() string {
	name, _ := os.Hostname()
	return name
}

func ShellCommand(shellStr string) (string, error) {
	// 3秒超时的ctx
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	cmd := exec.CommandContext(ctxt, "sh", "-c", shellStr)
	var buf bytes.Buffer
	// 标准错误重定向到标准输出
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return buf.String(), err
	}
	if err := cmd.Wait(); err != nil {
		return buf.String(), err
	}
	return buf.String(), nil

}

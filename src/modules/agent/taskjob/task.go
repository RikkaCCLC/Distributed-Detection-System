package taskjob

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/sys"
)

type Task struct {
	sync.Mutex

	Id      int64  //uid
	JobId   int64  //任务id
	Clock   int64  // 完成时间
	Action  string // 动作
	Status  string // 状态
	MetaDir string // 本地元信息目录

	alive  bool      //是否还在运行中的标志位
	Cmd    *exec.Cmd // cmd对象
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	// 元信息
	Script  string
	Args    string
	Account string
}

// 设置状态
func (t *Task) SetStatus(status string) {
	t.Lock()
	t.Status = status
	t.Unlock()
}

// 获取状态
func (t *Task) GetStatus() string {
	t.Lock()
	s := t.Status
	t.Unlock()
	return s
}

// 判断是否还在运行中
func (t *Task) GetAlive() bool {
	t.Lock()
	pa := t.alive
	t.Unlock()
	return pa
}

// 设置在运行中
func (t *Task) SetAlive(pa bool) {
	t.Lock()
	t.alive = pa
	t.Unlock()
}

// 获取标准输出
func (t *Task) GetStdout() string {
	t.Lock()
	out := t.Stdout.String()
	t.Unlock()
	return out
}

// 获取错误
func (t *Task) GetStderr() string {
	t.Lock()
	out := t.Stderr.String()
	t.Unlock()
	return out
}

// 重置错误和输出的buffer
func (t *Task) ResetBuff() {
	t.Lock()
	t.Stdout.Reset()
	t.Stderr.Reset()
	t.Unlock()
}

// 通过.done文件判断任务是否完成
func (t *Task) doneBefore() bool {
	doneFlag := path.Join(t.MetaDir, fmt.Sprint(t.Id), fmt.Sprintf("%d.done", t.Clock))
	return file.IsExist(doneFlag)
}

// 加载任务的记过
func (t *Task) loadResult() {
	metadir := t.MetaDir

	doneFlag := path.Join(metadir, fmt.Sprint(t.Id), fmt.Sprintf("%d.done", t.Clock))
	stdoutFile := path.Join(metadir, fmt.Sprint(t.Id), "stdout")
	stderrFile := path.Join(metadir, fmt.Sprint(t.Id), "stderr")

	var err error

	t.Status, err = file.ReadStringTrim(doneFlag)
	if err != nil {
		log.Printf("[E] read file %s fail %v", doneFlag, err)
	}
	stdout, err := file.ReadString(stdoutFile)
	if err != nil {
		log.Printf("[E] read file %s fail %v", stdoutFile, err)
	}
	stderr, err := file.ReadString(stderrFile)
	if err != nil {
		log.Printf("[E] read file %s fail %v", stderrFile, err)
	}

	t.Stdout = *bytes.NewBufferString(stdout)
	t.Stderr = *bytes.NewBufferString(stderr)
}
func (t *Task) meta() (script string, args string, account string) {
	return
}

// 准备任务的相关目录
func (t *Task) prepare() error {
	//if t.Account != "" {
	//	// already prepared
	//	return nil
	//}

	IdDir := path.Join(t.MetaDir, fmt.Sprint(t.Id))
	err := file.EnsureDir(IdDir)
	if err != nil {
		log.Printf("[E] mkdir -p %s fail: %v", IdDir, err)
		return err
	}

	writeFlag := path.Join(IdDir, ".write")
	if file.IsExist(writeFlag) {
		// 从磁盘读取
		argsFile := path.Join(IdDir, "args")
		args, err := file.ReadStringTrim(argsFile)
		if err != nil {
			log.Printf("[E] read %s fail %v", argsFile, err)
			return err
		}

		accountFile := path.Join(IdDir, "account")
		account, err := file.ReadStringTrim(accountFile)
		if err != nil {
			log.Printf("[E] read %s fail %v", accountFile, err)
			return err
		}

		t.Args = args
		t.Account = account
	} else {
		// 从远端读取，再写入磁盘
		script, args, account := t.Script, t.Args, t.Account

		scriptFile := path.Join(IdDir, "script")
		_, err = file.WriteString(scriptFile, script)
		if err != nil {
			log.Printf("[E] write script to %s fail: %v", scriptFile, err)
			return err
		}

		out, err := sys.CmdOutTrim("chmod", "+x", scriptFile)
		if err != nil {
			log.Printf("[E] chmod +x %s fail %v. output: %s", scriptFile, err, out)
			return err
		}

		argsFile := path.Join(IdDir, "args")
		_, err = file.WriteString(argsFile, args)
		if err != nil {
			log.Printf("[E] write args to %s fail: %v", argsFile, err)
			return err
		}

		accountFile := path.Join(IdDir, "account")
		_, err = file.WriteString(accountFile, account)
		if err != nil {
			log.Printf("[E] write account to %s fail: %v", accountFile, err)
			return err
		}

		_, err = file.WriteString(writeFlag, "")
		if err != nil {
			log.Printf("[E] create %s flag file fail: %v", writeFlag, err)
			return err
		}

		t.Args = args
		t.Account = account
	}

	return nil
}

// 启动任务
func (t *Task) start() {
	if t.GetAlive() {
		return
	}
	err := t.prepare()
	if err != nil {
		return
	}
	args := t.Args
	if args != "" {
		args = strings.Replace(args, ",,", "' '", -1)
		args = "'" + args + "'"
	}
	nowPath, _ := os.Getwd()
	scriptFile := path.Join(nowPath, t.MetaDir, fmt.Sprint(t.Id), "script")
	sh := fmt.Sprintf("%s %s", scriptFile, args)
	logger.Infof("[scriptFile:%+v][shCmd:%+v]", scriptFile, sh)
	var cmd *exec.Cmd
	if t.Account == "root" {
		cmd = exec.Command("sh", "-c", sh)
		cmd.Dir = "/root"
	} else {
		cmd = exec.Command("su", "-c", sh, "-", t.Account)
	}

	cmd.Stdout = &t.Stdout
	cmd.Stderr = &t.Stderr
	t.Cmd = cmd
	err = cmd.Start()
	if err != nil {
		log.Printf("[E] cannot start cmd of task[%d]: %v", t.Id, err)
		return
	}

	go runProcess(t)
}

func runProcess(t *Task) {
	t.SetAlive(true)
	defer t.SetAlive(false)

	err := t.Cmd.Wait()
	if err != nil {
		if strings.Contains(err.Error(), "signal: killed") {
			t.SetStatus("killed")
			logger.Debugf("process of task[%d] killed", t.Id)
		} else {
			t.SetStatus("failed")
			logger.Debugf("process of task[%d] return error: %v", t.Id, err)
		}
	} else {
		t.SetStatus("success")
		logger.Debugf("process of task[%d] done", t.Id)
	}

	persistResult(t)
}

// 任务结果落盘
func persistResult(t *Task) {
	metadir := t.MetaDir

	stdout := path.Join(metadir, fmt.Sprint(t.Id), "stdout")
	stderr := path.Join(metadir, fmt.Sprint(t.Id), "stderr")
	doneFlag := path.Join(metadir, fmt.Sprint(t.Id), fmt.Sprintf("%d.done", t.Clock))

	file.WriteString(stdout, t.GetStdout())
	file.WriteString(stderr, t.GetStderr())
	file.WriteString(doneFlag, t.GetStatus())
}

// kill任务
func (t *Task) kill() {
	go killProcess(t)
}

// 杀进程
func killProcess(t *Task) {
	t.SetAlive(true)
	defer t.SetAlive(false)

	logger.Debugf("begin kill process of task[%d]", t.Id)

	err := KillProcessByTaskID(t.Id, t.MetaDir)
	if err != nil {
		t.SetStatus("killfailed")
		logger.Debugf("kill process of task[%d] fail: %v", t.Id, err)
	} else {
		t.SetStatus("killed")
		logger.Debugf("process of task[%d] killed", t.Id)
	}

	persistResult(t)
}

func KillProcessByTaskID(id int64, metadir string) error {
	dir := strings.TrimRight(metadir, "/")
	arr := strings.Split(dir, "/")
	lst := arr[len(arr)-1]
	return sys.KillProcessByCmdline(fmt.Sprintf("%s/%d/script", lst, id))
}

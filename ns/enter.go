package ns

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"gitee.com/shangc1016/runc/cgroup"
	"gitee.com/shangc1016/runc/mount"
	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
	"github.com/sirupsen/logrus"
)

func ReenterConfig(it bool, mem, cpu string, command []string, volume []mount.VolumeInfo, name string) {
	self := EnterNs(command, it, volume)
	if err := self.Start(); err != nil {
		fmt.Println("<reenter error>", err)
		os.Exit(-1)
	}
	if it {
		self.Stdin = os.Stdin
		self.Stdout = os.Stdout
		self.Stderr = os.Stderr
	}
	fmt.Println("pid:", self.Process.Pid)

	// 设置cgroup
	cgroupLimit := cgroup.NewCgroupResource()
	cgroupLimit.SetCgroupPath("/sys/fs/cgroup")
	cgroupLimit.AddCgroupResource(*cgroup.NewCpuResource(cpu))
	cgroupLimit.AddCgroupResource(*cgroup.NewMemoryResource(mem))
	if err := cgroupLimit.InitNode(); err != nil {
		fmt.Println("initnode error", err)
		os.Exit(-1)
	}
	if err := cgroupLimit.SetQuota(); err != nil {
		fmt.Println("setquota error", err)
		os.Exit(-1)
	}
	if err := cgroupLimit.SetPid(self.Process.Pid); err != nil {
		fmt.Println("setpid error", err)
		os.Exit(-1)
	}
	//end cgroup

	if it {
		self.Wait()
		// 移除cgroup资源限制
		// cgroupLimit.RemoveNode()
		status, err := mount.RemoveVolmeMountPoint("/var/lib/runc", volume)
		fmt.Println(status, err)
		// 移除overlay分层文件系统
		mount.DeleteWorkSpace("/var/lib/runc")

	} else {
		// 后台运行，此时需要把子进程相关信息写入文件
		var childInfo status.ContainerInfo = status.ContainerInfo{
			Name:       name,
			Pid:        strconv.Itoa(self.Process.Pid),
			Status:     status.RUNNING,
			Id:         utils.GenRandString(10),
			Command:    command[0],
			CreateTime: time.Now().String(),
			Volumes:    volume,
		}
		data, _ := json.Marshal(childInfo)
		err := utils.WriteFile(data, "/var/lib/runc/status/"+childInfo.Id)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("write status ok, id:", childInfo.Id)
	}
	os.Exit(0)
	// 在处理完overlay挂载点等cleanup工作之后，退出
}

func EnterNs(args []string, it bool, volume []mount.VolumeInfo) *exec.Cmd {
	logrus.Info("config re-enter params")
	var argsInit []string = make([]string, len(args)+1)
	argsInit[0] = "init"
	for k, v := range args {
		argsInit[k+1] = v
	}
	cmd := exec.Command("/proc/self/exe", argsInit...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}
	if it {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// 在子进程启动之前就已经准备好了
	cmd.Dir = "/var/lib/runc/mnt"
	// 在进程运行之前准备好用到的分层的文件系统
	mount.NewWorkSpace("/var/lib/runc")
	mount.SetVolumeMountPoint("/var/lib/runc/mnt", volume)
	return cmd
}

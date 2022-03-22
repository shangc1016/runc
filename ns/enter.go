package ns

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"

	"gitee.com/shangc1016/runc/cgroup"
	"gitee.com/shangc1016/runc/fs"
	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
)

func ReenterConfig(it bool, mem, cpu string, command []string, volumes []fs.VolumeInfo, name string) {

	id := utils.GenRandString(10)
	self := EnterNs(command, it, id)
	// 新建文件系统
	fsys := fs.NewContainerFS(id)
	// 创建相应的目录结构
	fsys.Init()
	// 挂载基础镜像，设置overlay文件系统
	fsys.MkMountFs(path.Join(utils.Storage.Path, utils.Storage.Images, "busybox"))
	// 挂载数据卷
	fsys.MkVolumeFs(volumes)
	// 设置子进程的进入目录为/var/lib/runc/containers/`name`/mnt
	enterPoint := path.Join(utils.Storage.Path, utils.Storage.Containers, id, utils.Storage.Mnt)
	fmt.Println("emnterpoint:", enterPoint)
	self.Dir = enterPoint
	// 运行self.
	if err := self.Start(); err != nil {
		fmt.Println("<reenter error>", err)
		os.Exit(-1)
	}

	fmt.Println("pid:", self.Process.Pid)
	fmt.Println("id:", id)

	// 设置cgroup
	cgroupLimit := cgroup.NewCgroupResource(id, strconv.Itoa(self.Process.Pid))
	cgroupLimit.AddCgroupResource(cgroup.ResourceItem{
		Type: utils.Cgroup.Memory,
		File: utils.Cgroup.MemoryQuota, Quota: mem})

	cgroupLimit.AddCgroupResource(cgroup.ResourceItem{
		Type: utils.Cgroup.Cpu,
		File: utils.Cgroup.CpuQuota, Quota: cpu})

	cgroupLimit.Execute()

	if it {
		fmt.Println("waiting...")

		self.Wait()

		// 移除cgroup资源限制
		// cgroupLimit.RemoveNode()
		// status, err := mount.RemoveVolmeMountPoint("/var/lib/runc", volumes)
		// fmt.Println(status, err)
		// 移除overlay分层文件系统
		// mount.DeleteWorkSpace("/var/lib/runc")

	} else {
		// 后台运行，此时需要把子进程相关信息写入文件
		var childInfo status.ContainerInfo = status.ContainerInfo{
			Name:       name,
			Pid:        strconv.Itoa(self.Process.Pid),
			Status:     status.RUNNING,
			Id:         id,
			Command:    command[0],
			CreateTime: time.Now().String(),
			Volumes:    volumes,
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

func EnterNs(args []string, it bool, id string) *exec.Cmd {
	var argsInit []string = make([]string, len(args)+1)
	argsInit[0] = "init"
	// argsInit[1] = id
	for k, v := range args {
		argsInit[k+1] = v
	}
	cmd := exec.Command("/proc/self/exe", argsInit...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	}
	if it {
		// fmt.Println("it:", it)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// // 在子进程启动之前就已经准备好了
	// cmd.Dir = "/var/lib/runc/mnt"
	// // 在进程运行之前准备好用到的分层的文件系统
	// mount.NewWorkSpace("/var/lib/runc")
	// mount.SetVolumeMountPoint("/var/lib/runc/mnt", volume)
	return cmd
}

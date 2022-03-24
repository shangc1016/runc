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
	"gitee.com/shangc1016/runc/fsys"
	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
)

func ReenterConfig(it bool, mem, cpu string, command []string, volumes []fsys.VolumeInfo, name string) {

	id := utils.GenRandString(10)
	self := EnterNs(command, it, path.Join(utils.Storage.Path, utils.Storage.Containers, id, "mnt"))

	// 新建文件系统
	fsys := fsys.NewContainerFS(id)
	// 创建相应的目录结构
	fsys.Init()
	// 挂载基础镜像，设置overlay文件系统
	fsys.MkMountFs(path.Join(utils.Storage.Path, utils.Storage.Images, "busybox"))
	// 挂载数据卷
	fsys.MkVolumeFs(volumes)

	if !it {
		// 如果容器以后台进程的方式运行，需要把容器进程的标准输出重定向到/var/lib/runc/containers/`name`/output/output.log 文件中
		f, err := os.OpenFile(path.Join(fsys.Path, fsys.Name, fsys.Output, "output.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			fmt.Println("<output redirect error>", err)
		}
		defer f.Close()
		self.Stdout = f
		self.Stderr = f
	}
	// 设置子进程的进入目录为/var/lib/runc/containers/`name`/mnt
	enterPoint := path.Join(utils.Storage.Path, utils.Storage.Containers, id, utils.Storage.Mnt)
	self.Dir = enterPoint
	// 运行self.
	if err := self.Start(); err != nil {
		fmt.Println("<re-enter error>", err)
		os.Exit(-1)
	}

	fmt.Printf("pid:%v, id:%v\n", self.Process.Pid, id)
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
		//  it模式: 等待子进程返回
		self.Wait()
		cgroupLimit.RemoveNode(id)
		fsys.CleanUp()

	} else {
		// 后台运行，此时需要把子进程相关信息写入文件 /var/lib/runc/status/`id`
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
		err := utils.WriteFile(data, path.Join(utils.Storage.Path, utils.Storage.Status, childInfo.Id))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("write status ok, id:", childInfo.Id)
	}
	os.Exit(0)
	// 在处理完overlay挂载点等cleanup工作之后，退出
}

func EnterNs(args []string, it bool, rootPath string) *exec.Cmd {

	var argsInit []string = make([]string, len(args)+2)
	argsInit[0] = "init"
	argsInit[1] = rootPath
	for k, v := range args {
		argsInit[k+2] = v
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
	return cmd
}

package checker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"gitee.com/shangc1016/runc/cgroup"
	"gitee.com/shangc1016/runc/fsys"
	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
)

// check  模块检查/var/lib/runc/containers 目录下的所有容器进程，

func RollContainers(fn func(info status.ContainerInfo)) {
	files, err := utils.GetFileNameAll(path.Join(utils.Storage.Path, utils.Storage.Status))
	if err != nil {
		fmt.Println("<get container error>", err)
		fmt.Println(path.Join(utils.Storage.Path, utils.Storage.Status))
		os.Exit(-1)
	}
	for _, file := range files {
		data, _ := utils.ReadFile(path.Join(utils.Storage.Path, utils.Storage.Status, file))
		var info status.ContainerInfo
		json.Unmarshal(data, &info)
		fn(info)
	}
}

func PrintInfo(info status.ContainerInfo) {
	fmt.Println(info)
}

func DeleteTerminated(info status.ContainerInfo) {
	if info.Status == status.TERMINATED {
		os.Remove(path.Join(utils.Storage.Path, utils.Storage.Status, info.Id))
	}
}

func ChangeStatus(id, now, then string) func(status.ContainerInfo) {
	return func(info status.ContainerInfo) {
		if info.Id == id && info.Status == now {
			info.Status = then
			data, _ := json.Marshal(info)
			// 重置状态
			f, _ := os.OpenFile(path.Join(utils.Storage.Path, utils.Storage.Status, info.Id), os.O_TRUNC|os.O_WRONLY, os.ModeAppend|os.ModePerm)
			defer f.Close()
			f.Write(data)
			// 杀死进程

			fmt.Println("kill", info.Pid)
			cmd := exec.Command("kill", info.Pid)
			if err := cmd.Run(); err != nil {
				fmt.Println("stop container error")
			}
		}
	}
}

// used by rm command
func RemoveById(id string) func(status.ContainerInfo) {
	return func(info status.ContainerInfo) {
		fmt.Println("id:", id, "info.id:", info.Id)
		if info.Id == id && info.Status == status.TERMINATED {
			os.Remove(path.Join(utils.Storage.Path, utils.Storage.Status, info.Id))
			fmt.Println("remove success")
		}
		// 移除文件系统相关
		// 移除cgroup
		// FIXME: 这儿设计的有问题，因为是命令行，不能把原来的对象传过来，但是在此处新建对象太蠢了。新建的空对象没有cgroup的相关限制，所以删除不了cgroup配额
		cgroupLimit := cgroup.NewCgroupResource(info.Id, info.Pid)
		cgroupLimit.RemoveNode()
		// 移除layer fs
		// FIXME:此处同理
		fsys := fsys.NewContainerFS(info.Id)
		fsys.Volumes = info.Volumes
		fsys.CleanUp()
	}
}

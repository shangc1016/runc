package checker

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

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

func ChangeStatus(now, then string, info status.ContainerInfo) {
	if info.Status == now {
		info.Status = then
	}
}

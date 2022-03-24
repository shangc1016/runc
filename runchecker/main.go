package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"gitee.com/shangc1016/runc/checker"
	"gitee.com/shangc1016/runc/fsys"
	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
)

// runchecker 是一个后台进程，不断检索容器目录/var/lib/runc/container,并且更新他们的状态， 以及对进程挂载的文件系统的清理
// 更新状态基本只有一种，从running到terminated
// 然后umount这个进程用到的所有...

// 自定义的函数传入, 判断进程是否存在，不存在则改变info.Status 状态
func checkExistChangeStatus(info status.ContainerInfo) {
	exist, err := fsys.PathExists(path.Join("/proc", info.Pid))
	if err != nil {
		fmt.Println("<runchecker error>", err)
		os.Exit(-1)
	}
	if info.Status == status.TERMINATED {
		return
	}
	if !exist {
		info.Status = status.TERMINATED
		data, _ := json.Marshal(info)
		// 重新写入status
		f, _ := os.OpenFile(path.Join(utils.Storage.Path, utils.Storage.Status, info.Id), os.O_TRUNC|os.O_WRONLY, os.ModeAppend|os.ModePerm)
		defer f.Close()
		f.Write(data)
		log, err := os.OpenFile(path.Join(utils.Storage.Path, utils.Storage.Logs, utils.RunChecker.LogName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend|os.ModePerm)
		if err != nil {
			fmt.Println("err", err)
		}
		defer log.Close()
		str := "id:" + info.Id + ", name:" + info.Name + " terminated.\n"
		log.Write([]byte(str))
	}

}

func main() {
	for {
		checker.RollContainers(checkExistChangeStatus)
		time.Sleep(time.Second * 1)
	}
}

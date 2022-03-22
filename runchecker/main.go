package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"gitee.com/shangc1016/runc/mount"
	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
)

// runchecker 是一个后台进程，不断检索容器目录/var/lib/runc/container,并且更新他们的状态

func main() {
	for {
		files, err := utils.GetFileNameAll("/var/lib/runc/status")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(files)
		for _, file := range files {
			// 读状态文件到data
			data, _ := utils.ReadFile("/var/lib/runc/status/" + file)
			var info status.ContainerInfo
			// 转为结构体info
			json.Unmarshal(data, &info)
			// 如果状态为terminated，continue
			// fmt.Println("info:", info)
			if info.Status == status.TERMINATED {
				continue
			}
			exist, err := utils.PathExists(path.Join("/proc/" + info.Pid))
			if err != nil {
				fmt.Println(err)
			}
			// 如果进程终止，改变status为terminated
			if !exist {
				// TODO
				// 移除挂载的文件系统，
				// state, err := mount.RemoveVolmeMountPoint("/var/lib/runc", info.Volumes)
				// fmt.Println(state, err)
				// 移除overlay分层文件系统
				mount.DeleteWorkSpace("/var/lib/runc")

				info.Status = status.TERMINATED
				data, _ := json.Marshal(info)
				utils.WriteFile(data, "/var/lib/runc/status/"+file)
				fmt.Println("container:", info.Id, "terminated.")
				loginfo := info.Id + " terminated " + time.Now().Format("2006-01-02 15:04:05")
				f, err := os.OpenFile("/var/lib/runc/log/runchecker.log", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("open error:", err)
				}
				io.WriteString(f, loginfo+"\n")
				f.Close()
			}
		}
		time.Sleep(time.Second * 1)
	}
}

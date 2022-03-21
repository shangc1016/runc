package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gitee.com/shangc1016/runc/mount"
	"gitee.com/shangc1016/runc/utils"
)

type ContainerInfo struct {
	Name       string             `json:"name"`
	Id         string             `json:"id"`
	Pid        string             `json:"pid"`
	Status     string             `json:"status"`
	Command    string             `json:"command"`
	CreateTime string             `json:"create_time"`
	Volumes    []mount.VolumeInfo `json:"volumes"`
}

var (
	// 容器运行的状态
	RUNNING    string = "running"
	TERMINATED string = "terminated"
)

// 删除状态为terminated的容器
func RemoveTerminated(statePath string) error {
	var terminated_arr []string
	info, _ := GetStateInfo(statePath)
	for _, item := range info {
		if item.Status == TERMINATED {
			terminated_arr = append(terminated_arr, item.Id)
		}
	}
	for _, item := range terminated_arr {
		os.Remove(path.Join(statePath, item))
	}
	return nil
}

// 得到状态为running的所有信息
func GetAllQuietStatus(statePath string) ([]string, error) {
	var status []string
	info, _ := GetStateInfo(statePath)
	for _, item := range info {
		status = append(status, item.Id)
	}
	return status, nil
}

// 得到状态为running的信息， 只包括id
func GetQuietStatus(statePath string) ([]string, error) {
	var status []string
	info, _ := GetStateInfo(statePath)
	for _, item := range info {
		if item.Status == "running" {
			status = append(status, item.Id)
		}
	}
	return status, nil
}

// 得到所有状态信息,无论running或者terminated；
func GetAllStatus(statePath string) ([]ContainerInfo, error) {
	return GetStateInfo(statePath)
}

//
func GetStatus(statePath string) ([]ContainerInfo, error) {
	var state []ContainerInfo
	info, _ := GetStateInfo(statePath)
	for _, item := range info {
		if item.Status == "running" {
			state = append(state, item)
		}
	}
	return state, nil
}

// 底层函数
func GetStateInfo(statePath string) ([]ContainerInfo, error) {
	var Info []ContainerInfo
	var tmp ContainerInfo
	dirs, err := utils.GetFileNameAll(statePath)
	if err != nil {
		fmt.Println("read all file error", err)
	}
	for _, dir := range dirs {
		f, err := os.OpenFile(path.Join(statePath, dir), os.O_RDWR, os.ModeAppend)
		if err != nil {
			fmt.Println("open file error", err)
		}
		data, _ := ioutil.ReadAll(f)
		json.Unmarshal(data, &tmp)
		Info = append(Info, tmp)
	}
	return Info, nil
}

package fs

import (
	"fmt"
	"os"
	"path"

	"gitee.com/shangc1016/runc/utils"
)

func InitFs() {
	fmt.Println("init fs(storage, cgroup)...")

	status := InitializeStorageDir(utils.Storage.Path,
		[]string{utils.Storage.Containers,
			utils.Storage.Images, utils.Storage.Logs, utils.Storage.Status})
	if !status {
		fmt.Println("initialize storage error")
		os.Exit(-1)
	}

	status = InitializeCgroup(utils.Cgroup.Path, utils.Project.Name,
		[]string{utils.Cgroup.Memory, utils.Cgroup.Cpu})
	if !status {
		fmt.Println("initialize cgroup error")
		os.Exit(-1)
	}
}

func InitializeStorageDir(storagePath string, dirs []string) bool {
	for _, dir := range dirs {
		status := MkDirIfNotExist(path.Join(storagePath, dir))
		if !status {
			return false
		}
	}
	return true
}

func InitializeCgroup(cgroupPath, name string, resources []string) bool {
	for _, resource := range resources {
		status := MkDirIfNotExist(path.Join(cgroupPath, resource, name))
		if !status {
			return false
		}
	}
	return true
}

// 创建运行容器的目录,name 必须互不相同
func MkContainerDir(containerPath, name string, dirs []string) bool {
	container := path.Join(containerPath, name)
	exist, _ := PathExists(container)
	if exist {
		fmt.Printf("container name: `%v` already exist\n", name)
		return false
	}
	status := MkDirIfNotExist(container)
	if !status {
		return false
	}
	for _, dir := range dirs {
		status := MkDirIfNotExist(path.Join(container, dir))
		if !status {
			return false
		}
	}
	return true
}

func MkDirIfNotExist(dir string) bool {
	status, err := PathExists(dir)
	if err != nil {
		fmt.Println("error")
	}
	if status {
		return true
	}
	err = os.Mkdir(dir, 0644)
	return err == nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

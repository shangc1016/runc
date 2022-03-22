package fs

import (
	"fmt"
	"os"
	"path"
)

func MkFs(fsPath string) error {
	return nil
}

func InitializeStorageDir(cp, ip, sp, lp string) bool {
	// 创建存储相关目录
	s1, s2, s3, s4 := MkDirIfNotExist(cp), MkDirIfNotExist(ip), MkDirIfNotExist(sp), MkDirIfNotExist(lp)
	if s1 || s2 || s3 || s4 {
		return false
	}
	return true
}

func InitializeCgroupDir(name, cgroupPath string, resources []string) bool {
	for _, resource := range resources {
		status := MkDirIfNotExist(path.Join(cgroupPath, resource, name))
		if !status {
			return false
		}
	}
	return true
}

// 创建运行容器的目录
func MkContainerDir(containerPath, name, log string, dirs []string) bool {
	container := path.Join(containerPath, name)
	status := MkDirIfNotExist(container)
	if !status {
		return false
	}
	f, _ := os.Create(path.Join(container, log))
	f.Close()
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

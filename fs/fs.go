package fs

import (
	"fmt"
	"os"
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

func InitializeCgroupDir(name string, resourcesPath []string) bool {
	return false
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

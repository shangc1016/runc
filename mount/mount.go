package mount

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"gitee.com/shangc1016/runc/utils"
	"github.com/sirupsen/logrus"
)

// 使用overlay文件系统作为底层分出呢个文件系统的实现

// 需要准备好overlay文件系统的lowerdir、upperdir、workdir以及merged目录
// 分别将目录设置为：
// lowerdir: /var/lib/runc/busybox    (作为只读层)
// upperdir: /var/lib/runc/upper
// workdir: /var/lib/runc/work
// merged: /var/lib/runc/mnt

func NewWorkSpace(rootURL string) {
	status, err := CreateLayer(path.Join(rootURL, "busybox"))
	if err != nil {
		logrus.Errorf(err.Error())
		fmt.Println(status)
	}
	lowerPath := path.Join(rootURL, "busybox")
	upperPath := path.Join(rootURL, "upper")
	workPath := path.Join(rootURL, "work")
	mntPath := path.Join(rootURL, "mnt")
	CreateLayer(upperPath)
	CreateLayer(workPath)
	CreateLayer(mntPath)

	_, err = CreateMountPoint(lowerPath, upperPath, workPath, mntPath)
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteWorkSpace(rootURL string) {
	_, err := DeleteMountPoint(path.Join(rootURL, "mnt"))
	if err != nil {
		fmt.Println(err)
	}
	_, err = DeleteFile(path.Join(rootURL, "upper"))
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteMountPoint(mountPath string) (bool, error) {
	fmt.Println("delete mountpoint:", mountPath)
	exist, err := utils.PathExists(mountPath)
	if err != nil {
		return false, err
	}
	if exist {
		cmd := exec.Command("umount", mountPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return false, err
		}
		return true, nil
	}
	return true, nil
}

// 删除可写层
func DeleteFile(writablePath string) (bool, error) {
	fmt.Println("delete file:", writablePath)
	exist, err := utils.PathExists(writablePath)
	if err != nil {
		return false, err
	}
	if exist {
		// 级联删除所有此目录
		if err := os.RemoveAll(writablePath); err != nil {
			return false, err
		}
		return true, nil
	}
	return true, nil
}

func CreateLayer(path string) (bool, error) {
	exist, err := utils.PathExists(path)
	if err != nil {
		return false, err
	}
	if exist {
		return true, nil
	} else {
		if err := os.Mkdir(path, 0644); err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

// 挂载overlay的命令 eg:
// mount -t overlay overlay -o lowerdir=./lower/ -o upperdir=./upper/ -o workdir=./workdir/ ./merged/
func CreateMountPoint(lower, upper, work, merged string) (bool, error) {
	var args []string
	args = append(args, "-t")
	args = append(args, "overlay")
	args = append(args, "overlay")
	args = append(args, "-o")
	args = append(args, "lowerdir="+lower)
	args = append(args, "-o")
	args = append(args, "upperdir="+upper)
	args = append(args, "-o")
	args = append(args, "workdir="+work)
	args = append(args, merged)

	cmd := exec.Command("mount", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return false, err
	}
	return true, nil
}

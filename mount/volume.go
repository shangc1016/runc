package mount

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"gitee.com/shangc1016/runc/utils"
)

type VolumeInfo struct {
	OuterPath string `json:"outer_path"`
	InnerPath string `json:"inner_path"`
}

// volume的挂载放在分层的文件系统建立好了之后，挂载到最终的分层文件系统之上，从而实现容器内外通过volume互通。

// args是数组，每一项形如'/out/path:/inner/path'
func ParseVolume(args []string) ([]VolumeInfo, bool) {
	volumes := make([]VolumeInfo, len(args))
	for k, arg := range args {
		twoParts := strings.Split(arg, ":")
		if len(twoParts) != 2 || len(twoParts[0]) == 0 || len(twoParts[1]) == 0 {
			return []VolumeInfo{}, false
		}
		volumes[k] = VolumeInfo{
			OuterPath: twoParts[0],
			InnerPath: twoParts[1],
		}
	}
	return volumes, true
}

func SetVolumeMountPoint(mntPath string, volumes []VolumeInfo) (bool, error) {
	exist, err := utils.PathExists(mntPath)
	if err != nil {
		fmt.Println("err", err)
	}
	if !exist {
		fmt.Println("not exist")
		os.Exit(-1)
	}
	for _, volume := range volumes {
		inner := path.Join(mntPath, volume.InnerPath)
		outer := volume.OuterPath
		// 内部挂载点不存在就创建相应的文件，
		if exist, _ := utils.PathExists(inner); !exist {
			os.Mkdir(inner, 0644)
		}
		// 外部挂载点比存在就报错退出
		fmt.Println(outer)
		fmt.Println(inner)
		if exist, _ := utils.PathExists(outer); !exist {
			fmt.Println("out not exist")
			os.Exit(-1)
		}
		cmd := exec.Command("mount", "--bind", outer, inner)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("create mount point error", err)
			return false, err
		}
	}
	return true, nil
}

func RemoveVolmeMountPoint(mntPath string, volumes []VolumeInfo) (bool, error) {
	exist, err := utils.PathExists(mntPath)
	if err != nil {
		fmt.Println("path exist error", err)
		return false, err
	}
	if !exist {
		fmt.Println("file not exist")
		return false, nil
	}
	for _, volume := range volumes {
		inner := path.Join(mntPath, volume.InnerPath)
		cmd := exec.Command("umount", inner)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return false, err
		}
	}
	return true, nil
}

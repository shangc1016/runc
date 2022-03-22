package fs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"gitee.com/shangc1016/runc/utils"
)

type VolumeInfo struct {
	OuterPath string `json:"outerpath"`
	InnerPath string `json:"innerpath"`
}

type Containerfs struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Mnt      string `json:"container"`
	Upperdir string `json:"upperdir"`
	Workdir  string `json:"workdir"`
	Output   string `json:"output"`
}

func NewContainerFS(name string) *Containerfs {
	return &Containerfs{
		Name:     name,
		Path:     path.Join(utils.Storage.Path, utils.Storage.Containers),
		Mnt:      utils.Storage.Mnt,
		Upperdir: utils.Storage.UpperDir,
		Workdir:  utils.Storage.WorkDir,
		Output:   utils.Storage.Output,
	}
}

func (c *Containerfs) Init() {
	containerPath := path.Join(c.Path, c.Name)
	exist, _ := PathExists(containerPath)
	if exist {
		fmt.Println("container name already exist")
		os.Exit(-1)
	}
	MkDirIfNotExist(containerPath)
	MkDirIfNotExist(path.Join(containerPath, c.Mnt))
	MkDirIfNotExist(path.Join(containerPath, c.Upperdir))
	MkDirIfNotExist(path.Join(containerPath, c.Workdir))
	MkDirIfNotExist(path.Join(containerPath, c.Output))
}

func (c *Containerfs) MkMountFs(imagePath string) {
	// 把imagePath 挂载到当前容器目录的mnt下， 使用overlay文件系统挂载
	cmd := exec.Command("mount", "-t", "overlay", "overlay",
		"-o", "lowerdir="+imagePath,
		"-o", "upperdir="+path.Join(c.Path, c.Name, c.Upperdir),
		"-o", "workdir="+path.Join(c.Path, c.Name, c.Workdir),
		path.Join(c.Path, c.Name, c.Mnt))
	if err := cmd.Run(); err != nil {
		fmt.Println("mount error", err.Error())
		os.Exit(-1)
	}
}

func (c *Containerfs) MkVolumeFs(volumes []VolumeInfo) bool {
	for _, volume := range volumes {
		exist, _ := PathExists(path.Join(c.Mnt, volume.InnerPath))
		if !exist {
			os.Mkdir(path.Join(c.Mnt, volume.InnerPath), 0644)
		}
		exist, _ = PathExists(volume.OuterPath)
		if !exist {
			fmt.Println("outter path not exist")
			return false
		}
		cmd := exec.Command("mount", "--bind", volume.OuterPath, path.Join(c.Mnt, volume.InnerPath))
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	return true
}

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

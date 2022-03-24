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
	Name     string       `json:"name"`      // eg:"qwertyuiop"
	Path     string       `json:"path"`      // eg:"/var/lib/runc/containers"
	Mnt      string       `json:"container"` // eg:"mnt"
	Upperdir string       `json:"upperdir"`  // eg:"ipperdir"
	Workdir  string       `json:"workdir"`   // eg:"workdir"
	Output   string       `json:"output"`    // eg:"output"
	Volumes  []VolumeInfo `json:"volume"`    // eg:"[{/outer/path:/inner/path}...]"
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
	//  先判断imagePath镜像是否存在
	exist, _ := PathExists(imagePath)
	if !exist {
		fmt.Println("error, imagePath not exist.")
		os.Exit(-1)
	}
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

func (c *Containerfs) MkVolumeFs(volumes []VolumeInfo) {
	c.Volumes = volumes
	for _, volume := range volumes {
		fmt.Println("inner:", volume.InnerPath)
		fmt.Println("outer:", volume.OuterPath)
		inner := path.Join(c.Path, c.Name, c.Mnt, volume.InnerPath)
		fmt.Println("inner:", inner)
		exist, _ := PathExists(inner)
		if !exist {
			os.MkdirAll(inner, 0644)
		}
		exist, _ = PathExists(volume.OuterPath)
		if !exist {
			fmt.Println("挂载：外部路径不存在")
			os.Exit(-1)
		}
		cmd := exec.Command("mount", "--bind", volume.OuterPath, inner)
		if err := cmd.Run(); err != nil {
			fmt.Println("mount error", err)
			os.Exit(-1)
		}
	}
}

func (c *Containerfs) CleanUp() {
	c.cleanUpVolume()
	c.cleanUpLayer()
	err := os.RemoveAll(path.Join(c.Path, c.Name))
	if err != nil {
		fmt.Println("<remove container dir err>", err)
	}
	fmt.Println("cleanup fsys")
}

func (c *Containerfs) cleanUpVolume() bool {
	// cleanup all volume mountpoint.
	for _, volume := range c.Volumes {
		cmd := exec.Command("umount", path.Join(c.Path, c.Name, c.Mnt, volume.InnerPath))
		if err := cmd.Run(); err != nil {
			fmt.Println("<remove volume mount error>", err)
			return false
		}
	}
	return true
}

func (c *Containerfs) cleanUpLayer() bool {
	cmd := exec.Command("umount", "-l", path.Join(c.Path, c.Name, c.Mnt))
	if err := cmd.Run(); err != nil {
		fmt.Println("<remove layerfs mount>", err)
		return false
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

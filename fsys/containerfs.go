package fsys

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

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
		inner := path.Join(c.Path, c.Name, c.Mnt, volume.InnerPath)
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
	fmt.Println("fsys cleaned")
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

func PivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir(root); err != nil {
		return fmt.Errorf("chdir %v %v", root, err)
	}

	pivotDir = filepath.Join(root, ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}

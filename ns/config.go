package ns

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/sirupsen/logrus"
)

func Config(root string, command []string) {
	// 首先chroot到/var/lib/runc/mnt这个overlay的文件系统中

	fmt.Println("root::::", root)
	if err := syscall.Chroot(root); err != nil {
		fmt.Println("<chroot error>.", err.Error())
		os.Exit(-1)
	}

	// fsys.PivotRoot(root)

	// 然后把容器的proc文件系统挂载到 /proc目录
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NODEV | syscall.MS_NOSUID | syscall.MS_REC
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		fmt.Println(err)
	}
	err = syscall.Mount("tmpfs	", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		fmt.Println(err)
	}

	path, err := exec.LookPath(command[0])
	if err != nil {
		fmt.Println("<path not found>.", err.Error())
		os.Exit(-1)
	}
	if err := syscall.Exec(path, command[0:], os.Environ()); err != nil {
		pc, file, line, _ := runtime.Caller(1)
		logrus.Fatal(pc, file, line)
	}
}

package cmd

import (
	"fmt"
	"os"

	"gitee.com/shangc1016/runc/checker"
	"gitee.com/shangc1016/runc/formater"
	"gitee.com/shangc1016/runc/fsys"
	"gitee.com/shangc1016/runc/mount"
	"gitee.com/shangc1016/runc/ns"
	"gitee.com/shangc1016/runc/status"

	"github.com/spf13/cobra"
)

var (
	memQuota string
	cpuQuota string
	volume   []string
	name     string
	it       bool
	detach   bool
	all      bool
	quiet    bool
)

type RunFlags struct {
	Name     string
	It       bool
	Detach   bool
	Volumes  []mount.VolumeInfo
	MemQuota string
	CpuQuota string
}

type PsFlags struct {
	All   bool
	Quiet bool
}

var rootCmd *cobra.Command = &cobra.Command{
	Use: "runc",
}

var runCmd *cobra.Command = &cobra.Command{
	Use:   "run",
	Short: "execute process",
	Run: func(cmd *cobra.Command, args []string) {

		volumes, status := fsys.ParseVolume(volume)
		if !status {
			fmt.Println("volume format error(exit)")
			os.Exit(-1)
		}
		if it && detach || !it && !detach {
			// 要么是后台运行，要么是it模拟终端运行，不能都是
			fmt.Println("must run in it or detach mode")
			os.Exit(-1)
		}
		ns.ReenterConfig(it, memQuota, cpuQuota, args, volumes, name)
	},
}

var commitCmd *cobra.Command = &cobra.Command{
	Use:   "commit",
	Short: "commit container filesystem to achieved",
	Run: func(cmd *cobra.Command, args []string) {

		// TODO
		fmt.Println("args:", args)
		if len(args) == 0 {
			fmt.Println("error")
		}
		fmt.Println("commit container filesystem to achieved")
	},
}

// 只有在rm命令可以删除后台进程的文件系统
// it模式的进程在退出时删除文件系统
var rmCmd *cobra.Command = &cobra.Command{
	Use:   "rm",
	Short: "remove stopped container file system",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("args error")
			os.Exit(-1)
		}
		checker.RollContainers(checker.RemoveById(args[0]))
	},
}

// 杀死容器进程，设置状态为terminated
var killCmd *cobra.Command = &cobra.Command{
	Use:   "kill",
	Short: "kill running container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("args error")
			os.Exit(-1)
		}
		fmt.Println("stop:", args[0])
		checker.RollContainers(checker.ChangeStatus(args[0], status.RUNNING, status.TERMINATED))
	},
}

var psCmd *cobra.Command = &cobra.Command{
	Use:   "ps",
	Short: "print state of all container",
	Run: func(cmd *cobra.Command, args []string) {

		if all && quiet {
			info, _ := status.GetAllQuietStatus("/var/lib/runc/status")
			formater.PsQuiet(info)
		} else if quiet {
			info, _ := status.GetQuietStatus("/var/lib/runc/status")
			formater.PsQuiet(info)
		} else if all {
			info, _ := status.GetAllStatus("/var/lib/runc/status")
			formater.PsAll(info)
		} else {
			// 没有任何选项，打印正在运行的所有信息
			info, _ := status.GetStatus("/var/lib/runc/status")
			formater.PsAll(info)
		}
	},
}

// this sub-command should not invoke outside.
var initCmd *cobra.Command = &cobra.Command{
	Use:   "init",
	Short: "re-enter program and do some setting work.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:参数有效性验证
		ns.Config(args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(killCmd)

	runCmd.Flags().StringVarP(&memQuota, "mem", "m", "100m", "mem quota, range [...]")
	runCmd.Flags().StringVarP(&cpuQuota, "cpu", "c", "-1", "cpu quota, range [-1, 100000]")

	runCmd.Flags().StringVarP(&name, "name", "n", "none", "designate container name")

	//TODO: 只支持目录的挂载，不支持文件挂载
	runCmd.Flags().StringSliceVarP(&volume, "volume", "v", []string{}, "mount volume")

	runCmd.Flags().BoolVar(&it, "it", false, "with interactive termainl")
	runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "run background")

	psCmd.Flags().BoolVarP(&all, "all", "a", false, "print all container, no matter running or terminated.")
	psCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only print container's id.")

	initCmd.Hidden = true // only invoke internal.

}

func NewCmd() *cobra.Command {
	return rootCmd
}

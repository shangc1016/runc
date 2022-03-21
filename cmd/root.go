package cmd

import (
	"fmt"
	"os"

	"gitee.com/shangc1016/runc/formater"
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

var runFlags RunFlags
var psFlags PsFlags

var rootCmd *cobra.Command = &cobra.Command{
	Use: "run-c",
}

var runCmd *cobra.Command = &cobra.Command{
	Use:   "run",
	Short: "execute run-c process",
	Run: func(cmd *cobra.Command, args []string) {
		volume, status := mount.ParseVolume(volume)
		if !status {
			fmt.Println("volume format error(exit)")
			os.Exit(0)
		}
		if it && detach {
			// 要么是后台运行，要么是it模拟终端运行，不能都是
			fmt.Println("can not appear together")
			os.Exit(-1)
		}
		if detach {
			it = false
		}
		ns.ReenterConfig(it, memQuota, cpuQuota, args, volume, name)
	},
}

var commitCmd *cobra.Command = &cobra.Command{
	Use:   "commit",
	Short: "commit container filesystem to achieved",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args:", args)
		if len(args) == 0 {
			fmt.Println("error")
		}
		fmt.Println("commit container filesystem to achieved")
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
		ns.Config(args)
	},
}

func init() {

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&memQuota, "mem", "m", "100m", "mem quota, range [...]")
	runCmd.Flags().StringVarP(&cpuQuota, "cpu", "c", "-1", "cpu quota, range [-1, 100000]")

	runCmd.Flags().StringVarP(&name, "name", "n", "none", "designate container name")

	//TODO: 只支持目录的挂载，不支持文件挂载
	runCmd.Flags().StringSliceVarP(&volume, "volume", "v", []string{}, "mount volume")

	runCmd.Flags().BoolVar(&it, "it", false, "with interactive termainl")

	// 后台运行
	runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "run background")

	rootCmd.AddCommand(initCmd)
	initCmd.Hidden = true // only invoke internal.

	rootCmd.AddCommand(commitCmd)

	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVarP(&all, "all", "a", false, "print all container, no matter running or terminated.")
	psCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only print container's id.")

}

func NewCmd() *cobra.Command {
	return rootCmd
}

package main

import (
	"os"

	"gitee.com/shangc1016/runc/cmd"
	"gitee.com/shangc1016/runc/fsys"
	"github.com/sirupsen/logrus"
)

func main() {

	if os.Args[0] != "/proc/self/exe" {
		fsys.InitFs()
	}
	// fmt.Println(os.Args)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	rootCmd := cmd.NewCmd()
	rootCmd.Execute()
}

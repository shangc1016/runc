package main

import (
	"os"

	"gitee.com/shangc1016/runc/cmd"
	"gitee.com/shangc1016/runc/fs"
	"gitee.com/shangc1016/runc/utils"

	"github.com/sirupsen/logrus"
)

func main() {

	// fmt.Println("reenter1=========")
	// if os.Args[0] != "/proc/self/exe" {
	utils.InitConfig()
	cmd.InitCmd()
	fs.InitFs()

	// fmt.Println("reenter2=========")
	// fmt.Println("os.Args:", os.Args)
	// init logger config
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// init cli parser
	rootCMd := cmd.NewCmd()
	rootCMd.Execute()
}

package main

import (
	"os"
	"path"

	"gitee.com/shangc1016/runc/cmd"
	"gitee.com/shangc1016/runc/utils"
	"github.com/sirupsen/logrus"
)

func main() {

	// loading global config...
	if os.Args[0] != "/proc/self/exe" {
		wd, _ := os.Getwd()
		configPath := path.Join(wd, "config.json")
		utils.InitConfig(configPath)
	}

	// init logger config
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// init cli parser
	rootCMd := cmd.NewCmd()
	rootCMd.Execute()
}

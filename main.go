package main

import (
	"os"

	"gitee.com/shangc1016/runc/cmd"
	"gitee.com/shangc1016/runc/utils"
	"github.com/sirupsen/logrus"
)

func main() {

	// loading global config...
	utils.Init("config.json")

	// init logger config
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// init cli parser
	rootCMd := cmd.NewCmd()
	rootCMd.Execute()
}

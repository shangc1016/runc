package main

import (
	"os"

	"gitee.com/shangc1016/runc/cmd"
	"github.com/sirupsen/logrus"
)

func main() {

	// fmt.Println(os.Args)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	rootCmd := cmd.NewCmd()
	rootCmd.Execute()
}

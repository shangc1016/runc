package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type StorageConfig struct {
	Path       string `json:"path"`
	Status     string `json:"status"`
	Images     string `json:"images"`
	Logs       string `json:"logs"`
	Containers string `json:"containers"`
	Mnt        string `json:"mnt"`
	UpperDir   string `json:"upperdir"`
	WorkDir    string `json:"workdir"`
	Output     string `json:"Output"`
}

type CgroupConfig struct {
	Path        string `json:"path"`
	MemoryQuota string `json:"memory_quota"`
	CpuQuota    string `json:"cpu_quota"`
	Memory      string `json:"memory"`
	Cpu         string `json:"cpu"`
}

type RunCheckerConfig struct {
	Name    string `json:"name"`
	LogName string `json:"log_name"`
}

type ProjectConfig struct {
	Name string `json:"name"`
	Self string `json:"self"`
}

type GlobalConfig struct {
	ProjectConfig    `json:"project"`
	StorageConfig    `json:"storage"`
	RunCheckerConfig `json:"runchecker"`
	CgroupConfig     `json:"cgroup"`
}

// 项目所有配置
var (
	Global     GlobalConfig //全局配置
	Storage    StorageConfig
	Cgroup     CgroupConfig
	RunChecker RunCheckerConfig
	Project    ProjectConfig
)

// 初始化全局参数，从config.json中读入参数
func InitConfig() {
	fmt.Println("init utils, loading params...")
	wd, _ := os.Getwd()
	config, err := ParseJsonConfig(path.Join(wd, "config.json"))
	if err != nil {
		fmt.Println("loading configuration error:", err)
		os.Exit(-1)
	}

	Global = config
	Storage = config.StorageConfig
	Cgroup = config.CgroupConfig
	RunChecker = config.RunCheckerConfig
	Project = config.ProjectConfig
}

func ParseJsonConfig(jsonPath string) (GlobalConfig, error) {
	jsonConfig, err := os.Open(jsonPath)
	if err != nil {
		return GlobalConfig{}, err
	}
	defer jsonConfig.Close()
	var config GlobalConfig

	err = json.NewDecoder(jsonConfig).Decode(&config)
	if err != nil {
		return GlobalConfig{}, err
	}
	return config, nil
}

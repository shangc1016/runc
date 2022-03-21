package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Storage struct {
	Path       string `json:"path"`
	Status     string `json:"status"`
	Images     string `json:"images"`
	Logs       string `json:"logs"`
	Containers string `json:"containers"`
}

type Cgroup struct {
	Path        string `json:"path"`
	MemoryQuota string `json:"memory_quota"`
	CpuQuota    string `json:"cpu_quota"`
	Memory      string `json:"memory"`
	Cpu         string `json:"cpu"`
}

type RunChecker struct {
	Name    string `json:"name"`
	LogName string `json:"log_name"`
}

type Configuration struct {
	Name       string `json:"name"`
	Self       string `json:"self"`
	Storage    `json:"storage"`
	RunChecker `json:"runchecker"`
	Cgroup     `json:"cgroup"`
}

var (
	Name        string
	Self        string
	StoragePath string
	Status      string
	STATUS_PATH string
	Containers  string
	CTNERS_PATH string
	Images      string
	IMGS_PATH   string
	Logs        string
	LOGS_PATH   string
	CgroupPath  string
	MemQuota    string
	CpuQuota    string
	Mem         string
	Cpu         string
	RCName      string
	RCLogName   string
)

// 初始化全局参数，从config.json中读入参数
func InitConfig(configPath string) {
	config, err := ParseJsonConfig(configPath)
	if err != nil {
		fmt.Println("loading configuration error:", err)
		os.Exit(-1)
	}
	Name = config.Name
	Self = config.Self
	StoragePath = config.Storage.Path
	Status = config.Storage.Status
	Containers = config.Storage.Containers
	Images = config.Storage.Images
	Logs = config.Storage.Logs
	CgroupPath = config.Cgroup.Path
	MemQuota = config.Cgroup.MemoryQuota
	CpuQuota = config.Cgroup.CpuQuota
	Mem = config.Cgroup.Memory
	Cpu = config.Cgroup.Cpu
	RCName = config.RunChecker.Name
	RCLogName = config.RunChecker.LogName

	STATUS_PATH = path.Join(StoragePath, Status)
	CTNERS_PATH = path.Join(StoragePath, Containers)
	IMGS_PATH = path.Join(StoragePath, Images)
	LOGS_PATH = path.Join(StoragePath, Logs)

	// test
	// fmt.Println("Name:", Name)
	// fmt.Println("Self:", Self)
	// fmt.Println("storagePath:", StoragePath)
	// fmt.Println("status:", Status)
	// fmt.Println("images:", Images)
	// fmt.Println("containers:", Containers)
	// fmt.Println("logs:", Logs)
	// fmt.Println("RCName:", RCName)
	// fmt.Println("RCLogName:", RCLogName)
	// fmt.Println("CgroupPath:", CgroupPath)
	// fmt.Println("MemQuota:", MemQuota)
	// fmt.Println("CpuQuota:", CpuQuota)
	// fmt.Println("Mem:", Mem)
	// fmt.Println("Cpu:", Cpu)
}

func ParseJsonConfig(jsonPath string) (Configuration, error) {
	jsonConfig, err := os.Open(jsonPath)
	if err != nil {
		return Configuration{}, err
	}
	defer jsonConfig.Close()
	var config Configuration

	err = json.NewDecoder(jsonConfig).Decode(&config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}

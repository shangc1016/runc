package utils

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
	Storage StorageConfig = StorageConfig{
		Path:       "/var/lib/runc",
		Status:     "status",
		Images:     "images",
		Containers: "containers",
		Logs:       "logs",
		Mnt:        "mnt",
		UpperDir:   "upperdir",
		WorkDir:    "workdir",
		Output:     "output",
	}

	Cgroup CgroupConfig = CgroupConfig{
		Path:        "/sys/fs/cgroup",
		MemoryQuota: "memory.limit_in_bytes",
		CpuQuota:    "cpu.cfs_quota_us",
		Memory:      "memory",
		Cpu:         "cpu",
	}

	RunChecker RunCheckerConfig = RunCheckerConfig{
		Name:    "runchecker",
		LogName: "runchecker.log",
	}

	Project ProjectConfig = ProjectConfig{
		Name: "runc",
		Self: "/proc/self/exe",
	}
)

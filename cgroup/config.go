package cgroup

var (
	CgroupName      string = "runc"
	MemoryQuotaFile string = "memory.limit_in_bytes"
	CpuQuotaFile    string = "cpu.cfs_quota_us"
)

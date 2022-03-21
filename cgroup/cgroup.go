package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"gitee.com/shangc1016/run-c/utils"
	"github.com/sirupsen/logrus"
)

type ResourceItem struct {
	Type  string
	File  string
	Quota string
}

type CgroupResource struct {
	CgroupPath string
	CgroupName string
	Resource   []ResourceItem
}

func NewMemoryResource(quota string) *ResourceItem {
	return &ResourceItem{
		Type:  "memory",
		File:  MemoryQuotaFile,
		Quota: quota,
	}
}

func NewCpuResource(quota string) *ResourceItem {
	return &ResourceItem{
		Type:  "cpu",
		File:  CpuQuotaFile,
		Quota: quota,
	}
}

func NewCgroupResource() *CgroupResource {
	return &CgroupResource{
		CgroupName: CgroupName,
	}
}

func (c *CgroupResource) SetCgroupPath(path string) {
	c.CgroupPath = path
}

func (c *CgroupResource) AddCgroupResource(resource ResourceItem) {
	c.Resource = append(c.Resource, resource)
}

func (c *CgroupResource) InitNode() error {
	for _, resource := range c.Resource {
		path := path.Join(c.CgroupPath, resource.Type, c.CgroupName)
		// fmt.Println(1, path)
		exist, err := utils.PathExists(path)
		if err != nil {
			logrus.Fatal(err)
			return err
		}
		if !exist {
			return os.MkdirAll(path, 0644)
		}
	}
	return nil
}

func (c *CgroupResource) RemoveNode() error {
	for _, resource := range c.Resource {
		path := path.Join(c.CgroupPath, resource.Type, c.CgroupName)
		// fmt.Println(2, path)
		exist, _ := utils.PathExists(path)
		if exist {
			os.Remove(path)
		}
	}
	return nil
}

func (c *CgroupResource) SetQuota() error {
	for _, resource := range c.Resource {
		path := path.Join(c.CgroupPath, resource.Type, c.CgroupName, resource.File)
		fmt.Println(3, path, resource.Quota)
		if err := ioutil.WriteFile(path, []byte(resource.Quota), 0644); err != nil {
			fmt.Println(resource.Type, path, resource.Quota)
			logrus.Fatal(err)
			return err
		}
	}
	return nil
}

func (c *CgroupResource) SetPid(pid int) error {
	for _, resource := range c.Resource {
		path := path.Join(c.CgroupPath, resource.Type, c.CgroupName, "tasks")
		// fmt.Println(4, path)
		if err := ioutil.WriteFile(path, []byte(strconv.Itoa(pid)), 0644); err != nil {
			logrus.Fatal(err)
			return err
		}
	}
	return nil
}

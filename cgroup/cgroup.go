package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"gitee.com/shangc1016/runc/utils"
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
		Type:  utils.Mem,
		File:  utils.MemQuota,
		Quota: quota,
	}
}

func NewCpuResource(quota string) *ResourceItem {
	return &ResourceItem{
		Type:  utils.Cpu,
		File:  utils.CpuQuota,
		Quota: quota,
	}
}

func NewCgroupResource() *CgroupResource {
	return &CgroupResource{
		CgroupName: utils.Name,
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
		if err := ioutil.WriteFile(path, []byte(strconv.Itoa(pid)), 0644); err != nil {
			logrus.Fatal(err)
			return err
		}
	}
	return nil
}

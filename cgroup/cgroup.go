package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gitee.com/shangc1016/runc/fsys"
	"gitee.com/shangc1016/runc/utils"
	"github.com/sirupsen/logrus"
)

type ResourceItem struct {
	Type  string
	File  string
	Quota string
}

type CgroupResource struct {
	Path     string
	Name     string
	Root     string
	Pid      string
	Resource []ResourceItem
}

func NewCgroupResource(name, pid string) *CgroupResource {
	return &CgroupResource{
		Name:     name,               // 容器ID
		Path:     utils.Cgroup.Path,  //  `/sys/fs/cgroup`
		Root:     utils.Project.Name, // 项目名  runc
		Pid:      pid,                // 实施cgroup的进程
		Resource: []ResourceItem{},
	}
}

func (c *CgroupResource) AddCgroupResource(resource ResourceItem) {
	c.Resource = append(c.Resource, resource)
}

func (c *CgroupResource) Execute() {
	err := c.setQuota()
	if err != nil {
		fmt.Println("set quota error", err)
		os.Exit(-1)
	}
	err = c.setPid()
	if err != nil {
		fmt.Println("set pid error", err)
		os.Exit(-1)
	}
}

func (c *CgroupResource) setQuota() error {
	for _, resource := range c.Resource {
		if err := os.Mkdir(path.Join(c.Path, resource.Type, c.Root, c.Name), 0644); err != nil {
			return err
		}
		quota_path := path.Join(c.Path, resource.Type, c.Root, c.Name, resource.File)
		// FIXME:设置cpu限制的时候，写入文件报错？为啥
		if err := ioutil.WriteFile(quota_path, []byte(resource.Quota), 0644); err != nil {
			fmt.Println("set Quota error", err)
			return err
		}
	}
	return nil
}

func (c *CgroupResource) setPid() error {
	for _, resource := range c.Resource {
		tasks_path := path.Join(c.Path, resource.Type, c.Root, c.Name, "tasks")
		fmt.Println("tasks_path", tasks_path, c.Pid)
		if err := ioutil.WriteFile(tasks_path, []byte(c.Pid), 0644); err != nil {
			logrus.Fatal(err)
			return err
		}
	}
	return nil
}

func (c *CgroupResource) RemoveNode() bool {
	for _, resource := range c.Resource {
		node_path := path.Join(c.Path, resource.Type, c.Root, c.Name)
		fmt.Println("node_path:", node_path)
		exist, err := fsys.PathExists(node_path)
		if err != nil {
			fmt.Println("<removeNode error>:", err)
		}
		if exist {
			if err := os.RemoveAll(node_path); err != nil {
				fmt.Println("<remove cgroup node err>", err)
			}
		}
	}
	fmt.Println("node removed")
	return true
}

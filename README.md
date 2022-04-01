# runc

> 这是根据[自己动手写Docker](https://book.douban.com/subject/27082348/)这本书的一个toy.



### 1、进入容器环境的总体流程

#### 		1.1、run子命令设置clone的flag

#### 		1.2、init子命令设置挂载点

### 2、使用管道处理容器内部命令


### 2、overlay包裹busybox作为容器只读层

>- https://docs.docker.com/storage/storagedriver/overlayfs-driver/#how-the-overlay2-driver-works
>- https://blogs.cisco.com/developer/373-containerimages-03
>- https://www.cnblogs.com/FengZeng666/p/14173906.html
>- https://docs.kernel.org/filesystems/overlayfs.html
>- https://wiki.archlinux.org/title/Overlay_filesystem

<img src="https://note-img-1300721153.cos.ap-nanjing.myqcloud.com//md-imgimage-20220318232449145.png" style="zoom:50%;" />

> docker 用到的overlay文件系统示意



...

### 3、增加volume数据卷

...



### 5、实现镜像简单打包

> 文中的镜像打包，是直接去到挂载的目录，然后用`tar`命令，把mnt目录打包， 感觉这样有点过于简单

### 6、进入容器
需要setns这个系统调用，可以根据提供的pid进程号进入到进程所在的NS中，原理是读取/proc/<pid>/ns下面的ns号码，然后进入。
这个系统调用不能再多线程环境调用，因此需要使用到go的cgo机制，借助c语言实现，然后go代码调用c代码。









## chapter-5   

> - 5.1 容器后台运行
> - 5.2查看运行状态的容器
> - 5.3查看容器日志
> - 5.4进入容器namespace
> - 5.5停止容器
> - 5.6删除容器
> - 5.7通过容器制作镜像
> - 5.8容器指定环境变量

#### 5.1、容器后台运行  

书中5.1小节：在detach模式运行的容器进程，好像没有在容器进程被kill掉之后 移除文件系统挂载与数据卷挂载的相关逻辑？？？

#### 5.2、查看运行的容器

即`docker ps`，原理是把后台运行的容器信息保存在文件中, ps命令就是去读文件、然后打印

#### 5.3、查看容器日志

把后台运行容器的标准输出重定向到本地文件, 然后读取文件输出。

#### 5.4、进入容器Namespace, `exec`

利用setns系统调用，以及go、c之间的调用

FIXME：在这儿遇到了问题，使用setns进不去容器的命名空间，

#### 5.5、停止容器, `kill`

直接杀死进程，更新容器状态为`terminalted`，但是不销毁容器的文件系统

#### 5.6、删除容器，`rm`

删除容器的文件系统，并且删除容器的状态文件









<img src="https://note-img-1300721153.cos.ap-nanjing.myqcloud.com//md-imgimage-20220320223401065.png" alt="image-20220320223401065" style="zoom:67%;" />




#### 项目目录结构：

```
....
```






#### 存储目录结构：/var/lib/runc

```
/var/lib/runc/
							| -- containers                  // 容器文件系统相关
							|     | -- xxxx                  // 以容器ID命名
							|     |     | -- mnt             // 容器挂载的目录，即overaly文件系统中的merged
							|     |     | -- upperdir        // overlay文件系统的upperdir 
							|     |     | -- workdir         // overlay文件系统的workdir
							|     |     | -- output          // 把容器输出挂载到此处
							|     |   .. 
							|
							|
							| -- images                      // 基础镜像，即overlay文件系统中的lowerdir
							|     | -- busybox               // busybox基础镜像
							|     |   ..
							|
							| -- status                      // 容器状态
							|     | -- xxxx                  // 以容器ID命名，保存容器状态
							|     |   ..
							|
							| -- logs                        // 运行日日志
							      | -- runchecker.log        // runchecker的日志
							      |   ..
							
```



#### cgroup目录

```
/sys/fs/cgroup/xxx/runc    // xxx为不同的资源类型，cpu、memery、cpuset等
```








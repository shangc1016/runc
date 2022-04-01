package nsenter

/*

#define _GNU_SOURCE
#include <errno.h>
#include <fcntl.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>


// constructor 在程序启动的时候运行
__attribute__((constructor)) void enter_namespace(void) {
  char *pid;
  printf("shangchao\n");
  pid = getenv("runc_pid");
  // 没有run_pid环境变量，直接退出
  if (!pid) return;
  char *cmd;
  cmd = getenv("runc_cmd");
  // 没有runc_cmd环境变量，直接退出
  if (!cmd) return;

  printf("%s\n", pid);
  printf("%s\n", cmd);

  int i;
  char nspath[1024];
  char *namespace[] = {"ipc", "net", "pid", "uts", "mnt"};
  for (i = 0; i < 5; i++) {
    // 拼接路径
    sprintf(nspath, "/proc/%s/ns/%s", pid, namespace[i]);
	printf("nspath:%s\n", nspath);
    int fd = open(nspath, O_RDONLY, 0644);
    if (setns(fd, 0) == -1) {
      fprintf(stderr, "setns on %s namespace failed, %s\n", namespace[i],
              strerror(errno));
    } else {
      fprintf(stdout, "setns on %s namespace succeeded\n", namespace[i]);
    }
    close(fd);
  }
  int res = system(cmd);
  exit(0);
  return;
}

*/
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"gitee.com/shangc1016/runc/status"
	"gitee.com/shangc1016/runc/utils"
)

const ENV_EXEC_PID = "runc_pid"
const ENV_EXEC_CMD = "runc_cmd"

// 根据容器的Id得到容器的pid进程号
func getPidById(id string) string {
	files, err := utils.GetFileNameAll(path.Join(utils.Storage.Path, utils.Storage.Status))
	if err != nil {
		fmt.Println("<get container error>", err)
		fmt.Println(path.Join(utils.Storage.Path, utils.Storage.Status))
		os.Exit(-1)
	}
	for _, file := range files {
		data, _ := utils.ReadFile(path.Join(utils.Storage.Path, utils.Storage.Status, file))
		var info status.ContainerInfo
		json.Unmarshal(data, &info)
		if info.Id == id {
			return info.Pid
		}
	}
	return ""
}

// 设置好两个环境变量，再次运行程序，携带同样的自命令exec
func ExecContainer(id string, commandArr []string) {
	// 根据容器的Id得到容器的进程号, eg id:y8ien6qlei  --> pis:44757
	pid := getPidById(id)
	fmt.Println("pid:", pid)
	fmt.Println("commandArr:", commandArr)
	if pid == "" {
		fmt.Println("exec container get id error")
		os.Exit(-1)
	}

	cmdStr := strings.Join(commandArr, " ")
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	os.Setenv(ENV_EXEC_CMD, cmdStr)
	os.Setenv(ENV_EXEC_PID, pid)
	if err := cmd.Run(); err != nil {
		fmt.Println("exec run error")
		os.Exit(-1)
	}
}

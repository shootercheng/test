package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
//	"runtime"
	"syscall"
)

func main() {
	// 目标容器的新根目录，需确保该目录已存在
	newRoot := "/var/sandbox/sandbox-python"

	// 1. 将当前 goroutine 锁定到一个系统线程，因为命名空间是线程局部属性
//	runtime.LockOSThread()
//	defer runtime.UnlockOSThread()

	// 2. 创建新的用户命名空间 (User Namespace)
	//    普通用户创建后，即可在命名空间内获得完整权限 (CAP_SYS_ADMIN等)
	if err := syscall.Unshare(syscall.CLONE_NEWUSER); err != nil {
		log.Fatalf("Unshare CLONE_NEWUSER 失败: %v", err)
	}

	// 3. 创建新的挂载命名空间 (Mount Namespace)
	if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
		log.Fatalf("Unshare CLONE_NEWNS 失败: %v", err)
	}

	// 4. 将当前命名空间内的根目录 ("/") 重新挂载为私有模式 (MS_PRIVATE | MS_REC)
	//    这能防止在此命名空间内的挂载事件传播到宿主机，实现真正的隔离
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		log.Fatalf("重新挂载根目录为私有失败: %v", err)
	}

	// 5. 切换根目录到指定的容器目录
	if err := syscall.Chroot(newRoot); err != nil {
		log.Fatalf("Chroot 失败: %v", err)
	}

	// 6. 切换工作目录到新的根目录
	if err := os.Chdir("/"); err != nil {
		log.Fatalf("切换工作目录失败: %v", err)
	}

	fmt.Println("成功进入隔离环境，正在启动 shell...")

	cmd := exec.Command("cat", "/etc/passwd")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("命令执行失败: %v\n错误信息: %s", err, stderr.String())
	}

	fmt.Println("标准输出：")
	fmt.Println(stdout.String())

}

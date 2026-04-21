package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	newRoot := "/var/sandbox/sandbox-python"

	// nums := strings.Split("1,2,3", ",")
	// for num := range nums {
	// 	syscall, err := strconv.Atoi(nums[num])
	// 	if err != nil {
	// 		fmt.Printf("Atoi err:%v\n", err)
	// 		continue
	// 	}
	// 	fmt.Println(syscall)
	// }

	if err := syscall.Chroot(newRoot); err != nil {
		log.Fatalf("chroot失败: %v", err)
	}

	if err := os.Chdir("/"); err != nil {
		log.Fatalf("切换工作目录失败: %v", err)
	}

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

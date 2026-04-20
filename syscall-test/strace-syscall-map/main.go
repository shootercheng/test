package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Println("用法: go run main.go <python脚本>")
	// 	fmt.Println("示例: go run main.go json_print.py")
	// 	os.Exit(1)
	// }
	sysCallMap, err := getSyscallSMap()
	if err != nil {
		fmt.Printf("get syscall error:%s\n", err.Error())
		return
	}

	// 执行 strace 命令
	cmd := exec.Command("strace", "-c", "python", "test.py", "/var/sandbox/sandbox-python", "3<json_print.py")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行 strace 失败: %v\n", err)
		// 继续处理输出
	}

	fmt.Printf("syscall map info %v\n", sysCallMap)

	// 解析输出
	parseStraceOutput(string(output), sysCallMap)
}

func getSyscallSMap() (map[string]string, error) {
	filePath := "/usr/include/x86_64-linux-gnu/asm/unistd_64.h"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return nil, err
	}
	defer file.Close()

	// 正则表达式匹配 #define __NR_xxx 数字
	re := regexp.MustCompile(`^#define\s+__NR_(\w+)\s+(\d+)$`)

	syscalls := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]
			number := matches[2]
			syscalls[name] = number
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		return nil, err
	}
	return syscalls, nil
}

func parseStraceOutput(output string, sysCallMap map[string]string) {
	lines := strings.Split(output, "\n")

	// 正则表达式匹配 strace -c 输出行
	// 格式: 19.46    0.000652           4       149        30 newfstatat
	re := regexp.MustCompile(`^\s*[\d.]+\s+[\d.]+\s+(\d+)\s+(\d+)\s+(\d+|\s+)\s+(\w+)`)

	index := 1

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); matches != nil {
			calls := matches[2]
			errors := matches[3]
			syscall := matches[4]

			sysCallNum, ok := sysCallMap[syscall]
			if ok {
				result = append(result, sysCallNum)
			} else {
				fmt.Printf("not match syscall %s\n", syscall)
			}

			fmt.Printf("%-6d %-20s %-10s %-10s\n",
				index, syscall, calls, errors)
			index++
		} else {
			fmt.Printf("not match line %s\n", line)
		}
	}
	slices.SortFunc(result, func(a, b string) int {
		numA, _ := strconv.Atoi(a)
		numB, _ := strconv.Atoi(b)
		if numA < numB {
			return -1
		} else if numA > numB {
			return 1
		}
		return 0
	})
	sysCallResult := strings.Join(result, ",")
	fmt.Println(sysCallResult)
}

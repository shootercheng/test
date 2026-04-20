package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	filePath := "/usr/include/x86_64-linux-gnu/asm/unistd_64.h"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer file.Close()

	// 正则表达式匹配 #define __NR_xxx 数字
	re := regexp.MustCompile(`^#define\s+__NR_(\w+)\s+(\d+)$`)

	syscalls := make(map[int]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]
			number, _ := strconv.Atoi(matches[2])
			syscalls[number] = name
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		return
	}

	// 排序输出
	numbers := make([]int, 0, len(syscalls))
	for num := range syscalls {
		numbers = append(numbers, num)
	}
	sort.Ints(numbers)

	// 输出到文件
	outputFile, err := os.Create("syscalls.txt")
	if err != nil {
		fmt.Printf("无法创建输出文件: %v\n", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	fmt.Fprintf(writer, "%-6s | %s\n", "Number", "Syscall Name")
	fmt.Fprintf(writer, "-------+-------------------\n")

	for _, num := range numbers {
		fmt.Fprintf(writer, "%-6d | %s\n", num, syscalls[num])
	}
	writer.Flush()

	fmt.Printf("成功提取 %d 个系统调用到 syscalls.txt\n", len(syscalls))
}

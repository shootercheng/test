package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	sg "github.com/seccomp/libseccomp-golang"
)

const (
	SYS_SECCOMP            = 317
	SeccompSetModeFilter   = 0x1
	SeccompFilterFlagTSYNC = 0x1
)

func main() {
	allowedSysCall := []int{
		0, 1, 3, 5, 8, 9, 10, 11, 12, 13, 16, 17, 21, 59, 72, 79, 89, 158, 186, 202, 217, 218, 257, 262, 273, 302, 318, 334,
	}
	err := seccomp(allowedSysCall[:3], []int{})
	if err != nil {
		fmt.Printf("seccomp error:%v\n", err)
		return
	}

	cmd := exec.Command("python", "./python/json_print.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("exec python error: %v\n output: %s\n", err, output)
		return
	}
	fmt.Println(string(output))
}

func seccomp(allowed_syscalls []int, allowed_not_kill_syscalls []int) error {
	ctx, err := sg.NewFilter(sg.ActKillProcess)
	if err != nil {
		return err
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		return err
	}
	defer reader.Close()
	defer writer.Close()

	for _, syscall := range allowed_syscalls {
		ctx.AddRule(sg.ScmpSyscall(syscall), sg.ActAllow)
	}

	for _, syscall := range allowed_not_kill_syscalls {
		ctx.AddRule(sg.ScmpSyscall(syscall), sg.ActErrno)
	}

	file := os.NewFile(uintptr(writer.Fd()), "pipe")
	ctx.ExportBPF(file)

	// read from pipe
	data := make([]byte, 4096)
	n, err := reader.Read(data)
	if err != nil {
		return err
	}
	// load bpf
	sock_filters := make([]syscall.SockFilter, n/8)
	bytesBuffer := bytes.NewBuffer(data)
	err = binary.Read(bytesBuffer, binary.LittleEndian, &sock_filters)
	if err != nil {
		return err
	}

	bpf := syscall.SockFprog{
		Len:    uint16(len(sock_filters)),
		Filter: &sock_filters[0],
	}

	_, _, err2 := syscall.Syscall(
		SYS_SECCOMP,
		uintptr(SeccompSetModeFilter),
		uintptr(SeccompFilterFlagTSYNC),
		uintptr(unsafe.Pointer(&bpf)),
	)

	if err2 != 0 {
		return err2
	}

	return nil
}

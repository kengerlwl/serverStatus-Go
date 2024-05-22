package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// 执行命令并获取输出的函数
func execCommand(command string) (int, error) {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return 0, err
	}
	output := strings.TrimSpace(string(out))
	result, err := strconv.Atoi(output)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// tupd 函数：获取TCP连接数、UDP连接数、进程数和线程数
func tupd() (int, int, int, int, error) {
	// 获取TCP连接数
	tcpCount, err := execCommand("ss -t | wc -l")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	tcpCount--

	// 获取UDP连接数
	udpCount, err := execCommand("ss -u | wc -l")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	udpCount--

	// 获取进程数
	processCount, err := execCommand("ps -ef | wc -l")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	processCount -= 2

	// 获取线程数
	threadCount, err := execCommand("ps -eLf | wc -l")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	threadCount -= 2

	return tcpCount, udpCount, processCount, threadCount, nil
}

func main() {
	tcp, udp, process, thread, err := tupd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("TCP connections: %d\n", tcp)
	fmt.Printf("UDP connections: %d\n", udp)
	fmt.Printf("Processes: %d\n", process)
	fmt.Printf("Threads: %d\n", thread)
}

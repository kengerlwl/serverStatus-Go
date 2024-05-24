package status

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	pNet "github.com/shirou/gopsutil/v3/net"
)

var cachedFs = make(map[string]struct{})
var timer = 0.0
var prevNetIn uint64
var prevNetOut uint64

// GPUUsage 包含 GPU 使用情况的信息
type GPUUsage struct {
	GPUIndex          int // GPU索引
	GPUUtilization    int // GPU利用率（百分比）
	MemoryUtilization int // 显存利用率（百分比）
	MemoryTotal       int // 总显存（MiB）
	MemoryUsed        int // 已用显存（MiB）
}

// execCommand 执行给定的shell命令并返回其输出
func execCommand(command string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// GetGpuUtils 获取系统的 GPU 使用情况
func GetGpuUtils() ([]GPUUsage, error) {
	// 执行 nvidia-smi 命令
	output, err := execCommand("nvidia-smi --query-gpu=index,memory.used,memory.total,utilization.gpu --format=csv,noheader,nounits")
	if err != nil {
		return nil, fmt.Errorf("failed to execute nvidia-smi: %w", err)
	}

	// 解析 nvidia-smi 的输出
	var gpuUsages []GPUUsage
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ", ")
		if len(fields) < 4 {
			return nil, fmt.Errorf("unexpected output format: %s", line)
		}

		gpuIndex, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse GPU index: %w", err)
		}

		memoryUsed, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse memory used: %w", err)
		}

		memoryTotal, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse total memory: %w", err)
		}

		gpuUtilization, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse GPU utilization: %w", err)
		}

		gpuUsages = append(gpuUsages, GPUUsage{
			GPUIndex:          gpuIndex,
			GPUUtilization:    gpuUtilization,
			MemoryUtilization: memoryUsed * 100 / memoryTotal,
			MemoryTotal:       memoryTotal,
			MemoryUsed:        memoryUsed,
		})
	}

	return gpuUsages, nil
}

// // getGpuUtils 获取系统的GPU使用情况
// func GetGpuUtils() ([]GPUUsage, error) {
//     // 执行nvidia-smi命令
//     output, err := execCommand("nvidia-smi --query-gpu=memory.used,memory.total,utilization.gpu --format=csv,noheader,nounits")
//     if err != nil {
//         return nil, fmt.Errorf("failed to execute nvidia-smi: %w", err)
//     }

//     // 解析nvidia-smi的输出
//     gpuUsages, err := parseNvidiaSmiOutput(output)
//     if err != nil {
//         return nil, fmt.Errorf("failed to parse nvidia-smi output: %w", err)
//     }

//     return gpuUsages, nil
// }

func Uptime() uint64 {
	bootTime, _ := host.BootTime()
	return uint64(time.Now().Unix()) - bootTime
}

func Load() float64 {
	theLoad, _ := load.Avg()
	return theLoad.Load1
}

// execCommand 执行一个给定的 shell 命令，并返回结果的整数值。
func execNetCommand(command string) (int, error) {
	var out bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	// 去除输出中的换行符并转换为整数
	output := strings.TrimSpace(out.String())
	result, err := strconv.Atoi(output)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func Tupd() (int, int, int, int) {
	tcpCount, _ := execNetCommand("ss -t | wc -l")
	udpCount, _ := execNetCommand("ss -u | wc -l")
	processCount, _ := execNetCommand("ps -ef | wc -l")
	threadCount, _ := execNetCommand("ps -eLf | wc -l")
	return tcpCount - 1, udpCount - 1, processCount - 2, threadCount - 2
	// var a, b, c, d int
	// a = 0
	// b = 0
	// c = 0
	// d = 0
	// return a, b, c, d
	// return 0, 0, 0, 0, nil
}

func GetLoadAvg() (float64, float64, float64, error) {
	avg, err := load.Avg()
	if err != nil {
		return 0, 0, 0, err
	}
	return avg.Load1, avg.Load5, avg.Load15, nil
}

func DiskIO() (uint64, uint64) {
	ioStat, _ := disk.IOCounters()
	var read, write uint64
	for _, v := range ioStat {
		read += v.ReadBytes
		write += v.WriteBytes
	}
	return read, write
}

func Disk(INTERVAL float64) (uint64, uint64) {
	var (
		size, used uint64
	)
	if timer <= 0 {
		diskList, _ := disk.Partitions(false)
		devices := make(map[string]struct{})
		for _, d := range diskList {
			_, ok := devices[d.Device]
			if !ok && checkValidFs(d.Fstype) {
				cachedFs[d.Mountpoint] = struct{}{}
				devices[d.Device] = struct{}{}
			}
		}
		timer = 300.0
	}
	timer -= INTERVAL
	for k := range cachedFs {
		usage, err := disk.Usage(k)
		if err != nil {
			delete(cachedFs, k)
			continue
		}
		size += usage.Total / 1024.0 / 1024.0
		used += usage.Used / 1024.0 / 1024.0
	}
	return size, used
}

func Cpu(INTERVAL float64) float64 {
	cpuInfo, _ := cpu.Percent(time.Duration(INTERVAL*float64(time.Second)), false)
	return math.Round(cpuInfo[0]*10) / 10
}

func Network(checkIP int) bool {
	var HOST string
	if checkIP == 4 {
		HOST = "8.8.8.8:53"
	} else if checkIP == 6 {
		HOST = "[2001:4860:4860::8888]:53"
	} else {
		return false
	}
	conn, err := net.DialTimeout("tcp", HOST, 2*time.Second)
	if err != nil {
		return false
	}
	if conn.Close() != nil {
		return false
	}
	return true
}

func Traffic(INTERVAL float64) (uint64, uint64, uint64, uint64) {
	var (
		netIn, netOut uint64
	)
	netInfo, _ := pNet.IOCounters(true)
	for _, v := range netInfo {
		if checkInterface(v.Name) {
			netIn += v.BytesRecv
			netOut += v.BytesSent
		}
	}
	rx := uint64(float64(netIn-prevNetIn) / INTERVAL)
	tx := uint64(float64(netOut-prevNetOut) / INTERVAL)
	prevNetIn = netIn
	prevNetOut = netOut
	return netIn, netOut, rx, tx
}

func TrafficVnstat() (uint64, uint64, error) {
	buf, err := exec.Command("vnstat", "--oneline", "b").Output()
	if err != nil {
		return 0, 0, err
	}
	vData := strings.Split(BytesToString(buf), ";")
	if len(vData) != 15 {
		// Not enough data available yet.
		return 0, 0, nil
	}
	netIn, err := strconv.ParseUint(vData[8], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	netOut, err := strconv.ParseUint(vData[9], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return netIn, netOut, nil
}

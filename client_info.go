package main

import (
	"flag"
	"fmt"
	"log"

	// "strings"
	"time"

	status "kenger.work/kenger/serverStatusGo"

	jsoniter "github.com/json-iterator/go"
)

// GetLoadAvg 获取系统的 1 分钟、5 分钟和 15 分钟负载平均值

var (
	SERVER   = flag.String("h", "", "Input the host of the server")
	PORT     = flag.Int("port", 35601, "Input the port of the server")
	USER     = flag.String("u", "", "Input the client's username")
	PASSWORD = flag.String("p", "", "Input the client's password")
	INTERVAL = flag.Float64("interval", 2.0, "Input the INTERVAL")
	DSN      = flag.String("dsn", "", "Input DSN, format: username:password@host:port")
	isVnstat = flag.Bool("vnstat", false, "Use vnstat for traffic statistics, linux only")
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ServerStatus struct {
	Name     string `json:"name"` // 服务器名称
	Type     string `json:"type"` // 服务器类型
	Host     string `json:"host"` // 服务器地址
	Location string `json:"location"`

	Uptime      uint64          `json:"uptime"`
	Load1       float64         `json:"load_1"`
	Load5       float64         `json:"load_5"`
	Load15      float64         `json:"load_15"`
	Ping10010   float64         `json:"ping_10010"`
	Ping189     float64         `json:"ping_189"`
	Ping10086   float64         `json:"ping_10086"`
	Time10010   float64         `json:"time_10010"`
	Time189     float64         `json:"time_189"`
	Time10086   float64         `json:"time_10086"`
	TCP         int             `json:"tcp"`
	UDP         int             `json:"udp"`
	Process     int             `json:"process"`
	Thread      int             `json:"thread"`
	IORead      uint64          `json:"io_read"`
	IOWrite     uint64          `json:"io_write"`
	Custom      string          `json:"custom"`
	MemoryTotal uint64          `json:"memory_total"`
	MemoryUsed  uint64          `json:"memory_used"`
	SwapTotal   uint64          `json:"swap_total"`
	SwapUsed    uint64          `json:"swap_used"`
	HddTotal    uint64          `json:"hdd_total"`
	HddUsed     uint64          `json:"hdd_used"`
	CPU         jsoniter.Number `json:"cpu"`
	NetworkTx   uint64          `json:"network_tx"`
	NetworkRx   uint64          `json:"network_rx"`
	NetworkIn   uint64          `json:"network_in"`
	NetworkOut  uint64          `json:"network_out"`
	Online4     bool            `json:"online4"`
	Online6     bool            `json:"online6"`
}

// 获取当前服务器的状态信息
func getServerStatus() ServerStatus {
	timer := 0.0
	checkIP := 4

	item := ServerStatus{}

	CPU := status.Cpu(*INTERVAL)
	var netIn, netOut, netRx, netTx uint64
	if !*isVnstat {
		netIn, netOut, netRx, netTx = status.Traffic(*INTERVAL)
	} else {
		_, _, netRx, netTx = status.Traffic(*INTERVAL)
		netIn, netOut, _ = status.TrafficVnstat()
		// if err != nil {
		// 	log.Println("Please check if the installation of vnStat is correct")
		// }
		log.Printf("Traffic: %d %d\n", netIn, netOut)
	}
	memoryTotal, memoryUsed, swapTotal, swapUsed := status.Memory()
	hddTotal, hddUsed := status.Disk(*INTERVAL)
	uptime := status.Uptime()
	// load := status.Load()
	item.CPU = jsoniter.Number(fmt.Sprintf("%.1f", CPU))
	// item.Load = jsoniter.Number(fmt.Sprintf("%.2f", load))
	item.Load1, item.Load5, item.Load15, _ = status.GetLoadAvg()

	// to fix
	item.Ping10010, item.Time10010 = 0, 0
	item.Ping189, item.Time189 = 0, 0
	item.Ping10086, item.Time10086 = 0, 0
	item.TCP, item.UDP, item.Process, item.Thread = status.Tupd()
	item.IORead, item.IOWrite = status.DiskIO()
	item.Custom = "custom"

	// user info
	confs, err := LoadConfig("para.client.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	serverName := confs["serverName"].(string)
	serverType := confs["serverType"].(string)
	serverHost := confs["serverHost"].(string)
	serverLocation := confs["serverLocation"].(string)
	item.Name = serverName
	item.Type = serverType
	item.Host = serverHost
	item.Location = serverLocation

	item.Uptime = uptime
	item.MemoryTotal = memoryTotal
	item.MemoryUsed = memoryUsed
	item.SwapTotal = swapTotal
	item.SwapUsed = swapUsed
	item.HddTotal = hddTotal
	item.HddUsed = hddUsed
	item.NetworkRx = netRx
	item.NetworkTx = netTx
	item.NetworkIn = netIn
	item.NetworkOut = netOut
	if timer <= 0 {
		if checkIP == 4 {
			item.Online4 = status.Network(checkIP)
		} else if checkIP == 6 {
			item.Online6 = status.Network(checkIP)
		}
		timer = 150.0
	}
	timer -= *INTERVAL
	// data, _ := json.Marshal(item)

	return item

}

func connect() {

	for {

		item := getServerStatus()
		fmt.Println(item)
		time.Sleep(5 * time.Second)
	}

}

// func main() {

// 	for {
// 		connect()
// 	}
// }

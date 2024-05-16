package main

import (
	"flag"
	"fmt"
	"log"

	// "strings"
	"time"

	"kenger.work/kenger/serverStatusGo"

	jsoniter "github.com/json-iterator/go"
)

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
	Uptime      uint64          `json:"uptime"`
	Load        jsoniter.Number `json:"load"`
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

func connect() {

	timer := 0.0
	checkIP := 4

	item := ServerStatus{}
	for {
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
		load := status.Load()
		item.CPU = jsoniter.Number(fmt.Sprintf("%.1f", CPU))
		item.Load = jsoniter.Number(fmt.Sprintf("%.2f", load))
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
		data, _ := json.Marshal(item)
		fmt.Println(data)
		fmt.Println(item)
		time.Sleep(5 * time.Second)
	}
}

func main() {

	for {
		connect()
	}
}

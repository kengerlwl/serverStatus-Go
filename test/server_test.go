package main


import (
	"testing"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
	"fmt"
)

func TestRegistor(t *testing.T) {


	for i := 0; i < 10; i++ {
		// 读取配置文件
		confs, err := LoadConfigFromEnv()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		etcdHost := confs["etcd"].(map[string]interface{})["host"].(string)
		etcdPort := confs["etcd"].(map[string]interface{})["port"].(string)

		var endpoints = []string{etcdHost + ":" + etcdPort}

		nowStr := strconv.FormatInt(time.Now().Unix(), 10)
		log.Println(nowStr)

		// 加入时间参数
		ser, err := NewServiceRegister(endpoints, "/web/node/" + nowStr, "localhost:8000", 5) // 本地的8000端口，
		if err != nil {
			log.Fatalln(err)
		}

		//监听续租相应chan
		go ser.ListenLeaseRespChan()
		select {
		// case <-time.After(20 * time.Second):
		// 	ser.Close()
		}
	}


}


func TestFind(t *testing.T) {
	// 读取配置文件
	confs, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	etcdHost := confs["etcd"].(map[string]interface{})["host"].(string)
	etcdPort := confs["etcd"].(map[string]interface{})["port"].(string)

	var endpoints = []string{etcdHost + ":" + etcdPort}
	ser := NewServiceDiscovery(endpoints)
	defer ser.Close()
	ser.WatchService("/web/")
	ser.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}
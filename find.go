package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli              *clientv3.Client         //etcd client
	serverList       map[string]string        //服务列表
	allExistedServer map[string]string        //所有存在过的服务
	allExistedServe2 map[string]*ServerStatus //所有存在过的服务
	lock             sync.Mutex
}

// NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceDiscovery{
		cli:              cli,
		serverList:       make(map[string]string),
		allExistedServer: make(map[string]string),
		allExistedServe2: make(map[string]*ServerStatus),
		lock:             sync.Mutex{},
	}
}

// WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的server
	go s.watcher(prefix)
	return nil
}

// watcher 监听前缀，用于动态监听服务节点的变化
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// SetServiceList 新增或修改服务地址
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = string(val)

	if _, flagExist := s.allExistedServer[key]; flagExist {
		// log.Println("update key :", key, "val:", val)
	} else {
		log.Println("put key:", key, "val:", val)
	}

	s.allExistedServer[key] = string(val)
	var serverStatus ServerStatus
	json.Unmarshal([]byte(val), &serverStatus)
	// fmt.Println(serverStatus)
	s.allExistedServe2[key] = &serverStatus

}

// DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del key:", key)
}

// GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for k, v := range s.serverList {
		addrs = append(addrs, string(k)+string(v))
	}
	return addrs
}

// Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func (s *ServiceDiscovery) startWebServer() {
	// 设置处理函数
	http.HandleFunc("/json/stats.json", func(w http.ResponseWriter, r *http.Request) {
		// 创建一个响应对象
		// fmt.Println("get request")
		s.lock.Lock()
		items := make([]ServerStatus, 0)
		for _, v := range s.allExistedServe2 {
			items = append(items, *v)
		}
		s.lock.Unlock()
		nowTime := time.Now().Unix()
		// 创建包含服务器列表的对象
		responseData := map[string]interface{}{
			"servers": items,
			"updated": nowTime,
		}
		// 将响应对象编码为JSON
		jsonData, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置响应头
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// 写入响应体
		w.Write(jsonData)
	})

	// 前端页面
	fs := http.FileServer(http.Dir("./static"))

	// 代理到根目录
	http.Handle("/", http.StripPrefix("/", fs))

	// 启动服务器
	http.ListenAndServe(":8080", nil)
}

func main() {
	// 读取配置文件
	// confs, err := LoadConfigFromEnv()
	confs, err := LoadConfig("para.server.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	etcdHost := confs["etcd"].(map[string]interface{})["host"].(string)
	etcdPort := confs["etcd"].(map[string]interface{})["port"].(string)

	var endpoints = []string{etcdHost + ":" + etcdPort}
	ser := NewServiceDiscovery(endpoints)
	defer ser.Close()
	ser.WatchService("/server/")
	go ser.startWebServer()
	// ser.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}

package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	cli     *clientv3.Client //etcd client
	leaseID clientv3.LeaseID //租约ID
	//租约keepalieve相应chan
	leaseTTL      int64 //租约时间
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string //key
	val           string //value
}

// NewServiceRegister 新建注册服务
func NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: val,
	}
	ser.leaseTTL = lease

	//申请租约设置时间keepalive
	if err := ser.putKeyWithLease(ser.leaseTTL); err != nil {
		return nil, err
	}

	return ser, nil
}

// 设置key和租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {

	// ctx, cancel = context.WithTimeout(context.Background(), time.Second)

	// 如果key已经存在，则注册失败
	respGet, err := s.cli.Get(context.Background(), s.key)
	if err != nil {
		return err
	}
	if len(respGet.Kvs) != 0 {
		log.Printf("key %s already exists", s.key)
		return nil
	}

	//设置租约时间
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	//注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	log.Println(s.leaseID)
	s.keepAliveChan = leaseRespChan
	log.Printf("Put key:%s  val:%s  success!", s.key, s.val)
	return nil
}

func (s *ServiceRegister) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp, ok := <-s.keepAliveChan:
			if !ok {
				// 如果续租管道关闭，则尝试重新连接
				log.Println("续租管道关闭，尝试重新连接")
				s.retryKeepAlive()
				return
			}

			log.Println("续约成功", leaseKeepResp)

			// 测试每次更新value内容
			nowServerStatus := getServerStatus()
			// 序列化为json
			data, _ := json.Marshal(nowServerStatus)
			s.val = string(data)

			_, err := s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(s.leaseID))
			if err != nil {
				log.Println("更新value失败", err)
			}
		}
	}
	log.Println("关闭续租")
}

// retryKeepAlive 尝试重新连接并恢复续租
func (s *ServiceRegister) retryKeepAlive() {
	for {
		// 尝试重新获取租约ID
		leaseResp, err := s.cli.Grant(context.Background(), int64(s.leaseTTL))
		if err != nil {
			log.Println("重新获取租约ID失败", err)
			time.Sleep(2 * time.Second) // 等待2秒后重试
			continue
		}

		s.leaseID = leaseResp.ID

		// 开始续租
		ch, err := s.cli.KeepAlive(context.Background(), s.leaseID)
		if err != nil {
			log.Println("重新开始续租失败", err)
			time.Sleep(2 * time.Second) // 等待2秒后重试
			continue
		}

		s.keepAliveChan = ch
		log.Println("重新连接并恢复续租成功")
		go s.ListenLeaseRespChan() // 重新启动监听
		return
	}
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("撤销租约")
	return s.cli.Close()
}

func main() {

	// 读取配置文件
	confs, err := LoadConfig("para.client.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	etcdHost := confs["etcd"].(map[string]interface{})["host"].(string)
	etcdPort := confs["etcd"].(map[string]interface{})["port"].(string)
	serverName := confs["serverName"].(string)

	var endpoints = []string{etcdHost + ":" + etcdPort}

	nowStr := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(nowStr)

	// 获取当前服务器状态
	nowServerStatus := getServerStatus()

	//序列化为json
	data, _ := json.Marshal(nowServerStatus)

	// 加入时间参数
	ser, err := NewServiceRegister(endpoints, "/server/"+serverName, string(data), 5) // 本地的8000端口，
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

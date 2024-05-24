# serverStatus-Go
 基于go的服务器探针系统


# run
## client
`go run client.go config.go client_info.go`

## find

`go run find.go config.go client_info.go`


# to do

- 修改配置文件的读取。让client可以设置机器名字，以及其他参数(done) 
- 服务器端添加web接口，主要是一个读取接口和一个删除某个服务器节点的接口 done
- 删除接口高危，需要鉴权。但是先不管
- etcd的设置状态也应该鉴权
- 加入前端界面对接。建议直接套用 (done)
- 客户端网络不稳定，重新上线功能
- 



## 数据补全

```go
:{"uptime":128046,"load":0.00,"memory_total":201234888,"memory_used":19282060,"swap_total":201234888,"swap_used":24828188,"hdd_total":953159,"hdd_used":691802,"cpu":23.4,"network_tx":773327424,"network_rx":2685766416,"network_in":5371532832,"network_out":1546654849,"online4":true,"online6":false}
```

https://github.com/cppla/ServerStatus/blob/master/clients/client-linux.py
```py
  array['uptime'] = Uptime
                array['load_1'] = Load_1
                array['load_5'] = Load_5
                array['load_15'] = Load_15
                array['memory_total'] = MemoryTotal
                array['memory_used'] = MemoryUsed
                array['swap_total'] = SwapTotal
                array['swap_used'] = SwapTotal - SwapFree
                array['hdd_total'] = HDDTotal
                array['hdd_used'] = HDDUsed
                array['cpu'] = CPU
                array['network_rx'] = netSpeed.get("netrx")
                array['network_tx'] = netSpeed.get("nettx")
                array['network_in'] = NET_IN
                array['network_out'] = NET_OUT
                array['ping_10010'] = lostRate.get('10010') * 100
                array['ping_189'] = lostRate.get('189') * 100
                array['ping_10086'] = lostRate.get('10086') * 100
                array['time_10010'] = pingTime.get('10010')
                array['time_189'] = pingTime.get('189')
                array['time_10086'] = pingTime.get('10086')
                array['tcp'], array['udp'], array['process'], array['thread'] = tupd()
                array['io_read'] = diskIO.get("read")
                array['io_write'] = diskIO.get("write")
                array['custom'] = "<br>".join(f"{k}\\t解析: {v['dns_time']}\\t连接: {v['connect_time']}\\t下载: {v['download_time']}\\t在线率: <code>{v['online_rate']*100:.1f}%</code>" for k, v in monitorServer.items())
```


# 功能

## 客户端注册信息
实现客户端注册信息到服务器
然后定时更新客户端的信息


## 服务器端

### 获取信息
获取当前所有活跃客户端的信息
- 针对gpu服务器，添加gpu的具体信息

### 获取历史所有客户端的信息
实现客户端下线的功能，也就是，虽然当前活跃的系统里面没有，但是曾经有过





# windows 上编译


## 客户端
### 编译为exe
```
go build -o ./buildPackage/client_probe.exe client.go config.go client_info.go
client_probe.exe
```


### 打包为linux可执行
```cmd
set GOOS=linux
set GOARCH=amd64

go build -o ./buildPackage/client_linux_probe client.go config.go client_info.go

```



## 服务端
### 编译为exe
```
go build -o ./buildPackage/server_probe.exe find.go config.go client_info.go
client_probe.exe
```


### 打包为linux可执行
```cmd
set GOOS=linux
set GOARCH=amd64

go build -o ./buildPackage/server_linux_probe find.go config.go client_info.go

```
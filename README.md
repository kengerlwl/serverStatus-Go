# serverStatus-Go
 基于go的服务器探针系统


# run
## client
`go run client.go config.go client_info.go`

## find

`go run find.go config.go client_info.go`


# to do

- 修改配置文件的读取。让client可以设置机器名字，以及其他参数(done) 
- 服务器端添加web接口，主要是一个读取接口和一个删除某个服务器节点的接口
- 删除接口高危，需要鉴权。但是先不管
- etcd的设置状态也应该鉴权
- 加入前端界面对接。建议直接套用
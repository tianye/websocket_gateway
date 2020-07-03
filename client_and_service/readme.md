文件: 

MAC版本:websocket_gateway_mac.zip

Linux版本:websocket_gateway_linux_amd64.zip


使用需要解压

启动

MAC版本:
```
./websocket_gateway_mac -network_local_ip="127.0.0.1" -intranet_local_ip="127.0.0.1" -socket_listen_port="8301" -http_listen_port="8302" -callback_url_path="http://127.0.0.1:8808/service_callback.php"
```

Linux版本:
```
./websocket_gateway_linux_amd64 -network_local_ip="127.0.0.1" -intranet_local_ip="127.0.0.1" -socket_listen_port="8301" -http_listen_port="8302" -callback_url_path="http://127.0.0.1:8808/service_callback.php"
```



启动参数
```
 -callback_url_path string
        事件回调地址 (default "http://127.0.0.1:8808/service_callback.php")
  -conf string
        配置文件
  -http_listen_port string
        http监听端口 (default "8302")
  -intranet_local_ip string
        内网访问IP-没有请填写同外网IP (default "127.0.0.1")
  -network_local_ip string
        外网IP (default "127.0.0.1")
  -socket_listen_port string
        socket监听端口 (default "8301")
```

使用:
```
client_and_service文件夹

client.html     模拟客户端连接 js版本
client_php.php  模拟客户端连接 php版
client_go.go    模拟客户端连接 go版

service_api.php 模拟服务端主动调用接口
service_callback.php 模拟接受到(客户端消息和事件的处理)和当前gateway的事件处理
service_callback.log 接受到的日志

如果也是go服务的话也是直接请求接口和接受json事件就可以了.
解密管道ID的方法在文件:
common/structure/connection/connection.go:369
func DecodeConnection(connectionId string) (connectionInfo ConnectionInfo, err error)

```

打包:
```
Mac 下编译 Linux 和 Windows 64位可执行程序:
------------------------------------------------------------------
websocket_gateway_mac:
go1.14.4 build main.go
------------------------------------------------------------------
websocket_gateway_linux_amd64:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go1.14.4 build main.go
------------------------------------------------------------------
GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
GOARCH：目标平台的体系架构（386、amd64、arm）
交叉编译不支持 CGO 所以要禁用它
------------------------------------------------------------------
```

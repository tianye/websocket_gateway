package initialize

import (
	"encoding/json"
	"flag"
	"github.com/tianye/websocket_gateway/conf"
)

func InitConf() {
	confPath := ""
	confStruct := struct {
		NetworkLocalIp   string `json:"network_local_ip"`
		IntranetLocalIp  string `json:"intranet_local_ip"`
		SocketListenPort string `json:"socket_listen_port"`
		HttpListenPort   string `json:"http_listen_port"`
		CallbackUrlPath  string `json:"callback_url_path"`
	}{
		NetworkLocalIp:   "127.0.0.1",
		IntranetLocalIp:  "127.0.0.1",
		SocketListenPort: "8301",
		HttpListenPort:   "8302",
		CallbackUrlPath:  "http://127.0.0.1:8808/service_callback.php",
	}

	flag.StringVar(&confPath, "conf", "", "配置文件")

	if confPath != "" {
		json.Unmarshal([]byte(confPath), &confStruct)
	}

	flag.StringVar(&conf.NetworkLocalIp, "network_local_ip", confStruct.NetworkLocalIp, "外网IP")
	flag.StringVar(&conf.IntranetLocalIp, "intranet_local_ip", confStruct.IntranetLocalIp, "内网访问IP-没有填写同外网IP")
	flag.StringVar(&conf.SocketListenPort, "socket_listen_port", confStruct.SocketListenPort, "socket监听端口")
	flag.StringVar(&conf.HttpListenPort, "http_listen_port", confStruct.HttpListenPort, "http监听端口")
	flag.StringVar(&conf.CallbackUrlPath, "callback_url_path", confStruct.CallbackUrlPath, "事件回调地址")
	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Parse()

	//初始化配置
	conf.InitConf()
}

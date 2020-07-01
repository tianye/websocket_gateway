package conf

var (
	NetworkLocalIp   string //外网IP
	IntranetLocalIp  string //内网IP
	SocketListenPort string //socket监听端口
	HttpListenPort   string //http监听端口
	CallbackUrlPath	 string //事件回调地址
)

var ConfMap = make(map[string]string)

const NetworkLocalIP  = "NETWORK_LOCAL_IP"
const IntranetLocalIP = "INTRANET_LOCAL_IP"
const HttpPort 		  = "HTTP_PORT"
const WebSocketPort   = "SOCKET_PORT"
const CallbackUrl     = "CALLBACK_URL"

//每多少个用户下线触发一次清理通道GC操作
const ClearSocketConnectionNum = 1000

//第一次等待客户端请求超时时间
const FirstLinkActiveTimeOut = 60

//每次续期等待客户端下次响应超时时间
const EveryTimeActiveTimeOut = 600

//初始化CONF
func InitConf() {
	ConfMap[NetworkLocalIP]  = NetworkLocalIp
	ConfMap[IntranetLocalIP] = IntranetLocalIp
	ConfMap[WebSocketPort]   = SocketListenPort
	ConfMap[HttpPort]        = HttpListenPort
	ConfMap[CallbackUrl]     = CallbackUrlPath
}

//获取配置
func GetConf(key string) string {
	return ConfMap[key]
}

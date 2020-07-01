package main

import (
	_ "github.com/tianye/websocket_gateway/common/initialize" //初始化服务
	"github.com/tianye/websocket_gateway/conf"
	"github.com/tianye/websocket_gateway/request/http"
	"github.com/tianye/websocket_gateway/request/websocket"
	"github.com/tianye/websocket_gateway/service"
	"github.com/tianye/websocket_gateway/service/process_end"
	"github.com/tianye/websocket_gateway/service/push_worker"
	"os"
)

func main() {
	service.S.Log.Info("PUSH_GATEWAY-进程已启动", "PID:", os.Getpid())

	process_end.ProcessKillStart() //监听进程是否被杀死

	initHttp()           //Http监听
	initWebSocket()      //webSocket监听
	initProcessRunPush() //进程启动消息推送

	process_end.ProcessKillOver() //如果进程被杀死处理后事
}

func initHttp() {
	go http.Listen()
	service.S.Log.Info("PUSH_GATEWAY_HTTP-监听已启动", "端口:", conf.GetConf(conf.HttpPort))
}

func initWebSocket() {
	go websocket.Listen()
	service.S.Log.Info("PUSH_GATEWAY_SOCKET-监听已启动", "端口:", conf.GetConf(conf.WebSocketPort))
}

func initProcessRunPush() {
	push_worker.ProcessRun()
}

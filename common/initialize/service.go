package initialize

import (
	"github.com/labstack/gommon/log"
	"github.com/tianye/websocket_gateway/service"
)

func InitService() {
	//初始化服务
	service.S = new(service.Service)
	//日志
	service.S.Log = log.New("push_gateway")
}

package http

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/tianye/websocket_gateway/conf"
)

var e *echo.Echo

func apiRoute(e *echo.Echo) {
	//对内网API接口
	group := e.Group("api")
	group.POST("/push_connection_message", PushConnectionMessage)
	group.POST("/get_connection_is_online", GetConnectionIsOnline)
	group.POST("/kicked_out_connection", KickedOutConnection)
	group.POST("/get_online_num", GetOnlineNum)
}

func Listen() {
	e = echo.New()
	apiRoute(e)

	e.Server.Addr = ":" + conf.GetConf(conf.HttpPort)
	e.Logger.Error(gracehttp.Serve(e.Server))
}

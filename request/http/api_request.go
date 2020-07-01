package http

import (
	"github.com/labstack/echo"
	"github.com/tianye/websocket_gateway/common/structure/connection"
	"github.com/tianye/websocket_gateway/common/structure/http_struct"
	"github.com/tianye/websocket_gateway/controller"
	"github.com/tianye/websocket_gateway/service"
	"net/http"
)

//推送消息到管道
func PushConnectionMessage(c echo.Context) error {

	requestConnectionMessage := new(http_struct.RequestConnectionMessage)
	if err := c.Bind(requestConnectionMessage); err != nil {
		return err
	}

	responseData := controller.PushConnectionMessage(requestConnectionMessage.ConnectionId, []byte(requestConnectionMessage.PushMessage))

	service.S.Log.Info("请求方法:PushConnectionMessage", "-请求参数:", requestConnectionMessage, "-响应参数:", responseData)

	return c.JSON(http.StatusOK, responseData)
}

//获取管道在线状态
func GetConnectionIsOnline(c echo.Context) error {
	requestConnectionId := new(http_struct.RequestConnectionId)
	if err := c.Bind(requestConnectionId); err != nil {
		return err
	}

	responseData := controller.GetConnectionIsOnline(requestConnectionId.ConnectionId)

	service.S.Log.Info("请求方法:GetConnectionIsOnline", "-请求参数:", requestConnectionId, "-响应参数:", responseData)

	return c.JSON(http.StatusOK, responseData)
}

//踢出管道强制离线
func KickedOutConnection(c echo.Context) error {
	requestConnectionId := new(http_struct.RequestConnectionId)
	if err := c.Bind(requestConnectionId); err != nil {
		return err
	}

	responseData := controller.KickedOutConnection(requestConnectionId.ConnectionId)

	service.S.Log.Info("请求方法:KickedOutConnection", "-请求参数:", requestConnectionId, "-响应参数:", responseData)

	return c.JSON(http.StatusOK, responseData)
}

func GetOnlineNum(c echo.Context) error {
	responseData := map[string]int64{
		"max_connection_num":        connection.ConnectionCollection.ConnectionCounts.MaxConnectionNum,              //当前自增数
		"link_connection_num":       int64(connection.ConnectionCollection.ConnectionCounts.ConnectionLinkNum),      //链接数
		"wait_clear_connection_num": int64(connection.ConnectionCollection.ConnectionCounts.ConnectionWaitClearNum), //等待清理数
	}

	return c.JSON(http.StatusOK, responseData)
}

package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/tianye/websocket_gateway/common/structure/connection"
	"github.com/tianye/websocket_gateway/controller"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		conn   *connection.Connection
	)

	// 完成ws协议的握手操作 Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	//上线
	conn, err = onlineConnection(wsConn)

	//用户下线了
	defer closeConnection(conn)

	if err != nil {
		return
	}

	//消息监听
	err = messageConnection(conn)
	if err != nil {
		return
	}
}

//上线
func onlineConnection(wsConn *websocket.Conn) (conn *connection.Connection, err error) {
	if conn, err = connection.OnlineConnection(wsConn); err != nil {
		return
	}
	//上线事件触发
	controller.EventOnline(conn)

	return
}

//下线
func closeConnection(conn *connection.Connection) {
	conn.Close()

	//消息事件触发
	controller.EventOffline(conn)
}

//消息接受
func messageConnection(conn *connection.Connection) (err error) {
	var data = make([]byte, 0)

	for {
		if data, err = conn.ReadMessage(); err != nil {
			return
		}

		//消息事件触发
		controller.EventMessage(conn, data)
	}

	return
}

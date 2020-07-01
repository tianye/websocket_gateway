package websocket

import (
	"github.com/tianye/websocket_gateway/conf"
	"net/http"
)

func Listen() {
	http.HandleFunc("/web_socket", wsHandler)
	http.ListenAndServe(":"+conf.GetConf(conf.WebSocketPort), nil)
}

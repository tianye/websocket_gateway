package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

var ws *websocket.Conn //websocket连接
var messageChan = make(chan string, 20) //简单的限制处理速度
var addr = flag.String("addr", "127.0.0.1:8301", "http service address") //socket地址

func main() {
	//初始化socket连接
	initWebSocket()

	//发送消息
	WriteWebSocketValue([]byte("hello"))

	//监听socket
	listen()
}

//初始化WebSocket
func initWebSocket() error {
	var err error

	u := url.URL{Scheme: "ws", Host: *addr, Path: "web_socket"}
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		return err
	}

	return nil
}

//获取WebSocket单例
func GetWebSocketConn() (*websocket.Conn, error) {
	if nil == ws {
		err := initWebSocket()
		if nil != err {
			return nil, err
		}
	}
	return ws, nil
}

//获取socket数据
func GetWebSocketValue() (messageType int, p []byte, err error) {
	return ws.ReadMessage()
}

//发送socket数据
func WriteWebSocketValue(message [] byte) (error) {
	err := ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return err
}

//监听端口
func listen() {
	go requestMessage()
	//获取socket里面的数据
	for {
		_, message, err := GetWebSocketValue()
		if err != nil {
			continue
		}

		messageChan <- string(message)
	}
}

//接受的消息处理
func requestMessage() {
	for {
		messageInfo := <-messageChan

		doCommit(string(messageInfo))
	}
}

//根据接受的消息做什么
func doCommit(messageInfo string) {
	fmt.Print("接受到的消息:", messageInfo)
}

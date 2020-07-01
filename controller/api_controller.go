package controller

import (
	"fmt"
	"github.com/tianye/websocket_gateway/common/response_code"
	"github.com/tianye/websocket_gateway/common/structure/connection"
	"github.com/tianye/websocket_gateway/common/structure/http_struct"
	"github.com/tianye/websocket_gateway/conf"
	"time"
)

//推送消息到管道
func PushConnectionMessage(connectionId string, pushMessage []byte) (responseData http_struct.Response) {

	responseData = http_struct.Response{Code: response_code.FAIL_CODE, Data: http_struct.Message{Message: "消息推送失败"}}

	//获取到了管道地址
	if connectionItem, err := connection.GetConnectionLink(connectionId); err == nil {
		//写入管道信息
		if err := connectionItem.WriteMessage(pushMessage); err == nil {
			responseData = http_struct.Response{Code: response_code.SUCCESS_CODE, Data: http_struct.Message{Message: "消息推送成功"}}

			//续期心跳检测
			renewalHeartbeatTimer(connectionItem, conf.EveryTimeActiveTimeOut*time.Second)
		} else {
			responseData = http_struct.Response{Code: response_code.EquipmentPushFail, Data: http_struct.Message{Message: "消息推送失败", ErrorMessage: string(fmt.Sprintf("%s", err))}}
		}
	} else if connectionItem == nil {
		responseData = http_struct.Response{Code: response_code.EquipmentNotOnline, Data: http_struct.Message{Message: "客户端不在线"}}
	}

	return
}

//获取管道在线状态
func GetConnectionIsOnline(connectionId string) (responseData http_struct.Response) {

	responseData = http_struct.Response{Code: response_code.FAIL_CODE, Data: http_struct.Message{Message: "获取客户端在线情况失败"}}

	if connectionItem, err := connection.GetConnectionLink(connectionId); err == nil {
		responseData = http_struct.Response{Code: response_code.SUCCESS_CODE, Data: http_struct.Message{Message: "客户端在线"}}
	} else if connectionItem == nil {
		responseData = http_struct.Response{Code: response_code.EquipmentNotOnline, Data: http_struct.Message{Message: "客户端不在线"}}
	} else {
		responseData = http_struct.Response{Code: response_code.FAIL_CODE, Data: http_struct.Message{Message: "获取客户端在线情况失败", ErrorMessage: string(fmt.Sprintf("%s", err))}}
	}

	return
}

//踢出管道强制离线
func KickedOutConnection(connectionId string) (responseData http_struct.Response) {

	responseData = http_struct.Response{Code: response_code.FAIL_CODE, Data: http_struct.Message{Message: "踢出客户端失败"}}

	if connectionItem, err := connection.GetConnectionLink(connectionId); err == nil {
		//关闭管道
		connectionItem.Close()
		responseData = http_struct.Response{Code: response_code.SUCCESS_CODE, Data: http_struct.Message{Message: "已将客户端提下线"}}
	} else if connectionItem == nil {
		responseData = http_struct.Response{Code: response_code.SUCCESS_CODE, Data: http_struct.Message{Message: "客户已离线"}}
	} else {
		responseData = http_struct.Response{Code: response_code.FAIL_CODE, Data: http_struct.Message{Message: "踢出客户端失败", ErrorMessage: string(fmt.Sprintf("%s", err))}}
	}

	return
}

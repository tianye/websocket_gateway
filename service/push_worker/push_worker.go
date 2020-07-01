package push_worker

import (
	"encoding/json"
	"github.com/tianye/websocket_gateway/common/structure/connection"
	"github.com/tianye/websocket_gateway/common/structure/push_worker"
	"github.com/tianye/websocket_gateway/common/tools"
	"github.com/tianye/websocket_gateway/conf"
	"github.com/tianye/websocket_gateway/service"
	"github.com/tianye/websocket_gateway/service/client_worker"
	"os"
	"strings"
)

const ConstEventOnline = "EventOnline"   //上线事件
const ConstEventOffline = "EventOffline" //下线事件
const ConstEventMessage = "EventMessage" //消息事件

const ConstEventProcessRun = "EventProcessRun"             //进程启动
const ConstEventProcessStartKill = "EventProcessStartKill" //进程开始杀死消息
const ConstEventProcessEndKill = "EventProcessEndKill"     //进程结束杀死消息

//客户端事件

//上线业务工作者事件处理RPC
func EventOnline(connectionId string, eventData string) (responseInfo *push_worker.ResponseInfo, err error) {
	responseInfo, err = client_worker.EventOnline(&push_worker.EventInfo{
		GatewayIp:    conf.GetConf(conf.IntranetLocalIP),
		HttpPort:     conf.GetConf(conf.HttpPort),
		SocketPort:   conf.GetConf(conf.WebSocketPort),
		ConnectionId: connectionId,
		EventTime:    tools.GetNowTimeUnix(),
		EventType:    ConstEventOnline,
		EventData:    eventData,
		DataLength:   int64(strings.Count(eventData, "")),
	})

	if err != nil {
		service.S.Log.Error(err)

		return nil, err
	}

	return responseInfo, nil
}

//下线业务工作者事件处理RPC
func EventOffline(connectionId string, eventData string) (responseInfo *push_worker.ResponseInfo, err error) {
	responseInfo, err = client_worker.EventOffline(&push_worker.EventInfo{
		GatewayIp:    conf.GetConf(conf.IntranetLocalIP),
		HttpPort:     conf.GetConf(conf.HttpPort),
		SocketPort:   conf.GetConf(conf.WebSocketPort),
		ConnectionId: connectionId,
		EventTime:    tools.GetNowTimeUnix(),
		EventType:    ConstEventOffline,
		EventData:    eventData,
		DataLength:   int64(strings.Count(eventData, "")),
	})

	if err != nil {
		service.S.Log.Error(err)

		return nil, err
	}

	return responseInfo, nil
}

//消息业务工作者事件处理
func EventMessage(connectionId string, eventData string) (responseInfo *push_worker.ResponseInfo, err error) {
	responseInfo, err = client_worker.EventMessage(&push_worker.EventInfo{
		GatewayIp:  conf.GetConf(conf.IntranetLocalIP),
		HttpPort:   conf.GetConf(conf.HttpPort),
		SocketPort: conf.GetConf(conf.WebSocketPort),

		ConnectionId: connectionId,
		EventTime:    tools.GetNowTimeUnix(),
		EventType:    ConstEventMessage,
		EventData:    eventData,
		DataLength:   int64(strings.Count(eventData, "")),
	})

	if err != nil {
		service.S.Log.Error(err)

		return nil, err
	}

	return responseInfo, err
}

//服务端进程事件

//PUSH_GATEWAY-进程已经结束-等待释放", "\n死亡原因:
func ProcessKillStart(sig string) (responseInfo *push_worker.ResponseInfo, err error) {
	data := struct {
		Sig    string `json:"sig"`
		Pid    int    `json:"pid"`
		Prefix string `json:"prefix"`
	}{
		Sig:    sig,
		Pid:    os.Getpid(),
		Prefix: connection.ConnectionCollection.Prefix,
	}

	JsonEventData, _ := json.Marshal(data)
	eventData := string(JsonEventData)

	responseInfo, err = client_worker.ProcessMessage(&push_worker.EventInfo{
		GatewayIp:    conf.GetConf(conf.IntranetLocalIP),
		HttpPort:     conf.GetConf(conf.HttpPort),
		SocketPort:   conf.GetConf(conf.WebSocketPort),
		ConnectionId: "0",
		EventTime:    tools.GetNowTimeUnix(),
		EventType:    ConstEventProcessStartKill,
		EventData:    eventData,
		DataLength:   int64(strings.Count(eventData, "")),
	})

	if err != nil {
		service.S.Log.Error(err)

		return nil, err
	}

	return responseInfo, err
}

//PUSH_GATEWAY-进程已经结束-完全关闭
func ProcessKillOver() (responseInfo *push_worker.ResponseInfo, err error) {
	data := struct {
		Pid    int    `json:"pid"`
		Prefix string `json:"prefix"`
	}{
		Pid:    os.Getpid(),
		Prefix: connection.ConnectionCollection.Prefix,
	}

	JsonEventData, _ := json.Marshal(data)
	eventData := string(JsonEventData)

	responseInfo, err = client_worker.ProcessMessage(&push_worker.EventInfo{
		GatewayIp:    conf.GetConf(conf.IntranetLocalIP),
		HttpPort:     conf.GetConf(conf.HttpPort),
		SocketPort:   conf.GetConf(conf.WebSocketPort),
		ConnectionId: "0",
		EventTime:    tools.GetNowTimeUnix(),
		EventType:    ConstEventProcessEndKill,
		EventData:    eventData,
		DataLength:   int64(strings.Count(eventData, "")),
	})

	if err != nil {
		service.S.Log.Error(err)

		return nil, err
	}

	return responseInfo, err
}

//PUSH_GATEWAY-进程启动
func ProcessRun() (responseInfo *push_worker.ResponseInfo, err error) {
	data := struct {
		Pid    int    `json:"pid"`
		Prefix string `json:"prefix"`
	}{
		Pid:    os.Getpid(),
		Prefix: connection.ConnectionCollection.Prefix,
	}

	JsonEventData, _ := json.Marshal(data)
	eventData := string(JsonEventData)

	responseInfo, err = client_worker.ProcessMessage(&push_worker.EventInfo{
		GatewayIp:    conf.GetConf(conf.IntranetLocalIP),
		HttpPort:     conf.GetConf(conf.HttpPort),
		SocketPort:   conf.GetConf(conf.WebSocketPort),
		ConnectionId: "0",
		EventTime:    tools.GetNowTimeUnix(),
		EventType:    ConstEventProcessRun,
		EventData:    eventData,
		DataLength:   int64(strings.Count(eventData, "")),
	})

	if err != nil {
		service.S.Log.Error(err)

		return nil, err
	}

	return responseInfo, err
}

package controller

import (
	"github.com/tianye/websocket_gateway/common/response_code"
	"github.com/tianye/websocket_gateway/common/structure/connection"
	"github.com/tianye/websocket_gateway/conf"
	"github.com/tianye/websocket_gateway/service"
	"github.com/tianye/websocket_gateway/service/push_worker"
	"time"
)

//上线事件
func EventOnline(conn *connection.Connection) {
	service.S.Log.Info("用户上线啦:", conn.ConnectionId)

	responseInfo, err := push_worker.EventOnline(conn.ConnectionId, "")

	if err != nil {
		service.S.Log.Error("上线事件RPC处理失败:", conn.ConnectionId)
		//上线事件处理失败 直接离线
		conn.Close()
		return
	}

	if responseInfo.ResponseCode == response_code.SUCCESS_CODE {
		//响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
	} else if responseInfo.ResponseCode == response_code.FAIL_CODE {
		//响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
		//如果响应状态为失败 则用户直接下线
		conn.Close()

		return
	} else if responseInfo.ResponseCode == response_code.AUTHENTICATION_FAILED {
		//响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
		//如果响应状态为失败 则用户直接下线
		conn.Close()
	}

	//response_code.SUCCESS_NULL_CODE 如果是当前code204 则表示处理成功但是不返还任何消息内容

	//启动心跳检测
	startHeartbeatTimer(conn)

	return
}

//下线事件
func EventOffline(conn *connection.Connection) {
	//下线时强制清理定时器-应为用户已经下线了 所以直接清除定时器
	delHeartbeatTimer(conn, true)

	service.S.Log.Info("用户下线啦:", conn.ConnectionId)

	responseInfo, err := push_worker.EventOffline(conn.ConnectionId, "")

	if err != nil {
		service.S.Log.Error("下线事件RPC处理失败:", conn.ConnectionId)
		return
	}

	if responseInfo.ResponseCode == response_code.SUCCESS_CODE {
		//响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
	}

	//response_code.SUCCESS_NULL_CODE 如果是当前code204 则表示处理成功但是不返还任何消息内容
	//response_code.FAIL_CODE   吐血 就算是worker 处理下线消息失败了 但是用户已经下线了已无力回天

	return
}

//消息事件
func EventMessage(conn *connection.Connection, data []byte) {
	service.S.Log.Info("用户发来消息啦:", conn.ConnectionId, "消息:", string(data))
	//续期心跳检测
	renewalHeartbeatTimer(conn, conf.EveryTimeActiveTimeOut*time.Second)

	responseInfo, err := push_worker.EventMessage(conn.ConnectionId, string(data))

	if err != nil {
		service.S.Log.Error("消息事件RPC处理失败:", conn.ConnectionId)
		return
	}

	if responseInfo.ResponseCode == response_code.SUCCESS_CODE {
		//成功-响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
	} else if responseInfo.ResponseCode == response_code.FAIL_CODE {
		//失败-响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
	} else if responseInfo.ResponseCode == response_code.AUTHENTICATION_FAILED {
		//响应消息
		conn.WriteMessage([]byte(responseInfo.ResponseData))
		//如果响应状态为失败 则用户直接下线
		conn.Close()
	}

	//response_code.SUCCESS_NULL_CODE 如果是当前code204 则表示处理成功但是不返还任何消息内容

	return
}

//启动心跳定时器
func startHeartbeatTimer(conn *connection.Connection) {
	//增加定时检测第一次鉴活
	conn.Timer = time.NewTimer(time.Second * conf.FirstLinkActiveTimeOut)
	go func() {
		for {
			select {
			case <-conn.Timer.C:
				//到达触发时间 如果是非活跃的踢下线 用 Reset 不要用 Stop 不然不会触发这里的操作 定时器不会释放
				if conn.IsBreak == false {
					//如果还没端口但是超时了就断开它
					conn.Close()
					service.S.Log.Warn("这个用户超时离线了:", conn.ConnectionId)
				}
			}

			return
		}
	}()
}

//删除心跳定时器
func delHeartbeatTimer(conn *connection.Connection, mandatory bool) {

	if conn.Timer == nil {
		return
	}

	if mandatory == true {
		//如果是强制的只删除定时器
		conn.Timer.Reset(0)
	} else if conn.IsBreak == false {
		service.S.Log.Warn("用户触发了活跃状态删除了心跳检测:", conn.ConnectionId)
		//如果非活跃-删除定时器的时候直接踢掉用户
		conn.IsBreak = true //没错不要动就是 true
		conn.Timer.Reset(0)
	}
}

//续期心跳定时检测
func renewalHeartbeatTimer(conn *connection.Connection, d time.Duration) {

	if conn.Timer == nil {
		startHeartbeatTimer(conn)

		return
	}

	service.S.Log.Warn("用户续期了心跳检测时间:", conn.ConnectionId)
	conn.IsBreak = false //没错不要动就是 false
	conn.Timer.Reset(d)
}

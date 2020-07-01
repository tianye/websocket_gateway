package process_end

import (
	"github.com/tianye/websocket_gateway/common/structure/connection"
	"github.com/tianye/websocket_gateway/service/push_worker"
	"log"
	"sync"
)

var wg sync.WaitGroup

func closeConnectionAll() {
	//锁住一切 准备下线所以客户端
	needKillLen := len(connection.ConnectionCollection.ConnectionList)

	log.Println("KILL-需要下线的用户数量:", needKillLen)

	connection.ConnectionCollection.Mutex.Lock()
	defer connection.ConnectionCollection.Mutex.Unlock()

	for ConnectionId, ConnectionItem := range connection.ConnectionCollection.ConnectionList {
		wg.Add(1)
		//触发用户下线事件
		go closeConnectionItem(ConnectionItem, ConnectionId)
	}

	wg.Wait()
}

func closeConnectionItem(ConnectionItem *connection.Connection, ConnectionId string) {

	//通知worker 此用户已下线
	push_worker.EventOffline(ConnectionItem.ConnectionId, "")

	log.Println("KILL-下线客户端:", ConnectionId)

	defer wg.Done()
}

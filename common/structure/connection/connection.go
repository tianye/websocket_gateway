package connection

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/tianye/websocket_gateway/common/tools"
	"github.com/tianye/websocket_gateway/common/tools/pack"
	"github.com/tianye/websocket_gateway/common/tools/range_num"
	"github.com/tianye/websocket_gateway/conf"
	"github.com/tianye/websocket_gateway/service"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type ConnectionInfo struct {
	Ip           string `json:"ip"`            //网关连接ip
	HttpPort     string `json:"http_port"`     //http端口
	Pid          string `json:"pid"`           //进程id
	Cid          string `json:"cid"`           //自增id
	Range        string `json:"range"`         //随机1-65535
	ConnectionId string `json:"connection_id"` //加密后的管道id
}

//通道计数
type ConnectionCounts struct {
	MaxConnectionNum int64 //当前自增的通道记录

	ConnectionLinkNum      int //当前连接中的通道数量
	ConnectionWaitClearNum int //当前待清理的删除通道数量

	mutex sync.Mutex //防止重复清理MAP
}

//链接MAP结构体
type ConnectionMap struct {
	ConnectionList map[string]*Connection //所有的socket通道信息记录

	TmpConnectionList map[string]*Connection //当IsClearIng为true的时候 会把链接进来的用户存到当前MAP中一份 用来异步清理MAP

	DeletedWaitReuse *ConnectionQueue //已删除可以复用的通道数

	Prefix string //前缀

	ConnectionCounts *ConnectionCounts //计数器

	IsClearIng bool //是否正在清理

	Mutex sync.Mutex //锁防止生成重复的通道key值
}

//链接信息
type Connection struct {
	ConnectionId string //链接信息唯一ID
	wsConnect    *websocket.Conn
	inChan       chan []byte
	outChan      chan []byte
	closeChan    chan byte

	Timer   *time.Timer //心跳检测超时器
	IsBreak bool        //是否活跃

	mutex    sync.Mutex // 对closeChan关闭上锁
	isClosed bool       // 防止closeChan被关闭多次
}

var packProtocol = pack.Protocol{}

//通道记录集合
var ConnectionCollection = ConnectionMap{}

//清除通道记录信号
var ConnectionClearSignal = make(chan bool, 100)

//初始化前缀
func InitConnection() {

	//ip2lang http_port rpc_port pid
	packProtocol.Format = []string{"N8", "N2", "N2"}
	//初始化前缀变量
	ip2long := int64(tools.Ip2long(net.ParseIP(conf.GetConf(conf.IntranetLocalIP))))
	httpPort, _ := strconv.ParseInt(conf.GetConf(conf.HttpPort), 10, 64)
	pId := int64(os.Getpid())

	packPrefix := packProtocol.Pack16(ip2long, httpPort, pId)

	//初始化全局变量
	ConnectionCollection = ConnectionMap{
		Prefix:            packPrefix,
		ConnectionList:    make(map[string]*Connection),
		TmpConnectionList: make(map[string]*Connection),
		DeletedWaitReuse:  NewQueue(),
		ConnectionCounts:  &ConnectionCounts{ConnectionLinkNum: 0, ConnectionWaitClearNum: 0, MaxConnectionNum: 0},
		IsClearIng:        false,
	}

	//增加清除通道MAP-WORKER
	go clearCollectionMapSignalListen()
}

//初始化连接信息
func OnlineConnection(wsConn *websocket.Conn) (conn *Connection, err error) {

	conn = &Connection{
		wsConnect: wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
		IsBreak:   false,
	}

	conn.makeConnectionId()

	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()

	return
}

//获取管道记录
func GetConnectionLink(connectionId string) (conn *Connection, err error) {
	if connectionItem, ok := ConnectionCollection.ConnectionList[connectionId]; ok == true {
		return connectionItem, nil
	}

	if connectionItem, ok := ConnectionCollection.TmpConnectionList[connectionId]; ok == true {
		return connectionItem, nil
	}

	return nil, errors.New("NOT CONNECTION")
}

//Connection-读取链接内容
func (conn *Connection) ReadMessage() (data []byte, err error) {

	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

//Connection-写入链接内容
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

//Connection-关闭链接
func (conn *Connection) Close() {

	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClosed {
		//从map中移除
		conn.destructionConnectionId()
		//关闭chan管道
		close(conn.closeChan)

		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

//Connection-写锁
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.inChan <- data:
		case <-conn.closeChan: // closeChan 感知 conn断开
			goto ERR
		}

	}

ERR:
	conn.Close()
}

//Connection-读锁
func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}

		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

//生成通道唯一ID
func (conn *Connection) makeConnectionId() (connectionId string) {

	//加锁 - 之加计数器锁即可 不用加主map锁 切记 不然会导致 主map无法被读取和写入更新等操作
	ConnectionCollection.ConnectionCounts.mutex.Lock()

	//解锁
	defer func() {
		ConnectionCollection.ConnectionCounts.mutex.Unlock()
	}()

	//获取队中最先进来的一个可用的管道ID
	reuseConnection, reuseErr := ConnectionCollection.DeletedWaitReuse.DeQueue()

	//如果存在错误 或者 复用通道为空
	if reuseErr != nil || reuseConnection.connectionId == "" {
		//增加通道数
		ConnectionCollection.ConnectionCounts.MaxConnectionNum++
		//当前链接用户数增加
		ConnectionCollection.ConnectionCounts.ConnectionLinkNum++
		//通道ID = 前缀(base64 本机IP)+进程ID+通道递增数

		packProtocol.Format = []string{"N2", "N4"}
		rangeNum := int64(range_num.GenerateRangeNum(1, 65535))
		ConnectionNumString := packProtocol.Pack16(rangeNum, ConnectionCollection.ConnectionCounts.MaxConnectionNum)

		connectionId = ConnectionCollection.Prefix + ConnectionNumString
	} else {
		//如果存在可以复用的通道直接使用
		connectionId = reuseConnection.connectionId
	}

	//标记唯一ID
	conn.ConnectionId = connectionId
	//记录唯一ID HASH_MAP 中便于查找
	ConnectionCollection.ConnectionList[connectionId] = conn

	//如果正在清理MAP中
	if ConnectionCollection.IsClearIng == true {
		ConnectionCollection.TmpConnectionList[connectionId] = conn
	}

	//返回通道唯一ID
	return connectionId
}

//删除MAP值-但是不会释放内存
func (conn *Connection) destructionConnectionId() bool {

	//加锁
	ConnectionCollection.ConnectionCounts.mutex.Lock()

	//解锁
	defer func() {
		ConnectionCollection.ConnectionCounts.mutex.Unlock()
	}()

	//如果HASH_MAP中存在管道记录
	if _, ok := ConnectionCollection.ConnectionList[conn.ConnectionId]; ok == true {

		//从HASH_MAP中移除
		delete(ConnectionCollection.ConnectionList, conn.ConnectionId)

		//如果正在清理MAP中离线了
		if ConnectionCollection.IsClearIng == true {
			//并且临时记录中存在
			if _, ok := ConnectionCollection.TmpConnectionList[conn.ConnectionId]; ok == true {
				//从临时记录中移除当前链接通道信息
				delete(ConnectionCollection.TmpConnectionList, conn.ConnectionId)
			}
		}

		//减少当前链接记录数
		ConnectionCollection.ConnectionCounts.ConnectionLinkNum--
		ConnectionCollection.ConnectionCounts.ConnectionWaitClearNum++
		//如果触达了清除指标数-发送信号量清除
		if ConnectionCollection.ConnectionCounts.ConnectionWaitClearNum >= conf.ClearSocketConnectionNum {
			//标记为清除中
			ConnectionCollection.IsClearIng = true
			//待清除记录数清0
			ConnectionCollection.ConnectionCounts.ConnectionWaitClearNum = 0
			//发送给清除信号
			ConnectionClearSignal <- ConnectionCollection.IsClearIng
		}

		//添加到重复可使用通道ID中
		ConnectionCollection.DeletedWaitReuse.EnQueue(&ConnectionQueueBody{connectionId: conn.ConnectionId, connectionDelTime: tools.GetNowTimeUnix() + ConnectionNextAvailable.connectionNextReuseTime})
	}

	return true
}

//清除管道MAP信号监听
func clearCollectionMapSignalListen() {
	for {
		clearSignal := <-ConnectionClearSignal

		//如果清除信号的false 则不清除跳过
		if false == clearSignal {
			continue
		}

		clearCollectionMap()
	}
}

func clearCollectionMap() {
	logId := tools.GetLogId()

	service.S.Log.Info("开始清理MAP:", logId)
	//生成一个新的MAP
	newMap := make(map[string]*Connection)

	//变量旧的MAP复制到新的MAP上
	for key, val := range ConnectionCollection.ConnectionList {
		newMap[key] = val
	}

	//加锁 这个时候需要主MAP和计数器同时加锁 - 要把清除过程中的tmpConnectionList放到新的map中 替换主map
	ConnectionCollection.Mutex.Lock()
	ConnectionCollection.ConnectionCounts.mutex.Lock()

	//结束后释放锁
	defer func() {
		ConnectionCollection.ConnectionCounts.mutex.Unlock()
		ConnectionCollection.Mutex.Unlock()
	}()

	//把清除过程中-新增进来的链接通道合并到新的MAP中
	for key, val := range ConnectionCollection.TmpConnectionList {
		newMap[key] = val
	}

	//删除在清除过程中-临时存储的链接通道 老的map没有人使用会被异步GC
	ConnectionCollection.TmpConnectionList = make(map[string]*Connection)

	//清空MAP赋值新的map 老的map没有人使用会被异步GC
	ConnectionCollection.ConnectionList = newMap

	//清除完标记为清理完成
	ConnectionCollection.IsClearIng = false

	service.S.Log.Info("MAP清理结束", logId)
}

//解密管道信息
func DecodeConnection(connectionId string) (connectionInfo ConnectionInfo, err error) {

	if len(connectionId) != 36 {
		return connectionInfo, errors.New("管道ID非法")
	}

	packProtocol := pack.Protocol{}
	packProtocol.Format = []string{"N8", "N2", "N2", "N2", "N4"}
	unConnection := packProtocol.UnPack16(connectionId)

	connectionInfo = ConnectionInfo{
		Ip:           tools.Long2ip(unConnection[0]).String(),
		HttpPort:     strconv.FormatInt(unConnection[1], 10),
		Pid:          strconv.FormatInt(unConnection[3], 10),
		Range:        strconv.Itoa(range_num.GenerateRangeNum(1, 65535)),
		Cid:          strconv.FormatInt(unConnection[4], 10),
		ConnectionId: connectionId,
	}

	return connectionInfo, err
}

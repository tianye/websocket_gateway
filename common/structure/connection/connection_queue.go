package connection

import (
	"errors"
	"github.com/tianye/websocket_gateway/common/tools"
	"github.com/tianye/websocket_gateway/service"
	"sync"
)

//一个双线链表 队列操作 头插 尾取

//下线ID小于10000个则不使用
const QueueServedLength = 10000

type ConnectionQueueBody struct {
	connectionId      string //管道ID
	connectionDelTime int64  //管道ID下次可用时间
}

type ConnectionQueueItem struct {
	ConnectionQueueBody *ConnectionQueueBody //管道的内容

	Prev *ConnectionQueueItem //链表上一个元素
	Next *ConnectionQueueItem //链表下一个元素
}

type ConnectionQueue struct {
	Header *ConnectionQueueItem //第一个元素 减少复杂度
	Tail   *ConnectionQueueItem //最后一个元素 减少复杂度

	queueLen int //队的长度

	queueLock sync.RWMutex //队列操作锁

	lock sync.RWMutex //struct锁
}

//管道锁检测
var ConnectionNextAvailable = struct {
	availableTimeStart          bool
	connectionNextAvailableTime int64
	connectionNextReuseTime     int64
	queueLen                    int
}{
	availableTimeStart:          false,       //是否开启检测可用时间
	connectionNextAvailableTime: int64(0),    //管道下次可用时间-不准确 只是为了减少锁操作
	connectionNextReuseTime:     int64(1800), //管道失效后多少秒后可用
	queueLen:                    0,           //此值也是不准确的 准确的请依照ConnectionQueue.queueLen 为准 当前只是为了减少锁操作
}

// 创建队
func NewQueue() *ConnectionQueue {
	return &ConnectionQueue{queueLen: 0}
}

// 入队-头插
func (s *ConnectionQueue) EnQueue(body *ConnectionQueueBody) bool {

	s.lock.Lock()
	defer s.lock.Unlock()

	//入队操作
	s.addQueue(body)

	//如果是第一次push 开始启动可用检测
	if ConnectionNextAvailable.availableTimeStart == false {
		ConnectionNextAvailable.availableTimeStart = true
		ConnectionNextAvailable.connectionNextAvailableTime = tools.GetNowTimeUnix()

		service.S.Log.Info("这是第一次启动下线管道ID时间检测复用")
	}

	ConnectionNextAvailable.queueLen = s.queueLen

	return true
}

//真实入队操作
func (s *ConnectionQueue) addQueue(body *ConnectionQueueBody) bool {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	if s.queueLen <= 0 {
		//不可能小于0的
		s.queueLen = 0
		newItem := &ConnectionQueueItem{ConnectionQueueBody: body, Next: nil, Prev: nil}
		//应为是第一个元素所以头尾都是自己
		s.Header = newItem
		s.Tail = newItem
	} else if s.queueLen == 1 {
		//应为是第一个元素
		newItem := &ConnectionQueueItem{ConnectionQueueBody: body}
		s.Header.Prev = newItem
		s.Header = newItem
	} else {
		//第二个以上的元素
		newItem := &ConnectionQueueItem{ConnectionQueueBody: body}
		s.Header.Prev = newItem
		newItem.Next = s.Header
		s.Header = newItem
	}

	s.queueLen++

	return true
}

//真实操作出队操作
func (s *ConnectionQueue) outQueue() (responseItem *ConnectionQueueBody, err error) {
	s.queueLock.Lock()

	defer s.queueLock.Unlock()

	if s.Tail == nil {
		return responseItem, errors.New("QUEUE IS EMPTY")
	}

	//尾取
	responseConnectionQueue := s.Tail

	//真实存储的内容
	responseItem = responseConnectionQueue.ConnectionQueueBody

	//下个一Header
	prevItem := responseConnectionQueue.Prev

	if prevItem != nil {
		//上一个元素已经被移除
		prevItem.Next = nil
	}

	s.Tail = prevItem

	//移除队列
	s.queueLen--

	return responseItem, nil
}

// 出队-尾取
func (s *ConnectionQueue) DeQueue() (responseItem *ConnectionQueueBody, err error) {
	//如果没有push则不可用
	if ConnectionNextAvailable.availableTimeStart == false {
		return responseItem, errors.New("NOT AVAILABLE")
	}

	//如果长度小于QueueServedLength个则不使用
	if ConnectionNextAvailable.queueLen <= QueueServedLength {
		return responseItem, errors.New("QUEUE LEN ZERO")
	}

	//如果可用时间大于当前时间
	if ConnectionNextAvailable.connectionNextAvailableTime > tools.GetNowTimeUnix() {
		return responseItem, errors.New("NOT TO THE TIME AVAILABLE")
	}

	//如果以上三个条件都满足了 则开始锁了 获取一个可复用的管道ID
	s.lock.Lock()

	//结束时释放锁
	defer s.lock.Unlock()

	//如果小于等于0表示队中无数据
	if s.queueLen <= 0 {
		return responseItem, errors.New("QUEUE IS NULL")
	}

	//读取一下队中最先应该获取的值
	item := s.Tail.ConnectionQueueBody

	//如果这个值没有到生效时间 则不返回当前值
	if item.connectionDelTime > tools.GetNowTimeUnix() {
		//下一个可用时间
		ConnectionNextAvailable.connectionNextAvailableTime = item.connectionDelTime

		return responseItem, errors.New("QUEUE IS NOT AVAILABLE TIME")
	}

	//出队
	responseItem, err = s.outQueue()

	ConnectionNextAvailable.queueLen = s.queueLen

	//如果这个值到了当前时间 在队移除 返回当前值
	return responseItem, nil
}

//获取队的长度
func (s *ConnectionQueue) GetLen() int {
	return s.queueLen
}

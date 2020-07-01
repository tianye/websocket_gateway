package process_end

import (
	"github.com/tianye/websocket_gateway/conf"
	"github.com/tianye/websocket_gateway/service/push_worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var sigs = make(chan os.Signal, 1)
var done = make(chan bool, 1)

func ProcessKillStart() {
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		push_worker.ProcessKillStart(sig.String())


		log.Println("PUSH_GATEWAY_SOCKET-进程被杀死", "PID:", os.Getpid(), "端口:", conf.GetConf(conf.WebSocketPort), "死亡原因:", sig.String())
		done <- true
	}()
}

func ProcessKillOver() {
	<-done
	log.Println("PUSH_GATEWAY-正在释放客户端链接-请稍后", "PID:", os.Getpid(), "端口:", conf.GetConf(conf.WebSocketPort))
	//进程结束之前
	closeConnectionAll()

	push_worker.ProcessKillOver()

	log.Println("PUSH_GATEWAY-进程已结束", "PID:", os.Getpid(), "端口:", conf.GetConf(conf.WebSocketPort))

	//退出
	os.Exit(0)
}

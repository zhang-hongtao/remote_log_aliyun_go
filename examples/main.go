package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"remote_log_aliyun_go"
	"syscall"
	"time"
)

/**
 * 测试
 */
func main() {
	log := remote_log_aliyun_go.NewLogger("winner-test-project")
	err := log.Init()
	if err != nil {
		fmt.Print("初始化错误:", err.Error())
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		time.Sleep(5 * time.Second)
		GracefullExit(ch)
	}()
	log.Error(errors.New("错误日志6"))
	log.Info("测试6")
	for {
		s := <-ch
		switch s {
		case syscall.SIGINT:
			//SIGINT 信号，在程序关闭时会收到这个信号
			fmt.Println("SIGINT", "退出程序，执行退出前逻辑")
			log.Access("测试Access6")
			log.Close()
			// time.Sleep(5 * time.Second)
			fmt.Println("system end")
			os.Exit(0)
		case syscall.SIGKILL:
			fmt.Println("SIGKILL")
		default:
			fmt.Println("default")
		}
	}
}

func GracefullExit(ch chan os.Signal) {
	ch <- syscall.SIGINT
}

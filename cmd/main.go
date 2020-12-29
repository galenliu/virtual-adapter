package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yeelight"
)



func main() {

	var adapter = yeelight.NewYeeAdapter()

	var systemCallCloseFunc = func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-c
		if adapter != nil {
			adapter.CloseProxy()
		}
		os.Exit(0)
	}

	go systemCallCloseFunc()

	for {
		if adapter.ProxyRunning() {
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Print("main running .....\r\n")
			continue
		} else {
			return
		}
	}

}

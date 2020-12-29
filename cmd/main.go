package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"virtual_adpater"
)

func main() {

	var yeelightAdapter = virtual_adapter.NewYeeAdapter()

	var virtualAdapter = virtual_adapter.NewVirtualAdapter()

	var systemCallCloseFunc = func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-c
		yeelightAdapter.CloseProxy()
		virtualAdapter.CloseProxy()

		os.Exit(0)
	}

	go systemCallCloseFunc()

	for {
		if yeelightAdapter.ProxyRunning() || virtualAdapter.ProxyRunning() {
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Print("main running .....\r\n")
			continue
		} else {
			return
		}
	}

}

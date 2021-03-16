package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	virtualAdapter "virtual_adpater"
)

type AdapterHandler interface {
}

type DeviceHandler interface {
}

type PropertyHandler interface {
}

func main() {

	var YeeAdapter = virtualAdapter.NewYeeAdapter()

	var systemCallCloseFunc = func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-c
		YeeAdapter.CloseProxy()

		os.Exit(0)
	}

	go systemCallCloseFunc()

	for {
		if YeeAdapter.ProxyRunning() || YeeAdapter.ProxyRunning() {
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Print("main running .....\r\n")
			continue
		} else {
			return
		}
	}

}

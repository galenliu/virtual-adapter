package main

import (
	"context"
	"errors"
	addon "github.com/galenliu/gateway-addon-golang"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yeelight-adapter/lib"
)

var (
	On               = "on"
	Brightness       = "brightness"
	Hue              = "hue"
	ColorTemperature = "ct"
	ColorModel       = "ColorMode"
)

func main() {

	var adapter = NewVirtualAdapter("Virtual-adapter", "Virtual-adapter")
	adapter.StartPairing(2000)

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
		if !adapter.IsProxyRunning() {
			time.Sleep(time.Duration(2) * time.Second)
			return
		}
	}

}

type VirtualAdapter struct {
	*addon.AdapterProxy
}

type VirtualDevice struct {
	*addon.DeviceProxy
}

type VirtualProperty struct {
	*addon.PropertyProxy
}

func NewVirtualAdapter(id, packageName string) *VirtualAdapter {
	adapter := &VirtualAdapter{
		addon.NewAdapterProxy(id, packageName),
	}
	adapter.StartPairing(10)
	return adapter
}

func NewVirtualProperty(proxy *addon.PropertyProxy) *VirtualProperty {
	return &VirtualProperty{PropertyProxy: proxy}
}

func NewVirtualDevice(proxy *addon.DeviceProxy) *VirtualDevice {
	return &VirtualDevice{proxy}
}

func (adapter *VirtualAdapter) StartPairing(timeout int) {
	ctx := context.Background()
	ctx1, _ := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	var pairing = func(ctx context.Context) {
		select {
		case <-ctx1.Done():
			adapter.CancelParing()
			return
		default:
			adapter.AdapterProxy.StartPairing(timeout)
			lights := lib.Discover()
			devs := initDevices(adapter, lights)
			for _, d := range devs {
				adapter.HandleDeviceAdded(d)
				return
			}

		}
	}
	go pairing(ctx1)
}



func (prop *VirtualProperty) SetValue(value interface{}) error {
	switch prop.Name {
	case On:
		_ = prop.PropertyProxy.SetValue(value)
		v, ok := value.(bool)
		if !ok {
			return errors.New("value err")

		}
		if v {
			_, _ = prop.lightLib.PowerOn(0)
		} else {
			_, _ = prop.lightLib.PowerOff(0)
		}
	case ColorTemperature:
		_ = prop.PropertyProxy.SetValue(value)
		v, ok := value.(int)
		if !ok {
			return errors.New("value err")

		}
		prop.lightLib.SetTemp(v, 0)
	case Hue:
		prop.PropertyProxy.SetValue(value)
		v, ok := value.(string)
		if !ok {
			return errors.New("value err")

		}
		rgb, err := lib.HTMLToRGB(v)
		if err != nil {
			return err
		}
		prop.lightLib.SetRGB(rgb.R, rgb.G, rgb.B, 0)

	}
	return nil
}

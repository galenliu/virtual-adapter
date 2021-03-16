package virtual_adapter

import (
	"addon"
	"addon/devices"
	"addon/properties"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"virtual_adpater/lib"
)

type YeeAdapter struct {
	locker *sync.Mutex
	*addon.AdapterProxy
}

func NewYeeAdapter() *YeeAdapter {
	adapter := &YeeAdapter{
		AdapterProxy: addon.NewAdapterProxy("yeelight-adapter", "yeelight-adapter", "yeelight-adapter"),
	}
	adapter.OnPairing = adapter.onPairing
	adapter.OnPairing(3000)
	adapter.locker = new(sync.Mutex)
	return adapter
}

func (adapter *YeeAdapter) onPairing(timeout float64) {
	if adapter.IsPairing {
		return
	}
	adapter.IsPairing = true
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	foundDevices := make(chan *addon.DeviceProxy, 20)
	defer func() {
		cancelFunc()
		adapter.IsPairing = false
	}()

	for {
		if !adapter.IsPairing {
			return
		}
		select {
		case device := <-foundDevices:
			adapter.HandleDeviceAdded(device)
		case <-ctx.Done():
			return
		default:
			adapter.Discover(foundDevices)
		}
	}
}

func (adapter *YeeAdapter) Discover(listDevice chan *addon.DeviceProxy) {

	lights := lib.Discover()
	if len(lights) == 0 {
		return
	}
	for _, light := range lights {
		if light.ID == "" {
			break
		}
		fmt.Printf("light id %s，name:%s", light.ID, light.Name)
		title := light.Name
		if title == "" {
			title = "YeeLight" + light.ID[6:]
		}
		lightBulb := NewYeeLight(light.ID, title, &light)

		lightBulb.On.OnValueRemoteUpdate(func(newValue bool) {
			if newValue == true {
				_, err := lightBulb.light.PowerOn(0)
				if err != nil {
					return
				}
			} else {
				_, err1 := lightBulb.light.PowerOff(0)
				if err1 != nil {
					return
				}
			}
			log.Printf("light-prop(%s-%s) changed,new value %v", light.ID, lightBulb.On.Name, newValue)
			return
		})

		for _, prop := range light.Support {
			switch prop {
			case "set_bright":
				lightBulb.Bright = properties.NewBrightnessProperty()
				lightBulb.Bright.Value = lightBulb.light.Bright
				lightBulb.Bright.OnValueRemoteUpdate(func(newValue int) {
					_, _ = lightBulb.light.SetBrightness(newValue, 0)
					log.Printf("light-prop(%s-%s) changed ,new value %v", lightBulb.ID, lightBulb.Bright.Name, newValue)

				})
				lightBulb.AddProperty(lightBulb.Bright.Property)
			case "set_rgb":
				lightBulb.Color = properties.NewColorProperty()
				lightBulb.Color.OnValueRemoteUpdate(func(newValue string) {
					r, g, b, _ := devices.Color16ToRGB(newValue)

					_, _ = lightBulb.light.SetRGB(r, g, b, 0)
				})
				lightBulb.AddProperty(lightBulb.Color.Property)
			}
		}
		//test light
		{
			lightBulb.Toggle()
			time.Sleep(time.Duration(800) * time.Millisecond)
			lightBulb.Toggle()
			_, _ = light.SetBrightness(30, 0)
		}
		listDevice <- lightBulb.DeviceProxy
	}
}

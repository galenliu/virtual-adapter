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
	adapter.locker = new(sync.Mutex)
	return adapter
}

func (adapter *YeeAdapter) onPairing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			adapter.Discover()
		}
	}

}

func (adapter *YeeAdapter) Discover() {
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
		lightBulb := NewYeeLight(light.ID, title)

		lightBulb.On.OnValueRemoteUpdate(func(newValue bool) {
			if newValue == true {
				_, err := light.PowerOn(0)
				if err != nil {
					return
				}
			} else {
				_, err1 := light.PowerOff(0)
				if err1 != nil {
					return
				}
			}
			log.Printf("light-prop(%s-%s) changed,new value %v", lightBulb.ID, lightBulb.On.Name, newValue)
			return
		})

		for _, prop := range light.Support {
			switch prop {
			case "set_bright":
				bright := properties.NewBrightnessProperty()
				bright.Value = light.Bright
				bright.OnValueRemoteUpdate(func(newValue int) {
					_, _ = light.SetBrightness(newValue, 0)
					log.Printf("light-prop(%s-%s) changed ,new value %v", lightBulb.ID, bright.Name, newValue)

				})
				lightBulb.AddProperty(bright.Property)
			case "set_rgb":
				color := properties.NewColorProperty()
				color.OnValueRemoteUpdate(func(newValue string) {
					r, g, b, _ := devices.Color16ToRGB(newValue)

					_, _ = light.SetRGB(r, g, b, 0)
				})
				lightBulb.AddProperty(color.Property)
			}
		}
		//test light
		{
			lightBulb.Toggle()
			time.Sleep(time.Duration(800) * time.Millisecond)
			lightBulb.Toggle()
			_, _ = light.SetBrightness(30, 0)
		}
		go adapter.HandleDeviceAdded(lightBulb.Device)
	}
}

//func (adapter *YeeAdapter) update(light *lib.Light) {
//	device, err := adapter.FindDevice(light.ID)
//	if err != nil {
//		log.Printf(err.Error())
//		return
//	}
//	on, err1 := device.GetProperty(addon.On)
//	if err1 != nil {
//		log.Printf(err1.Error())
//		return
//	}
//	bright, err2 := device.GetProperty(addon.Brightness)
//	if err2 != nil {
//		log.Printf(err2.Error())
//	}
//	var handler = func(message json.Any) {
//		if message.Get("power").ToString() == "on" {
//
//		}
//		if message.Get("power").ToString() == "off" {
//			on.SetValueAndNotify(false)
//		}
//		br, err := strconv.Atoi(message.Get("bright").ToString())
//		if err != nil {
//			log.Printf(err.Error())
//		} else {
//			bright.SetValueAndNotify(br)
//		}
//
//	}
//	go light.Listen(handler)
//}

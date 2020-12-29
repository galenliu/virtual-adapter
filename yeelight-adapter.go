package virtual_adapter

import (
	"context"
	"fmt"
	addon "gateway_addon_golang"
	"log"
	"time"
	"virtual_adpater/lib"
)

type YeeAdapter struct {
	*addon.AdapterProxy
	devices map[string]*addon.Device
}


func NewYeeAdapter() *YeeAdapter {
	adapter := &YeeAdapter{
		AdapterProxy: addon.NewAdapterProxy("yeelight-adapter", "yeelight-adapter"),
		devices:      make(map[string]*addon.Device),
	}
	adapter.OnStartPairing = adapter.StartPairing
	return adapter
}

func (adapter *YeeAdapter) StartPairing(timeout int) {
	fmt.Println("start pairing")

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	var pairing = func(ctx context.Context) {
		select {
		case <-ctx.Done():
			cancelFunc()
			return
		default:
			adapter.Discover()
		}
	}
	go pairing(ctx)
}

func (adapter *YeeAdapter) Discover() {
	lights := lib.Discover()
	for _, light := range lights {
		devInfo := addon.NewDeviceInfo(light.ID, "Yeelight Hue")
		lightBulb := addon.NewLightBulb(devInfo)
		_, _ = light.PowerOff(0)
		time.Sleep(time.Duration(1) * time.Second)
		_, _ = light.PowerOn(0)
		_, _ = light.SetBrightness(100, 2)

		prop := lightBulb.Properties[addon.On]
		if light.Power == "on" {
			prop.Value = true
		} else {
			prop.Value = false
		}
		prop.OnRemoteUpdate = func(oldValue, newValue interface{}) {
			boolValue, _ := newValue.(bool)
			if boolValue == true {
				_, _ = light.PowerOn(0)
			} else {
				_, _ = light.PowerOff(0)
			}
			log.Printf("light-prop(%s-%s) changed,old value: %s ,new value %s", lightBulb.ID, addon.On, oldValue, newValue)
		}

		lightBulb.AddProperties(addon.NewColorModeProperty(addon.ColorModel, light.ColorMode-1))
		for _, prop := range light.Support {
			switch prop {
			case "set_bright":
				p := addon.NewBrightnessProperty(addon.Brightness, 0, 100)
				p.Value = light.Bright
				p.OnRemoteUpdate = func(oldValue, newValue interface{}) {
					intValue, _ := newValue.(int)
					_, _ = light.SetBrightness(intValue, 0)
					log.Printf("light-prop(%s-%s) changed,old value: %s ,new value %s", lightBulb.ID, p.Name, oldValue, newValue)
				}
				lightBulb.AddProperties(p)
			case "set_rgb":
				p := addon.NewColorProperty(addon.ColorTemperature)
				p.Value = light.Hue
				p.OnRemoteUpdate = func(oldValue, newValue interface{}) {
					strValue, _ := newValue.(string)
					r, g, b, err := addon.Color16ToRGB(strValue)
					if err != nil {
						return
					}
					_, _ = light.SetRGB(r, g, b, 0)
					log.Printf("light-prop(%s-%s) changed,old value: %s ,new value %s", lightBulb.ID, p.Name, oldValue, newValue)
				}

				lightBulb.AddProperties(p)
			case "set_ct_abx":
				p := addon.NewColorTemperatureProperty(addon.ColorTemperature, 170000, 680000)
				p.Value = light.ColorTemp
				p.OnRemoteUpdate = func(oldValue, newValue interface{}) {
					intValue, _ := newValue.(int)
					_, _ = light.SetTemp(intValue, 0)
					log.Printf("light-prop(%s-%s) changed,old value: %s ,new value %s", lightBulb.ID, p.Name, oldValue, newValue)
				}
				lightBulb.AddProperties(p)
			default:
				continue
			}
		}
		adapter.devices[lightBulb.ID] = lightBulb.Device
		adapter.HandleDeviceAdded(lightBulb.Device)
	}
}

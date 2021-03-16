package virtual_adapter

import (
	"addon"
	"addon/devices"
	"virtual_adpater/lib"
)

type YeeLight struct {
	*devices.LightBulb
	light *lib.Light
}

func NewYeeLight(deviceId, title string, lib *lib.Light) *YeeLight {
	lightProxy := devices.NewLightBulb(deviceId, title)
	yee := &YeeLight{lightProxy, lib}
	return yee
}

func (y *YeeLight) SetPin(pin addon.PIN) error {
	return nil
}

func (y *YeeLight) TurnOn() {
	y.light.PowerOn(0)
}

func (y *YeeLight) TurnOff() {
	y.light.PowerOff(0)
}

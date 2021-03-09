package virtual_adapter

import (
	"addon/devices"
	"virtual_adpater/lib"
)

type YeeLight struct {
	*devices.LightBulb
	light *lib.Light
}

func NewYeeLight(deviceId, title string) *YeeLight {
	return &YeeLight{
		devices.NewLightBulb(deviceId, title), nil,
	}
}

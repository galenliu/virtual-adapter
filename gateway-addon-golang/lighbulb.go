package gateway_addon_golang

import (
	"strconv"
	"strings"
)

const (
	On               = "on"
	Brightness       = "brightness"
	Hue              = "hue"
	ColorTemperature = "ct"
	ColorModel       = "ColorMode"
)

type LightBulb struct {
	*Device
}

func NewLightBulb(devInfo *DeviceInfo)*LightBulb{
	lightBulb := &LightBulb{}
	lightBulb.Device = NewDevice(devInfo)
	on := NewOnOffProperty(On)
	lightBulb.AddProperties(on)
	lightBulb.AddTypes(Light,OnOffSwitch)
	return lightBulb
}



func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}

package gateway_addon_golang

import _ "github.com/asaskevich/govalidator"

const (
	Alarm                    = "Alarm"
	AirQualitySensor         = "AirQualitySensor"
	BarometricPressureSensor = "BarometricPressureSensor"
	BinarySensor             = "BinarySensor"
	Camera                   = "Camera"
	ColorControl             = "ColorControl"
	ColorSensor              = "ColorSensor"
	DoorSensor               = "DoorSensor"
	EnergyMonitor            = "EnergyMonitor"
	HumiditySensor           = "HumiditySensor"
	LeakSensor               = "LeakSensor"
	Light                    = "Light"
	Lock                     = "Lock"
	MotionSensor             = "MotionSensor"
	MultiLevelSensor         = "MultiLevelSensor"
	MultiLevelSwitch         = "MultiLevelSwitch"
	OnOffSwitch              = "OnOffSwitch"
	SmartPlug                = "SmartPlug"
	SmokeSensor              = "SmokeSensor"
	TemperatureSensor        = "TemperatureSensor"
	Thermostat               = "Thermostat"
	VideoCamera              = "VideoCamera"

	Context = "https://webthings.io/schemas"
)

type DeviceInfo struct {
	AtContext   string   `json:"@context" valid:"url,required"`
	Title       string   `json:"title,required"`
	AtType      []string `json:"@type"`
	Description string   `json:"description,omitempty"`
	ID          string   `json:"id"`
	Version     string   `json:"version,omitempty"`
	Modified    string   `json:"modified,omitempty"`
	Support     string   `json:"support,omitempty"`
	Properties  map[string]*Property
}

func NewDeviceInfo(id string, title string) *DeviceInfo {
	devInfo := &DeviceInfo{
		AtContext:   Context,
		AtType:      nil,
		Title:       title,
		Description: "",
		ID:          id,
		Version:     "",
		Support:     "",
	}
	return devInfo
}

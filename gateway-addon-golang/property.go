package gateway_addon_golang

import "fmt"

type UpdateFunc func(oldValue, newValue interface{})
type PropertyChangedNotification func(property *Property)

const (
	STRING  = "string"
	BOOLEAN = "boolean"
	INTEGER = "integer"
	NUMBER  = "number"

	UnitHectopascal = "hectopascal"
	UnitKelvin      = "kelvin"

	AlarmProperty                    = "AlarmProperty"
	BarometricPressureProperty       = "BarometricPressureProperty"
	BooleanProperty                  = "BooleanProperty"
	BrightnessProperty               = "BrightnessProperty"
	ColorModeProperty                = "ColorModeProperty"
	ColorProperty                    = "ColorProperty"
	ColorTemperatureProperty         = "ColorTemperatureProperty"
	ConcentrationProperty            = "ConcentrationProperty"
	CurrentProperty                  = "CurrentProperty"
	DensityProperty                  = "DensityProperty"
	FrequencyProperty                = "FrequencyProperty"
	HeatingCoolingProperty           = "HeatingCoolingProperty"
	HumidityProperty                 = "HumidityProperty"
	ImageProperty                    = "ImageProperty"
	InstantaneousPowerFactorProperty = "InstantaneousPowerFactorProperty"
	InstantaneousPowerProperty       = "InstantaneousPowerProperty"
	LeakProperty                     = "LeakProperty"
	LevelProperty                    = "LevelProperty"
	LockedProperty                   = "LockedProperty"
	MotionProperty                   = "MotionProperty"
	OnOffProperty                    = "OnOffProperty"
	OpenProperty                     = "OpenProperty"
	PushedProperty                   = "PushedProperty"
	SmokeProperty                    = "SmokeProperty"
	TargetTemperatureProperty        = "TargetTemperatureProperty"
	TemperatureProperty              = "TemperatureProperty"
	ThermostatModeProperty           = "ThermostatModeProperty"
	VideoProperty                    = "VideoProperty"
	VoltageProperty                  = "VoltageProperty"
)

type Property struct {
	AtType      string `json:"@type"` //引用的类型
	Type        string `json:"type"`  //数据的格式
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`

	Unit     string `json:"unit,omitempty"` //属性的单位
	ReadOnly bool   `json:"read_only"`
	Visible  bool   `json:"visible"`

	Minimum interface{} `json:"minimum,omitempty,string"`
	Maximum interface{} `json:"maximum,omitempty,string"`
	Value   interface{}
	Enum    []interface{} `json:"-"`

	//func
	OnRemoteUpdate                UpdateFunc                  `json:"-"`
	OnPropertyChangedNotification PropertyChangedNotification `json:"-"`
}

//set and notify
func (prop *Property) SetValueAndNotify(newValue interface{}) {
	prop.setCachedValue(newValue)
	prop.OnPropertyChangedNotification(prop)
}


//set prop value
func (prop *Property) SetValue(newValue interface{}) {
	prop.setCachedValue(newValue)
}

func (prop *Property) setCachedValue(newValue interface{}) error {
	if prop.ReadOnly{
		return fmt.Errorf("property(%s) read only",prop.Name)
	}
	oldValue := prop.Value
	if oldValue == newValue {
		return fmt.Errorf("new value equal old value")
	}
	if prop.Type == NUMBER || prop.Type == INTEGER {
		intValue, ok := newValue.(int)
		if !ok {
			return fmt.Errorf("value need a integer or number")
		}
		if prop.Maximum != nil {
			maxInt, _ := prop.Maximum.(int)
			minInt, _ := prop.Maximum.(int)
			if intValue > maxInt || intValue < minInt {
				return fmt.Errorf("value range err")
			}
		}
		prop.Value = newValue
	}
	if prop.Type == STRING {
		prop.Value = newValue
	}

	return nil
}

func NewStringProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        STRING,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
		Enum:        nil,
	}
	return p
}

func NewBooleanProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        BOOLEAN,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
	}

	return p
}

func NewNumberProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        NUMBER,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
	}
	return p
}

func NewIntegerProperty(name string, atType string) *Property {
	p := &Property{
		AtType:      atType,
		Type:        INTEGER,
		Title:       "",
		Description: "",
		Name:        name,
		Unit:        "",
		ReadOnly:    false,
		Visible:     true,
		Minimum:     nil,
		Maximum:     nil,
		Value:       nil,
	}
	return p
}

func NewColorTemperatureProperty(name string, min int, max int) *Property {
	p := NewIntegerProperty(name, "ColorTemperatureProperty")
	p.Minimum = min
	p.Maximum = max
	p.Unit = UnitKelvin
	return p
}

func NewColorProperty(name string) *Property {
	p := NewStringProperty(name, "ColorProperty")

	return p
}

func NewOnOffProperty(name string) *Property {
	p := NewBooleanProperty(name, "OnOffProperty")

	p.Value = true
	return p
}

func NewBrightnessProperty(name string, min int, max int) *Property {
	p := NewIntegerProperty(name, "BrightnessProperty")
	p.Minimum = min
	p.Maximum = max
	p.Value = min

	return p
}

func NewColorModeProperty(name string, mode int) *Property {
	var enum = []interface{}{"color", "temperature", "HSV"}
	p := NewIntegerProperty(name, "ColorModeProperty")
	p.ReadOnly = true
	p.Enum = enum
	p.Value = p.Enum[mode]
	return p
}

package gateway_addon_golang

type Device struct {
	*DeviceInfo
	adaper     *AdapterProxy        `json:"-"`
	pin        interface{}          `json:"-"`
	username   string               `json:"—"`
	password   string               `json:"-"`
	Properties map[string]*Property `json:"properties"`
	Actions    map[string]*Action   `json:"-"`
}

func NewDevice(deviceInfo *DeviceInfo) *Device {
	device := &Device{}
	device.Properties = make(map[string]*Property)
	device.DeviceInfo = deviceInfo
	return device
}

func (device *Device) AddProperties(props ...*Property) {
	for _, p := range props {
		device.Properties[p.Name] = p
		p.OnPropertyChangedNotification = device.adaper.sendPropertyChangedNotification
	}
}

func (device *Device) AddActions(actions ...*Action) {
	for _, a := range actions {
		a.ActionFunc = device.adaper.sendActionNotification
		device.Actions[a.Name] = a
	}
}

func (device *Device) setPin(pin interface{}) error {
	device.pin = pin
	return nil
}

func (device *Device) setCredentials(username, password string) error {
	device.username = username
	device.password = password
	return nil
}

func (device *Device) AddTypes(types ...string) {
	for _, t := range types {
		device.AtType = append(device.AtType, t)
	}
}

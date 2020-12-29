package gateway_addon_golang

import (
	"errors"
	"fmt"
	"time"
)

type OnPairingFunc func(timeout int)

type AdapterProxy struct {
	id             string
	packageName    string
	managerProxy   *AddonManagerProxy
	Devices        map[string]*Device
	Name           string
	OnStartPairing OnPairingFunc
	pairing        bool
}

func NewAdapterProxy(_id string, _packageName string) (adp *AdapterProxy) {
	adp = &AdapterProxy{id: _id, packageName: _packageName}
	adp.Devices = make(map[string]*Device)
	adp.managerProxy = NewAddonManagerProxy(_packageName)
	adp.pairing = false
	time.Sleep(time.Duration(1) * time.Second)
	adp.handleAdapterAdded(adp)
	return adp
}

func (adapter *AdapterProxy) HandleDeviceAdded(device *Device) {
	device.adaper = adapter
	adapter.managerProxy.handleDeviceAdded(device, adapter.id)
}

func (adapter *AdapterProxy) handleAdapterAdded(ad *AdapterProxy) {
	adapter.managerProxy.addAdapter(ad)
}

//配对方案，主要用于发现设备，初始化设备。
func (adapter *AdapterProxy) startPairing(timeout int) {
	fmt.Println("start pairing ....")
	adapter.OnStartPairing(timeout)
}

func (adapter *AdapterProxy) CancelParing() {

}

func (adapter *AdapterProxy) HandleDeviceSaved(devId string, dev interface{}) {
	fmt.Print("on device saved on the gateway")
}

func (adapter *AdapterProxy) Run(retryIntervalSeconds int) {
	adapter.managerProxy.run()
}

func (adapter *AdapterProxy) Unload() {
	fmt.Printf("adapter unload, adapterId:%v", adapter.id)
}

func (adapter *AdapterProxy) removeDevice(deviceId string) {
	dev := adapter.Devices[deviceId]
	if dev != nil {
		adapter.handleDeviceRemoved(dev.ID)
	}
}
func (adapter *AdapterProxy) cancelRemoveDevice(deviceId string) {
	dev := adapter.Devices[deviceId]
	if dev != nil {
		fmt.Print("cancel remove device, Id:", dev.ID)
	}
}

func (adapter *AdapterProxy) handleDeviceRemoved(devId string) {
	delete(adapter.Devices, devId)
	adapter.managerProxy.handleDeviceRemoved(adapter.id, devId)
}

func (adapter *AdapterProxy) getID() string {
	return adapter.id
}

func (adapter *AdapterProxy) getPackageName() string {
	return adapter.packageName
}

func (adapter *AdapterProxy) getName() string {
	if adapter.Name != "" {
		return adapter.Name
	} else {
		return adapter.id
	}

}
func (adapter *AdapterProxy) SetPin(devId, pin string) error {
	dev := adapter.getDevice(devId)
	if dev == nil {
		return errors.New("device no find")
	}
	return nil
}

func (adapter *AdapterProxy) SetCredentials(devId, username, password string) error {
	dev := adapter.getDevice(devId)
	if dev == nil {
		return errors.New("device no find")
	}
	return nil

}

func (adapter *AdapterProxy) getDevice(id string) *Device {
	return adapter.Devices[id]
}

func (adapter *AdapterProxy) sendPropertyChangedNotification(p *Property) {
	adapter.managerProxy.sendPropertyChangedNotification(p)
}

func (adapter *AdapterProxy) sendActionNotification(a *Action) {

}

func (adapter *AdapterProxy) CloseProxy() {
	adapter.managerProxy.close()
	fmt.Print("do some thing while adapter close")
}

func (adapter *AdapterProxy) ProxyRunning() bool {
	return adapter.managerProxy.running
}

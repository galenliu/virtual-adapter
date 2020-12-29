package gateway_addon_golang

import (
	"fmt"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"log"
	"sync"
	"time"
)

type AddonManagerProxy struct {
	mu        *sync.Mutex
	Adapters  map[string]*AdapterProxy
	pluginId  string
	ipcClient *IpcClient
	verbose   bool
	running   bool
}

var once sync.Once
var addonManager *AddonManagerProxy

func NewAddonManagerProxy(packetName string) *AddonManagerProxy {
	once.Do(
		func() {
			addonManager = &AddonManagerProxy{}
			addonManager.Adapters = make(map[string]*AdapterProxy)
			addonManager.pluginId = packetName
			addonManager.ipcClient = NewClient(packetName, addonManager.OnMessage)
			addonManager.running = true
			addonManager.mu = new(sync.Mutex)
		},
	)
	return addonManager
}

func (proxy *AddonManagerProxy) addAdapter(adapter *AdapterProxy) {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	proxy.Adapters[adapter.id] = adapter
	proxy.send(AdapterAddedNotification,
		struct {
			Name      string `json:"name"`
			PluginId  string `json:"pluginId"`
			AdapterId string `json:"adapterId"`
		}{
			AdapterId: adapter.getID(),
			PluginId:  proxy.pluginId,
			Name:      adapter.getName(),
		},
	)
}

func (proxy *AddonManagerProxy) OnMessage(data []byte, conn *websocket.Conn) {

	var messageType = json.Get(data, "messageType").ToInt()

	switch messageType {
	//卸载plugin
	case PluginUnloadRequest:
		proxy.send(PluginUnloadResponse, struct {
			PluginId string
		}{PluginId: proxy.pluginId})
		proxy.running = false
		var closeFun = func() {
			time.AfterFunc(500*time.Millisecond, func() { proxy.close() })
		}
		go closeFun()
		return
	}

	var adapterId = json.Get(data, "data", "adapterId").ToString()
	adapter, ok := proxy.Adapters[adapterId]
	if !ok {
		log.Fatal("addon manager can not find adapter : %s", adapterId)
		return
	}

	switch messageType {

	//adapter pairing command
	case AdapterStartPairingCommand:
		timeout := json.Get(data, "data", "timeout").ToInt()
		go adapter.startPairing(timeout)
		return

	case AdapterCancelPairingCommand:
		go adapter.CancelParing()
		return

		//adapter unload request
	case AdapterUnloadRequest:
		adapter.Unload()
		unloadFunc := func(proxy *AddonManagerProxy, adapter *AdapterProxy) {
			proxy.send(AdapterUnloadResponse, struct {
				AdapterId string `json:"adapterId"`
			}{AdapterId: adapter.id})
		}
		go unloadFunc(proxy, adapter)
		delete(proxy.Adapters, adapter.getID())
		return
	}

	var deviceId = json.Get(data, "data", "device_id").ToString()
	device, ok := adapter.Devices[deviceId]
	if !ok {
		log.Fatal("addon manager can not find device :%s", deviceId)
		return
	}
	switch messageType {
	case AdapterCancelRemoveDeviceCommand:
		adapter := proxy.Adapters[adapterId]
		adapter.cancelRemoveDevice(deviceId)

	case DeviceSavedNotification:
		adapter := proxy.Adapters[adapterId]
		go adapter.HandleDeviceSaved(deviceId, device)
		return

		//adapter remove device request
	case AdapterRemoveDeviceRequest:
		go adapter.removeDevice(deviceId)

		//device set property command
	case DeviceSetPropertyCommand:
		propName := json.Get(data, "data", "propertyName").ToString()
		newValue := json.Get(data, "data", "propertyValue").GetInterface()
		prop, ok := device.Properties[propName]
		if !ok {
			log.Fatal("DeviceSetPropertyCommand: device(%s) have not propName(%s): ", deviceId, propName)
			return
		}
		propChanged := func(oldValue, newValue interface{}) {
			err := prop.setCachedValue(newValue)
			if err != nil {
				return
			}
			prop.OnRemoteUpdate(oldValue, newValue)
		}
		go propChanged(prop.Value, newValue)

	case DeviceSetPinRequest:
		pin := json.Get(data, "data", "pin").GetInterface()
		if pin == nil {
			log.Fatal("DeviceSetPinRequest: not find pin form message")
			return
		}
		messageId := json.Get(data, "data", "message_id").ToInt()
		if messageId == 0 {
			log.Fatal("DeviceSetPinRequest:  non  messageId")
		}
		handleFunc := func() {
			err := device.setPin(pin)
			if err == nil {
				proxy.send(DeviceSetPinResponse, struct {
					PluginId  string
					AdapterId string
					MessageId int
					DeviceId  string
					Device    *Device
					Success   bool
				}{
					PluginId:  proxy.pluginId,
					AdapterId: adapterId,
					MessageId: messageId,
					DeviceId:  deviceId,
					Device:    device,
					Success:   true,
				})

			} else {
				proxy.send(DeviceSetPinResponse, struct {
					PluginId  string
					AdapterId string
					MessageId int
					DeviceId  string
					Device    *Device
					Success   bool
				}{
					PluginId:  proxy.pluginId,
					AdapterId: adapterId,
					MessageId: messageId,
					DeviceId:  deviceId,
					Device:    device,
					Success:   false,
				})
			}
		}
		go handleFunc()

	case DeviceSetCredentialsRequest:
		messageId := json.Get(data, "data", "message_id").ToInt()
		username := json.Get(data, "data", "username").ToString()
		password := json.Get(data, "data", "password").ToString()

		handleFunc := func() {
			err := device.setCredentials(username, password)
			if err != nil {
				fmt.Printf(err.Error())
				proxy.send(DeviceSetCredentialsResponse, struct {
					PluginId  string
					AdapterId string
					MessageId int
					DeviceId  string
					Device    *Device
					Success   bool
				}{
					PluginId:  proxy.pluginId,
					AdapterId: adapterId,
					MessageId: messageId,
					DeviceId:  deviceId,
					Device:    device,
					Success:   false,
				})
				return
			}
			proxy.send(DeviceSetCredentialsResponse, struct {
				PluginId  string
				AdapterId string
				MessageId int
				DeviceId  string
				Device    *Device
				Success   bool
			}{
				PluginId:  proxy.pluginId,
				AdapterId: adapterId,
				MessageId: messageId,
				DeviceId:  deviceId,
				Device:    device,
				Success:   true,
			})
		}
		go handleFunc()
	}

}

func (proxy *AddonManagerProxy) sendPropertyChangedNotification(p *Property) {
	data := struct {
		AdapterId string    `json:"adapterId"`
		DeviceId  string    `json:"device_id"`
		Property  *Property `json:"property"`
	}{
		AdapterId: proxy.pluginId,
		DeviceId:  proxy.pluginId,
		Property:  p,
	}
	proxy.send(DevicePropertyChangedNotification, data)
}

func (proxy *AddonManagerProxy) run() {
	proxy.ipcClient.runLoop()
}

func (proxy *AddonManagerProxy) handleDeviceAdded(dev *Device, adapterId string) {
	message := struct {
		PluginId  string      `json:"pluginId"`
		AdapterId string      `json:"adapterId"`
		Device    interface{} `json:"device"`
	}{
		PluginId:  proxy.pluginId,
		AdapterId: adapterId,
		Device:    dev,
	}
	proxy.send(DeviceAddedNotification, message)

}

func (proxy *AddonManagerProxy) handleDeviceRemoved(adapterId, devId string) {
	if proxy.verbose {
		fmt.Printf("addon manager handle device added, deviceId:%v\n", devId)
	}
	message := struct {
		PluginId  string `json:"pluginId"`
		AdapterId string `json:"adapterId"`
	}{
		PluginId:  proxy.pluginId,
		AdapterId: adapterId,
	}
	proxy.send(AdapterRemoveDeviceResponse, message)
}

func (proxy *AddonManagerProxy) send(messageType int, data interface{}) {

	var message = struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{MessageType: messageType, Data: data}

	d, er := json.MarshalIndent(message, "", "")
	if er != nil {
		log.Fatal(er)
		return
	}
	proxy.ipcClient.sendMessage(d)
}

func (proxy *AddonManagerProxy) close() {
	proxy.ipcClient.close()
	proxy.running = false
}

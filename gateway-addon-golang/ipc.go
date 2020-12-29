package gateway_addon_golang

import (
	"fmt"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"net/url"
	"sync"

	"log"
)

const (
	Disconnect = "Disconnect"
	Connected  = "Connected"
	Registered = "Registered"
)

type UserProfile struct {
	BaseDir        string `validate:"required" json:"base_dir"`
	DataDir        string `validate:"required" json:"data_dir"`
	AddonsDir      string `validate:"required" json:"addons_dir"`
	ConfigDir      string `validate:"required" json:"config_dir"`
	UploadDir      string `validate:"required" json:"upload_dir"`
	MediaDir       string `validate:"required" json:"media_dir"`
	LogDir         string `validate:"required" json:"log_dir"`
	GatewayVersion string
}

type Preferences struct {
	Language string `validate:"required" json:"language"`
	Units    Units  `validate:"required" json:"units"`
}

type Units struct {
	Temperature string `validate:"required" json:"temperature"`
}

const PORT = "9600"
const RetrySecond = 5 //retry dail before second

type OnMessage func(data []byte, conn *websocket.Conn)

//为Plugin提供和gateway Server进行消息的通信

type IpcClient struct {
	ws *websocket.Conn

	url         string
	preferences Preferences
	userProfile UserProfile

	writeCh   chan []byte
	readCh    chan []byte
	closeChan chan interface{}
	reConnect chan interface{}

	gatewayVersion string

	onMessage OnMessage
	mu        *sync.Mutex

	status   string
	pluginID string
	origin   string
	verbose  bool
}

//新建一个Client，注册消息Handler
func NewClient(PluginId string, handler OnMessage) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + PORT, Path: "/"}
	client := &IpcClient{}
	client.pluginID = PluginId
	client.url = u.String()
	client.status = Disconnect
	client.mu = new(sync.Mutex)

	client.closeChan = make(chan interface{})
	client.reConnect = make(chan interface{})

	client.readCh = make(chan []byte)
	client.writeCh = make(chan []byte)

	client.onMessage = handler

	//读协程
	go client.readLoop()

	//写协程
	go client.writeLoop()

	go client.runLoop()

	return client
}

func (client *IpcClient) onData(data []byte) {

	if json.Get(data, "messageType").ToInt() == PluginRegisterResponse {
		client.preferences.Language = json.Get(data, "data", "preferences", "language").ToString()
		client.preferences.Units.Temperature = json.Get(data, "data", "preferences", "units", "temperature").ToString()
		client.userProfile.AddonsDir = json.Get(data, "data", "user_profile", "addons_dir").ToString()
		client.userProfile.BaseDir = json.Get(data, "data", "user_profile", "base_dir").ToString()
		client.userProfile.ConfigDir = json.Get(data, "data", "user_profile", "config_dir").ToString()
		client.userProfile.DataDir = json.Get(data, "data", "user_profile", "data_dir").ToString()
		client.userProfile.GatewayVersion = json.Get(data, "data", "user_profile", "gateway_version").ToString()
		client.userProfile.LogDir = json.Get(data, "data", "user_profile", "log_dir").ToString()
		client.userProfile.MediaDir = json.Get(data, "data", "user_profile", "media_dir").ToString()
		client.userProfile.UploadDir = json.Get(data, "data", "user_profile", "upload_dir").ToString()
		client.status = Registered
	} else {
		client.onMessage(data, client.ws)
	}
}

//发送Message Struct
func (client *IpcClient) sendMessage(data []byte) {
	log.Printf("send message:  %s\r\n", string(data))
	if client.ws != nil {
		client.writeCh <- data
		return
	}
}

//循环往readCh中读取 Message
func (client *IpcClient) readLoop() {

	for {
		if client.ws != nil {
			_, data, err := client.ws.ReadMessage()
			if err == nil {
				log.Printf("loop read data : %s", string(data))
				client.onData(data)
			} else {
				log.Printf("read loop err : %s", err.Error())
			}

		}
	}
}

//循环发送writeChan中的Message
func (client *IpcClient) writeLoop() {
	defer client.close()
	for {
		select {
		case msg := <-client.writeCh:
			fmt.Printf("witre data: %s \r\n", string(msg))
			if client.ws != nil && client.status != Disconnect {
				err := client.ws.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					fmt.Printf("write loop err =%v", err)
					client.status = Disconnect
				}
			}
		case <-client.closeChan:
			return
		}

	}
}

//client 失去连接后，重新连接
func (client *IpcClient) runLoop() {
	for {
		if client.status == Disconnect {
			err := client.dial()
			if err != nil {
				client.status = Disconnect
				client.ws = nil
				fmt.Printf("pluginID:%v, err:%v ,retry after %v second \r\n", client.pluginID, err, 5)
				continue
			}

		}
		client.status = Connected
	}
}

func (client *IpcClient) close() {
	if client.ws != nil {
		err := client.ws.Close()
		if err != nil {
			fmt.Println("client close-----")
		}
	}
	client.closeChan <- ""
}

func (client *IpcClient) dial() error {

	var err error

	client.ws, _, err = websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		fmt.Printf("dial err: %s \r\n", err.Error())
		return err
	}
	message := struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		MessageType: PluginRegisterRequest,
		Data: struct {
			PluginID string `json:"pluginId"`
		}{PluginID: client.pluginID},
	}

	d, er := json.Marshal(message)
	if er != nil {
		return er
	}
	client.sendMessage(d)
	return nil
}

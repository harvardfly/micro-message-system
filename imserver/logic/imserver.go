package logic

/*
im 相当于是logic层 im
流程：
1. 首先用户登录进websocket，发送消息告知已连接websocket token与conn绑定
2. 通过网关发送消息的时候去向kafka发送消息，订阅者在不断地监听，当有消息过来的时候就通过SendMessage发送消息
3. kafka订阅的消息发送到websocket显示
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"

	"micro-message-system/common/baseerror"
	"micro-message-system/imserver/util"
)

type (
	ImServer struct {
		kafkaBroker *util.KafkaBroker
		clients     map[string]*websocket.Conn
		Address     string
		lock        sync.Mutex
		upgraer     *websocket.Upgrader
	}
	SendMsgRequest struct {
		FromToken     string `json:"fromToken"`
		ToToken       string `json:"toToken"`
		Body          string `json:"body"`
		TimeStamp     int64  `json:"timeStamp"`
		RemoteAddress string `json:"remoteAddress"`
	}
	LoginRequest struct {
		Token string `json:"token"`
	}
	SendMsgResponse struct {
		FromToken     string `json:"fromToken"`     // 消息来自谁
		Body          string `json:"body"`          // 消息内容
		RemoteAddress string `json:"remoteAddress"` // 消息远程地址
	}
	ImServerOptions func(im *ImServer)
)

var (
	DefaultAddress = ":7272" // websocket 默认地址
	UserNoLoginErr = baseerror.NewBaseError("此用户没有登录！")
	SendMessageErr = baseerror.NewBaseError("发送消息失败！")
)

func NewImServer(kafkaBroker *util.KafkaBroker, opts ImServerOptions) (*ImServer, error) {
	// 初始化
	if err := broker.Init(); err != nil {
		return nil, err
	}
	if err := broker.Connect(); err != nil {
		return nil, err
	}
	imServer := &ImServer{
		kafkaBroker: kafkaBroker,
		clients:     make(map[string]*websocket.Conn, 0), // 用户：多个websocket链接
		// 初始化websocket的读取大小和写入大小
		upgraer: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	if opts != nil {
		opts(imServer)
	}
	if imServer.Address == "" {
		imServer.Address = DefaultAddress
	}
	return imServer, nil
}

// kafka订阅的消息 发送到websocket
func (l *ImServer) SendMsg(r *SendMsgRequest) (*SendMsgResponse, error) {
	// websocket发送消息
	l.lock.Lock()
	defer l.lock.Unlock()
	log.Printf("send SendMsgRequest  %+v", r)
	conn := l.clients[r.ToToken] //c获取用户的websocket链接
	if conn == nil {
		return nil, UserNoLoginErr
	}
	r.TimeStamp = time.Now().Unix()
	r.RemoteAddress = conn.RemoteAddr().String()
	bodyMsg, err := json.Marshal(r)
	if err != nil {
		return nil, SendMessageErr
	}
	// 向websocket发送消息
	if err := conn.WriteMessage(websocket.TextMessage, bodyMsg); err != nil {
		log.Printf("send message err %v", err)
		l.clients[r.ToToken] = nil
		//log.Println(conn.Close())
		return nil, err
	}
	log.Printf("send message succes  %v", r.Body)
	return &SendMsgResponse{}, nil
}

// 订阅消息
func (l *ImServer) Subscribe() {
	fmt.Println("开始Subscribe")
	l.kafkaBroker.Subscribe(func(msg []byte) error {
		r := new(SendMsgRequest)
		fmt.Printf("r:s%", r)
		if err := json.Unmarshal(msg, r); err != nil {
			log.Printf("[Unmarshal msg err] : %+v", err)
			return err
		}
		if _, err := l.SendMsg(r); err != nil {
			log.Printf("[SendMsg err] : %+v", err)
			return err
		}
		log.Printf("has Subscribe msg %+v", string(msg))
		return nil
	})
}

func (l *ImServer) Run() {
	log.Printf("websocket has listens at %s", l.Address)
	http.HandleFunc("/ws", l.login)
	log.Fatal(http.ListenAndServe(l.Address, nil))
}

// 用户和websocket关联
func (l *ImServer) login(w http.ResponseWriter, r *http.Request) {
	// 连接websocket
	conn, err := l.upgraer.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 用户登录后会发一个消息 用户已连接websocket
	msgType, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("read login message err %+v", err)
		return
	}
	log.Printf("用户已连接websocket:s%", message)
	log.Printf("msgType:%s", msgType)
	// 只发送文本类型的消息
	if msgType != websocket.TextMessage {
		log.Printf("read login msgType err %+v", err)
		return
	}
	fmt.Println(string(message))
	loginMsgRequest := new(LoginRequest)

	// 从连接读到的消息中获取登录的token
	if err := json.Unmarshal(message, loginMsgRequest); err != nil {
		log.Printf("json.Unmarshal msg err %+v", err)
		return
	}
	l.clients[loginMsgRequest.Token] = conn
	fmt.Println(l.clients)
	return
}

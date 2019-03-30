package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type server struct {
	ip string
	port int
}

/**
	code :1 // 加入测试数据
	code :2 // 登录

	code :1001 // 返回队伍状态变化
	code :1002 // 返回可以登录状态
	code :1003 // 返回客户端名
 */
type message struct {
	Code   int    	`json:"code"`
	Data  string    `json:"data"`
	Param string 	`json:"param"`
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端读写消息
type wsMessage struct {
	messageType int
	data []byte
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan chan *wsMessage	// 读队列
	outChan chan *wsMessage // 写队列

	mutex sync.Mutex	// 避免重复关闭管道
	isClosed bool
	closeChan chan byte  // 关闭通知
}

func (wsConn *wsConnection)wsReadLoop() {
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &wsMessage{
			msgType,
			data,
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- req:
		case <- wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection)wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <- wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				goto error
			}
		case <- wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection)procLoop() {
	// 启动一个gouroutine发送心跳
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if wsConn.isClosed {
				os.Exit(1)
				break
			}
			//handlerLen()
		}
	}()

	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			fmt.Println("read fail")
			break
		}

		//fmt.Println(string(msg.data)) // 接收消息

		// decode
		decodeHandler(msg)

		//err = wsConn.wsWrite(msg.messageType, msg.data)
		//if err != nil {
		//	fmt.Println("write fail")
		//	break
		//}
	}
}
var wsConn *wsConnection
func wsHandler(resp http.ResponseWriter, req *http.Request) {
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	wsConn = &wsConnection{
		wsSocket: wsSocket,
		inChan: make(chan *wsMessage, 1000),
		outChan: make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed: false,
	}

	// 处理器
	go wsConn.procLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}

func (wsConn *wsConnection)wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data,}:
	case <- wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

func (wsConn *wsConnection)wsRead() (*wsMessage, error) {
	select {
	case msg := <- wsConn.inChan:
		return msg, nil
	case <- wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

func (wsConn *wsConnection)wsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

func NewServer(ip string,port int) *server {
	return &server{
		ip:     ip,
		port:	port,
	}
}

func (s *server)Start(){
	address := s.ip +":" + strconv.Itoa(s.port)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("wsserver start..",address )
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

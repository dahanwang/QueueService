package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

var Myuid string
type ChannelInfo struct {
	Id string
	Name string
	LoginTime int64
}

/**
decode
 */
func decodeHandler(netMessage *wsMessage)  {
	var jsonMessage message
	if err := json.Unmarshal([]byte(netMessage.data), &jsonMessage); err == nil {
		if jsonMessage.Code == 1 {
			test()
		}else if jsonMessage.Code == 2 {
			login()
		}else {
			fmt.Println("消息码错误!!!")
		}
	} else {
		fmt.Println("decodeHandler fail",err)
	}
}

/**
encode
 */
func encodeHandler(dataMessage *message) {
	data,_ := json.Marshal(dataMessage)
	var err = wsConn.wsWrite(websocket.TextMessage, data)
	if err != nil {
		fmt.Println("encodeHandler fail", err)
	}
}

func test()  {
	fmt.Println("test")

	for i:=0;i<10;i++ {
		pushQueue("")
	}
}

func login() {
	myUUID,_ := uuid.NewV4()
	Myuid = myUUID.String()
	pushQueue(Myuid)

	if queueNodeMap == nil {
		queueNodeMap = make(map[string]int)
	}
	queueNodeMap[Myuid] = ChannelQueue.Len()

	// 发送队伍状态
	HandlerLen()

	// 发送客户端名
	nameMessage := &message{
		Code : 1003,
		Data: Myuid,
	}
	encodeHandler(nameMessage)
	fmt.Println("login")
}

func HandlerLen() {
	totalLen := ChannelQueue.Len()

	if queueNodeMap == nil {
		queueNodeMap = make(map[string]int)
	}

	// 判断当前客户端在队伍中状态
	currentLen := queueNodeMap[Myuid]
	if _, ok := queueNodeMap[Myuid]; ok {
		currentLen := currentLen - 1
		queueNodeMap[Myuid] = currentLen
	}

	// 发送队伍状态
	lenMessage := &message{
		Code : 1001,
		Data:strconv.Itoa(totalLen),
		Param:strconv.Itoa(currentLen),
	}
	encodeHandler(lenMessage)

	// 如果排队到了，发送可以登录消息
	if _, ok := queueNodeMap[Myuid]; ok {
		if currentLen <= 0 {
			loginMessage := &message{
				Code : 1002,
				Data : Myuid,
			}
			encodeHandler(loginMessage)
		}
	}
}

func pushQueue(uid string) {
	if uid == "" {
		tempUid,_ := uuid.NewV4()
		uid = tempUid.String()
	}
	queueNode := &ChannelInfo{
		Id : uid,
		Name : GetRandomString(),
		LoginTime: time.Now().Unix(),
	}

	ChannelQueue.Push(queueNode)
}

func  GetRandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

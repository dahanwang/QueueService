package main

import (
	"./service"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type ChannelInfoFIFO struct {
	InfoMap   map[string]bool
	InfoFIFO  []service.ChannelInfo
	RwMutex   sync.RWMutex
	MaxLength int
}

var (
	logProducer = log.WithFields(log.Fields{
		"Modul": "ChannelInfoProducer",
	})
	logConsumer = log.WithFields(log.Fields{
		"Modul": "ChannelInfoConsumer",
	})

	channelInfoFIFO ChannelInfoFIFO
)

func initLog() {
	//log format
	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02T15:04:05.000000"
	log.SetFormatter(formatter)

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

func init() {
	initLog()
	channelInfoFIFO.InfoMap = make(map[string]bool)

}

func ChannelInfoConsumer(queue *service.Queue) {
	logProducer.Debug(fmt.Sprintf("consume start"))
	for {
		//get head elem
		qNode := queue.Pull()
		if qNode == nil {
			//logConsumer.Warn("empty queue!")
			time.Sleep(time.Second)
			continue
		}

		_, ok := qNode.(*service.ChannelInfo)
		if !ok {
			//logConsumer.Warn(fmt.Sprintf("qNode Type:%T, qNode:%+v\n", qNode, qNode))
			time.Sleep(time.Second)
			continue
		}

		//logConsumer.Debug(fmt.Sprintf("consume a elem, get_time:%+v, info:%+v\n",
			//cinfo.LoginTime, cinfo))

		// 发送队列
		service.HandlerLen()

		randomSecond := rand.Intn(4) + 1
		time.Sleep(time.Duration(randomSecond) * time.Second )
	}

}

func main() {
	server := service.NewServer("127.0.0.1",8001)
	go server.Start()

	go ChannelInfoConsumer(service.ChannelQueue)
	for {
		time.Sleep(time.Minute * 100)
	}
}

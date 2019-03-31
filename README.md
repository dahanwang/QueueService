# QueueService

**安装依赖库**

```
go get github.com/sirupsen/logrus
go get github.com/gorilla/websocket
go get github.com/satori/go.uuid
```

**启动服务器**

```
go run main.go
```

**启动客户端**

```
浏览器打开client.html
```

**测试方法**
1. 点击Open与服务器连接
1. 点击Test加入机器人(随时可点击多次)
1. 点击Login加入客户端

**QueueService原理**
1. 客户端与服务端websocket连接
2. 为测试排队效果加入了机器人
3. 服务端queue每1-5秒处理一个排队数据
4. queue记录当前客户端所在位置，并实时广播给客户端当前queue状态
5. 排到当前客户端会发送Login确认消息。服务端关闭。

**运行效果**

![image](https://github.com/dahanwang/QueueService/blob/master/20190330154343.png)

**TODOLIST**
1. 由于时间比较紧，这边的项目工作比较多。未进行压测。
1. 有部分值得优化的地方
1. 目前QueueService还属于单点，以下是根据SOA思想对架构的初步设计




## 仅用于FunPlus线下编程作业

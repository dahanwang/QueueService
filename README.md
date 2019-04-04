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

**QueueService架构设计**

目前QueueService还属于单点，以下是根据SOA思想对架构的初步设计

![image](https://github.com/dahanwang/QueueService/blob/master/15540223181.png)

1. LoginServer(登录服务器)向zookeeper注册
1. QueueServer(队列服务器)向zookeeper订阅
1. 客户端请求QueueServer，QueueServer根据LoginServer综合健康指标(注册人数，在线人数，排队状态等)获取具体LoginServer
1. QueueServer将LoginServer状态实时通知给客户端
1. 达到登录条件后，客户端连接LoginServer登录

由于这边的项目工作比较多，业余时间完成的本次作业。还有些功能没有完善以及可以优化的地方

#### *TODOLIST*
1. 未进行登录认证
1. 未进行压测
1. 未实现多客户端连接


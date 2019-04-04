# QueueService

1. 根据要求去掉Queue类，直接使用channel
1. websocket参数修改为1

## 运行效果
![image](https://github.com/dahanwang/QueueService/blob/queue1/20190401211402.png)

#### 疑问(只是想到下面的疑问，可能存在理解错误）
1. channel是线程安全的，加锁的确没有必要
1. 修改websocket channel 缓冲长度为1，是否会降低吞吐量？



第一版 [QueueService](https://github.com/dahanwang/QueueService) 切换master分支


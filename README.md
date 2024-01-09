#IM-CLIENT

##目的
由于目前网络上主流的websocket测试工具仅支持对标准的websocket连接进行
测试，对于我们的项目来说并不满足，所以需要重新编写一个客户端，满足连接的
稳定性、高并发性测试等。

##实现
通过gorilla/websocket包对client进行实现，完成连接验证、心跳监测、
建立连接等流程。
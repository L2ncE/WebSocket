# ChatRoom——对WebSocket的初步运用

**[参考项目](https://github.com/gorilla/websocket/blob/master/examples/chat)**

## 项目结构

### API

将登陆注册以及聊天室的实现封装在了里面

### DAO

对mysql进行连接并且有其他所有数据访问操作，主要是登陆注册相关

### HUB（借用原项目说明）

The code for the `Hub` type is in [hub.go](https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go). The application's `main` function starts the hub's `run` method as a goroutine. Clients send requests to the hub using the `register`, `unregister` and `broadcast` channels.

The hub registers clients by adding the client pointer as a key in the `clients` map. The map value is always true.

The unregister code is a little more complicated. In addition to deleting the client pointer from the `clients` map, the hub closes the clients's `send` channel to signal the client that no more messages will be sent to the client.

The hub handles messages by looping over the registered clients and sending the message to the client's `send` channel. If the client's `send` buffer is full, then the hub assumes that the client is dead or stuck. In this case, the hub unregisters the client and closes the websocket.

### MAIN

主函数

### MIDDLEWARE

中间件，其中有JWT相关以及CORS跨域相关中间件

### MODULE

模型文件，里面含有项目所有结构类代码

### RESPONSE

封装了回应函数

### ROUTER

路由，使用GIN框架

### SERVICE

服务层，封装了登陆注册以及聊天室所需服务

## 功能实现样例


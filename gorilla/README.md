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

### 登录注册

#### `/register`

![image-20220413103000532](https://s2.loli.net/2022/04/13/DMThijHXsWdBoPx.png)

![image-20220413103105473](https://s2.loli.net/2022/04/13/7UtkMV1P2uLvi3w.png)

#### `/login`

![image-20220413103149716](https://s2.loli.net/2022/04/13/3gMHpkO4L9bh1le.png)

通过登录所得到的token就可以进入聊天室进行聊天

### 聊天室

#### `ws://localhost:8080/?room_id="房间名"`

##### 基础功能

在header中填写token，连接后系统会识别你的name，并欢迎你

![image-20220413103612870](https://s2.loli.net/2022/04/13/76LkHpi8WBMJ5ZS.png)

其他人进入此房间，也会收到通知

![image-20220413103723020](https://s2.loli.net/2022/04/13/Dz9ih8AJt7bLgOr.png)

在房间中发送消息，房间中的其他用户也能够看到

![image-20220413105057379](https://s2.loli.net/2022/04/13/jJfuOtiqXnZ9ImC.png)

![image-20220413105108059](https://s2.loli.net/2022/04/13/k6v1jhWuNKDbfIi.png)

但是处于其他房间的用户是不会收到的（多房间同时开启）

![image-20220413105158526](https://s2.loli.net/2022/04/13/VAdkUKRHSC68BXN.png)

这个时候blue房间的小红是无法看到的



退出房间其他用户也会受到消息

![image-20220413105720012](https://s2.loli.net/2022/04/13/hadFvKLp3o7MTBX.png)

##### 限制发言功能（禁言）

```go
baseStr := "死傻操" //违法字符
		testStr := string(msg[:])
		for _, word := range testStr {
			//遍历是否有违法字符
			res := strings.Contains(baseStr, string(word))
			if res == true {
				c.limitNum += 1
				c.forbiddenWord = true //禁言
				//记录禁言开始时间
				c.timeLog = time.Now().Unix()
				h.warnings <- m
				break
```

发送有限制词的内容会禁言五分钟（使用time包），超过三次会永久禁言

![image-20220413105636879](https://s2.loli.net/2022/04/13/Nymqn5eQsvJw6tc.png)

至此作业要求的功能基本都已实现

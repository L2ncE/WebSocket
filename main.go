package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// 实现多房间聊天室思路
//创建map以roomName为key，访问者请求连接为value;value可以是map，也可以是切片数组;
//因为房间内有多个请求连接，为记录用户和请求，此处以用户id为key，请求为value;
// {roomid：{uid：conn}}
var rooms = make(map[string]map[int]*websocket.Conn)

var (
	// 服务器应用程序从HTTP请求处理程序调用Upgrader.Upgrade方法以获取* Conn;
	upgraderdd = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 连接websocket进入房间，首先在请求里指定要进入的房间
//  （也可以先进入个人房间，再转入其他房间都可以的，就看怎么处理连接对象conn，在这里先指定房间了）
// 所以在请求里需要带两个参数，房间名roomid和用户id
func wsEcho(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomid := r.Form["roomid"][0] //从请求里获取房价名roomid
	uidd := r.Form["uid"][0]      // 从请求里获取用户id
	uid, _ := strconv.Atoi(uidd)
	fmt.Println("roomid:===", roomid)
	fmt.Println("uid:==", uid)
	// conn就是建立连接后的连接对象
	conn, _ := upgraderdd.Upgrade(w, r, nil)

	defer conn.Close()

	func(conn *websocket.Conn) {
		if rooms[roomid] == nil {
			rooms[roomid] = map[int]*websocket.Conn{uid: conn}
		} else {
			rooms[roomid][uid] = conn
		}
	}(conn)
	fmt.Println("rooms:==", rooms)

	for {
		// 某个链接发来数据data
		msgTy, data, _ := conn.ReadMessage()
		rooms[roomid][uid] = conn //把房间名和用户id及连接对象conn保存到map里

		//把数据返给当前房间里(除自己)所有的链接
		for k, v := range rooms[roomid] {
			if k == uid {
				continue
			}
			err := v.WriteMessage(msgTy, data)
			if err != nil {
				log.Println("error:==", err)
			}
			log.Printf("Write msg to client: recved: %s \n", data)
		}
	}

}

func main() {
	http.HandleFunc("/", wsEcho) // ws://127.0.0.1:8888/rom
	// 监听 地址 端口
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

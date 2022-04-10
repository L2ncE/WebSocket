package main

import (
	"fmt"
	"log"
	"net/http"
)

type message struct {
	data   []byte
	roomid string
	conn   *connection
}

type hub struct {
	rooms       map[string]map[*connection]bool
	broadcast   chan message
	broadcastss chan message
	warnings    chan message
	register    chan message
	unregister  chan message
	kickoutroom chan message
	warnmsg     chan message
}

var h = hub{
	broadcast:   make(chan message),
	broadcastss: make(chan message),
	warnings:    make(chan message),
	warnmsg:     make(chan message),
	register:    make(chan message),
	unregister:  make(chan message),
	kickoutroom: make(chan message),
	rooms:       make(map[string]map[*connection]bool),
}

func (h *hub) run() {

	for {
		select {
		case m := <-h.register: //传输链接
			conns := h.rooms[m.roomid]
			if conns == nil { // 链接保存到相应的房间
				conns = make(map[*connection]bool)
				h.rooms[m.roomid] = conns
				fmt.Println("在线人数:==", len(conns))
				fmt.Println("rooms:==", h.rooms)
			}
			h.rooms[m.roomid][m.conn] = true
			fmt.Println("在线人数:==", len(conns))
			fmt.Println("rooms:==", h.rooms)

			for con := range conns {
				delmsg := "系统消息：欢迎新伙伴加入" + m.roomid + "聊天室！！！"
				data := []byte(delmsg)
				select {
				case con.send <- data:
				}
			}

		case m := <-h.unregister: //断开链接
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn) //删除链接
					close(m.conn.send)
					for con := range conns {
						delmsg := "系统消息：有小伙伴离开了" + m.roomid + "聊天室！欢送！！！"
						data := []byte(delmsg)
						select {
						case con.send <- data:
						}
						if len(conns) == 0 { // 链接都断开，删除房间
							delete(h.rooms, m.roomid)
						}
					}
				}
			}

		case m := <-h.kickoutroom: //3次不合法信息后，被踢出群聊
			conns := h.rooms[m.roomid]
			notice := "由于您多次发送不合法信息,已被踢出群聊！！！"
			select {
			case m.conn.send <- []byte(notice):
			}
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn)
					close(m.conn.send)
					if len(conns) == 0 {
						delete(h.rooms, m.roomid)
					}
				}
			}

		case m := <-h.warnings: //不合法信息警告
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "警告:您发布不合法信息，将禁言5分钟，三次后将被踢出群聊！！！"
					//starttime:=
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}

		case m := <-h.warnmsg: //禁言中提示
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "您还在禁言中,暂时不能发送信息！！！"
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}

		case m := <-h.broadcast: //传输群信息/房间信息
			conns := h.rooms[m.roomid]
			for con := range conns {
				if con == m.conn { //自己发送的信息，不用再发给自己
					continue
				}
				select {
				case con.send <- m.data:
				default:
					close(con.send)
					delete(conns, con)
					if len(conns) == 0 {
						delete(h.rooms, m.roomid)
					}
				}
			}

		case m := <-h.broadcastss: //传输全员广播信息
			for _, conns := range h.rooms {
				for con := range conns {
					if con == m.conn { //自己发送的信息，不用再发给自己
						continue
					}
					select {
					case con.send <- m.data:
					default:
						close(con.send)
						delete(conns, con)
						if len(conns) == 0 {
							delete(h.rooms, m.roomid)
						}
					}
				}

			}

		}
	}
}

// 要运行同包下的几个文件，都要执行
//go run server.go hub.go
func main() {

	go h.run()

	http.HandleFunc("/", serverWs)

	err := http.ListenAndServe(":8899", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

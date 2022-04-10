package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type connection struct {
	ws            *websocket.Conn
	send          chan []byte
	numberv       int
	forbiddenword bool
	timelog       int64
}

func (m message) readPump() {
	c := m.conn

	defer func() {
		h.unregister <- m
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			//log.Printf("error: %v", err)
			break
		}
		go m.Kickout(msg)

	}
}

// 信息处理，不合法言论 禁言警告，超过3次，踢出群聊；
func (m message) Kickout(msg []byte) {
	c := m.conn
	// 判断是否有禁言时间,并超过5分钟禁言时间,没有超过进入禁言提醒
	nowT := int64(time.Now().Unix())
	if nowT-c.timelog < 300 {
		h.warnmsg <- m
	}
	// 不合法信息3次，判断是否有不合法信息，没有进行信息发布
	if c.numberv < 3 {
		// 信息过滤，设置含有此字符串内字符为不合法信息，
		basestr := "死亡崩薨"
		teststr := string(msg[:])
		for _, ev := range teststr {
			// 判断字符串中是否含有某个字符，true/false
			reslut := strings.Contains(basestr, string(ev))
			if reslut == true {
				c.numberv += 1
				c.forbiddenword = true // 禁言为真
				// 记录禁言开始时间,禁言时间内任何信息不能发送
				c.timelog = int64(time.Now().Unix())
				h.warnings <- m
				break
			}
		}
		// 不禁言，消息合法 可以发送
		if c.forbiddenword != true {
			// 设置广播消息, 所有房间内都可以收到信息;给广播消息开头加一个特定字符串为标识，当然也有其他方法;
			// 此例 设置以开头0为标识, 之后去掉0 ;
			if msg[0] == 48 {
				head := string("所有玩家请注意:")
				data := head + string(msg[1:])
				m := message{[]byte(data), m.roomid, c}
				h.broadcastss <- m
			} else if msg[0] != 48 { //不是0，就是普通消息
				m := message{msg, m.roomid, c}
				h.broadcast <- m
			}
		}

		// 不合法信息超过三次，踢出群
	} else {
		h.kickoutroom <- m
		log.Println("要被踢出群聊了...")
		c.ws.Close() // 此处关闭了踢出的连接,也可以不关闭做其他操作,
	}

}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *message) writePump() {
	c := s.conn

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serverWs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomid := r.Form["roomid"][0]

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	m := message{nil, roomid, c}

	h.register <- m

	go m.writePump()
	go m.readPump()
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (m message) readPump() {
	c := m.conn

	defer func() {
		h.unregister <- m
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	c.ws.SetPongHandler(func(string) error {
		err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("unexpected close error:", err)
			}
			fmt.Println("err:", err)
			break
		}
		go m.Limit(msg)
	}
}

// Limit 发言违规进行时间限制禁言,超过三次禁言
func (m message) Limit(msg []byte) {
	c := m.conn
	//查看是否禁言，并判断是否超过了禁言时间
	timeNow := time.Now().Unix()
	if timeNow-c.timeLog < 300 {
		h.warnmsg <- m
	}

	// 不合法信息3次，判断是否有不合法信息，没有进行信息发布
	if c.limitNum >= 3 {
		h.kickoutroom <- m
		log.Println("素质太低，给你踢出去")
		c.ws.Close() //
	} else //没有超过三次，可以继续
	{
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
			}
		}
		// 不禁言，消息合法 可以发送
		if c.forbiddenWord != true {
			// 通过所有检查，进行广播

			if msg[0] == 48 {
				head := "所有玩家请注意:"
				data := head + string(msg[1:])
				m := message{[]byte(data), m.roomId, c}
				h.broadcastss <- m
			} else if msg[0] != 48 { //不是0，就是普通消息
				m := message{msg, m.roomId, c}
				h.broadcast <- m
			}
		}
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

func serverWs(ctx *gin.Context) {
	ctx.Request.ParseForm()
	roomid := ctx.Request.Form["roomid"][0]

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	m := message{nil, roomid, c}

	h.register <- m

	go m.writePump()
	go m.readPump()
}

package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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

			m := message{msg, m.roomId, m.name, c}
			h.broadcast <- m
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

// ChangePassword 修改密码服务
func ChangePassword(username, newPassword string) error {
	err := UpdatePassword(username, newPassword)
	return err
}

// IsPasswordCorrect 判断密码是否正确服务
func IsPasswordCorrect(username, password string) (bool, error) {
	user, err := SelectUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if user.Password != password {
		return false, nil
	}

	return true, nil
}

// IsRepeatUsername 判断用户名是否重复
func IsRepeatUsername(username string) (bool, error) {
	_, err := SelectUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Register 注册服务
func Register(user User) error {
	err := Insert(user)
	return err
}

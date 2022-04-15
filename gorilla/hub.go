package main

func (h *hub) run() {
	for {
		select {
		case m := <-h.register:
			conns := h.rooms[m.roomId]
			if conns == nil {
				conns = make(map[*connection]bool)
				h.rooms[m.roomId] = conns
			}
			h.rooms[m.roomId][m.conn] = true

			for con := range conns {
				sysmsg := "系统消息：欢迎新伙伴" + m.name + "加入" + m.roomId + "聊天室！！！"
				data := []byte(sysmsg)
				select {
				case con.send <- data:
				}
			}

		case m := <-h.unregister: //断开链接
			conns := h.rooms[m.roomId]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn) //删除链接
					close(m.conn.send)
					for con := range conns {
						delMsg := "系统消息：" + m.name + "离开了" + m.roomId + "聊天室"
						data := []byte(delMsg)
						select {
						case con.send <- data:
						}
						if len(conns) == 0 {
							delete(h.rooms, m.roomId)
						}
					}
				}
			}

		case m := <-h.kickoutroom: //3次不合法信息后，被踢出群聊
			conns := h.rooms[m.roomId]
			notice := "由于您多次发送不合法信息,已被踢出群聊！！！"
			select {
			case m.conn.send <- []byte(notice):
			}
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn)
					close(m.conn.send)
					if len(conns) == 0 {
						delete(h.rooms, m.roomId)
					}
				}
			}

		case m := <-h.warnings:
			conns := h.rooms[m.roomId]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "警告:您发布不合法信息，将禁言5分钟，三次后将被踢出群聊！！！"
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}

		case m := <-h.warnmsg: //禁言中提示
			conns := h.rooms[m.roomId]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "您还在禁言中,暂时不能发送信息！！！"
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}

		case m := <-h.broadcast: //传输群信息/房间信息
			conns := h.rooms[m.roomId]
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
						delete(h.rooms, m.roomId)
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
							delete(h.rooms, m.roomId)
						}
					}
				}
			}
		}
	}
}

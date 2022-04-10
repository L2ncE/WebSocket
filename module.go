package main

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type connection struct {
	ws            *websocket.Conn
	send          chan []byte
	numberv       int
	forbiddenword bool
	timelog       int64
}

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

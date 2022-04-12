package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.GET("/", serverWs)
	go h.run()
	err := engine.Run()
	if err != nil {
		fmt.Printf("init error:%v\n", err)
		return
	}
}

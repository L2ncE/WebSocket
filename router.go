package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.GET("/", serverWs)
	err := engine.Run(":8899")
	if err != nil {
		fmt.Printf("init error:%v\n", err)
		return
	}
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.POST("/register", register)
	engine.POST("/login", login)
	userGroup := engine.Group("/user")
	userGroup.Use(CORS())
	{
		userGroup.Use(JWTAuth)
		userGroup.POST("/password", changePassword)
		userGroup.GET("/ws", serverWs)
	}
	go h.run()
	err := engine.Run()
	if err != nil {
		fmt.Printf("init error:%v\n", err)
		return
	}
}

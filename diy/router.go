package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.Use(CORS())
	engine.POST("/register", register)
	engine.POST("/login", login)

	userGroup := engine.Group("/user")
	{
		userGroup.Use(JWTAuth)
		userGroup.PUT("/password", changePassword)
	}

	wsGroup := engine.Group("/")
	{
		wsGroup.Use(JWTAuth)
		wsGroup.GET("/", serverWs)
	}

	go h.run()
	err := engine.Run()
	if err != nil {
		fmt.Printf("init error:%v\n", err)
		return
	}
}

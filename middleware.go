package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var mySigningKey = []byte("RedRock")

// JWTAuth JWT登录
func JWTAuth(ctx *gin.Context) {
	token := ctx.Request.Header.Get("token")
	if token == "" {
		RespErrorWithData(ctx, "游客你好！没有您的信息,请先登录!")
		ctx.Abort()
		return
	}
	ctx.Set("username", ParseToken(token))
	ctx.Next()
}

func ParseToken(s string) string {
	//解析传过来的token
	tokenClaims, err := jwt.ParseWithClaims(s, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return tokenClaims.Claims.(*MyClaims).Username
}

//CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, token")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}

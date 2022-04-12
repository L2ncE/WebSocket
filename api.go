package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

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

func changePassword(ctx *gin.Context) {
	oldPassword := ctx.PostForm("old_password")
	newPassword := ctx.PostForm("new_password")
	iUsername, _ := ctx.Get("username")
	l1 := len([]rune(newPassword))
	if l1 <= 16 && l1 >= 6 { //强制规定密码小于16位并大于6位
		username := iUsername.(string)

		//检验旧密码是否正确
		flag, err := IsPasswordCorrect(username, oldPassword)
		if err != nil {
			fmt.Println("judge password correct err: ", err)
			RespInternalError(ctx)
			return
		}
		if !flag {
			RespErrorWithData(ctx, "旧密码输入错误")
			return
		}
		//修改新密码
		err = ChangePassword(username, newPassword)
		if err != nil {
			fmt.Println("change password err: ", err)
			RespInternalError(ctx)
			return
		}
		RespSuccessfulWithData(ctx, "已修改")
	} else {
		RespErrorWithData(ctx, "密码请在6位到16位之内")
		return
	}
}

func login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	flag, err := IsPasswordCorrect(username, password)
	if err != nil {
		fmt.Println("judge password correct err: ", err)
		RespInternalError(ctx)
		return
	}
	if !flag {
		RespErrorWithData(ctx, "密码错误")
		return
	}
	//jwt
	c := MyClaims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,
			ExpiresAt: time.Now().Unix() + 6000000,
			Issuer:    "YuanXinHao",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, err := t.SignedString(mySigningKey)
	if err != nil {
		RespInternalError(ctx)
	}
	RespSuccessfulWithData(ctx, s)
}

func register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	//输入信息不能为空
	if username != "" && password != "" {
		l1 := len([]rune(username))
		l2 := len([]rune(password))
		if l1 <= 8 && l1 >= 1 { //强制规定用户名长度小于8位大于1位
			if l2 <= 16 && l2 >= 6 { //强制规定密码小于16位大于6位
				user := User{
					Name:     username,
					Password: password,
				}
				flag, err := IsRepeatUsername(username)
				if err != nil {
					fmt.Println("judge repeat username err: ", err)
					RespInternalError(ctx)
					return
				}
				if flag {
					RespErrorWithData(ctx, "用户名已经存在")
					return
				}
				err = Register(user)
				if err != nil {
					fmt.Println("register err: ", err)
					RespInternalError(ctx)
					return
				}
				RespSuccessfulWithData(ctx, "已注册")
			} else {
				RespErrorWithData(ctx, "密码请在6位到16位之内")
				return
			}
		} else {
			RespErrorWithData(ctx, "用户名请在1位到8位之内")
			return
		}
	} else {
		RespErrorWithData(ctx, "请将信息输入完整")
		return
	}
}

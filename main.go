package main

import "fmt"

func main() {
	err := InitGormDB()
	if err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
	} else {
		fmt.Println("连接GORM MySQL数据库成功!")
	}
	InitEngine()
}

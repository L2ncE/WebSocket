package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func InitGormDB() (err error) {
	dB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      "chatroom1:yxh030714@tcp(42.192.155.29:3306)/chatroom1?charset=utf8mb4&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:        171,
		DisableDatetimePrecision: true,
		DontSupportRenameIndex:   true,
	}), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Printf("连接失败：%v\n", err)
	}
	db = dB
	return err
}

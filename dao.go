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

func UpdatePassword(Name string, newPassword string) error {
	deRes := db.Model(&User{}).Where("Name = ?", Name).Update("Password", newPassword)
	err := deRes.Error
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return err
	}
	return err
}

func SelectUserByUsername(Name string) (User, error) {
	var user User
	dbRes := db.Model(&User{}).Select("id", "password").Where("Name = ?", Name).First(&user)
	err := dbRes.Error
	if err != nil {
		return user, err
	}
	fmt.Println(user)
	return user, nil
}

func Insert(user User) error {
	deres := db.Select("Name", "Password").Create(&User{Name: user.Name, Password: user.Password})
	err := deres.Error
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return err
	}
	return err
}

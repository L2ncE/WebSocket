package main

import (
	"database/sql"
	"gorm.io/gorm"
)

// ChangePassword 修改密码服务
func ChangePassword(username, newPassword string) error {
	err := UpdatePassword(username, newPassword)
	return err
}

// IsPasswordCorrect 判断密码是否正确服务
func IsPasswordCorrect(username, password string) (bool, error) {
	user, err := SelectUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if user.Password != password {
		return false, nil
	}

	return true, nil
}

// IsRepeatUsername 判断用户名是否重复
func IsRepeatUsername(username string) (bool, error) {
	_, err := SelectUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Register 注册服务
func Register(user User) error {
	err := Insert(user)
	return err
}

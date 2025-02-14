package service

import (
	"database/sql"
	"decleration/dao"
	"decleration/model"
	"decleration/utils1"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password string) error {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	//哈希加密
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 插入数据
	query := "INSERT INTO user (username, password) VALUES (?, ?)"
	_, err = db.Exec(query, username, hashPassword)
	if err != nil {
		return fmt.Errorf("插入用户失败: %w", err)
	}

	return nil
}

func LoginUser(username, password string) (string, error) {
	// 查询用户信息
	storedUser, err := dao.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("用户名或密码错误")
		}
		return "", fmt.Errorf("查询数据库失败: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password)) != nil {
		return "", fmt.Errorf("用户名或密码错误")
	}

	// 生成 Token
	token, err := utils1.GenerateToken(username)
	if err != nil {
		return "", fmt.Errorf("生成 Token 失败: %w", err)
	}

	return token, nil
}
func ChangePasswordService(username, oldPassword, newPassword string) error {
	db, err := dao.InitDB()
	if err != nil {
		return fmt.Errorf("数据库连接失败")
	}
	defer db.Close()

	// 查询用户的旧密码
	var storedPassword string
	err = db.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("用户名不存在")
		}
		return fmt.Errorf("查询用户失败")
	}

	// 检查旧密码是否匹配
	if storedPassword != oldPassword {
		return fmt.Errorf("旧密码错误")
	}

	// 更新密码
	_, err = db.Exec("UPDATE user SET password=? WHERE username=?", newPassword, username)
	if err != nil {
		return fmt.Errorf("更新密码失败")
	}

	return nil
}

func GetUserInfoService(userID string) (*model.UserInfo, error) {
	db, err := dao.InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败")
	}
	defer db.Close()

	var user model.UserInfo
	query := "SELECT id, nickname, introduction, phone, qq, gender, email, birthday FROM userinfo WHERE id = ?"
	err = db.QueryRow(query, userID).Scan(
		&user.Id, &user.Nickname, &user.Introduction, &user.Phone, &user.Qq,
		&user.Gender, &user.Email, &user.Birthday,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户信息失败")
	}

	return &user, nil
}

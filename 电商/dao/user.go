package dao

import (
	"database/sql"
	"decleration/model"
	"fmt"
)

// InitDB 初始化数据库连接
func InitDB() (*sql.DB, error) {
	// 数据库连接字符串
	dsn := "root:hj2005691@tcp(127.0.0.1:3306)/my_db_01" // 注意 MySQL 的地址和端口需要正确配置

	// 尝试打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接错误: %v", err) // 返回具体的错误信息
	}

	// 确保连接有效
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err) // 返回具体的错误信息
	}

	return db, nil // 返回数据库连接和 nil 错误
}

func DeleteComment(db *sql.DB, commentID string) error {
	query := "DELETE FROM comment WHERE comment_id = ?"
	_, err := db.Exec(query, commentID)
	return err
}
func GetUserByUsername(username string) (*model.User, error) {
	db, err := InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	var user model.User
	err = db.QueryRow("SELECT username, password FROM user WHERE username = ?", username).
		Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

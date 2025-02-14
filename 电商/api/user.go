package api

import (
	"context"
	"decleration/dao"
	"decleration/model"
	"decleration/service"
	"decleration/utils1"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

func Register(ctx context.Context, c *app.RequestContext) {
	// 解析参数
	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "数据解析失败"})
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, utils.H{"message": "用户名和密码不能为空"})
		return
	}

	// 调用 service 层处理注册逻辑
	err := service.RegisterUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "注册失败"})
		return
	}

	// 注册成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "注册成功",
	})
}

func Login(ctx context.Context, c *app.RequestContext) {
	// 解析请求体
	var req model.User

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "请求体解析失败"})
		return
	}

	// 调用 service 处理登录逻辑
	token, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "登录失败"})
		return
	}

	// 返回成功信息和 Token
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "登录成功",
		"token":   token,
	})
}

func RefreshToken(ctx context.Context, c *app.RequestContext) {
	type TokenRequest struct {
		OldToken string `json:"oldToken"`
	}
	// 从请求中获取旧的 Token
	var oldToken TokenRequest
	if err := c.Bind(&oldToken); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, utils.H{"message": "请求体解析失败"})
		return
	}

	if oldToken.OldToken == "" {
		c.JSON(400, utils.H{"error": "Token 不能为空"})
		return
	}

	// 调用刷新 Token 的工具函数
	newToken, err := utils1.RefreshToken(oldToken.OldToken)
	if err != nil {
		c.JSON(401, utils.H{"error": err.Error()})
		return
	}

	// 返回新的 Token
	c.JSON(200, utils.H{
		"info":          "success",
		"status":        10000,
		"refresh_token": newToken})
}

func ChangePassword(ctx context.Context, c *app.RequestContext) {
	auth, _ := utils1.ExtractToken(c)
	username, _ := utils1.ValidateToken(auth)

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind(&request); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "解析请求体失败"})
		return
	}

	// 调用 service 处理修改密码
	if err := service.ChangePasswordService(username, request.OldPassword, request.NewPassword); err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"inf":     "success",
		"status":  10000,
		"message": "密码修改成功"})
}

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	auth, _ := utils1.ExtractToken(c)
	username, _ := utils1.ValidateToken(auth)

	user, err := service.GetUserInfoService(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"info":   "success",
		"status": 10000,
		"user":   user,
	})
}

func ChangeUserInfo(ctx context.Context, c *app.RequestContext) {
	// 获取 Bearer Token
	auth, _ := utils1.ExtractToken(c)
	userID, err := utils1.ValidateToken(auth) // 验证 Token
	if err != nil {
		c.JSON(consts.StatusUnauthorized, utils.H{"message": "无效的认证令牌"})
		return
	}

	// 连接数据库
	db, _ := dao.InitDB()
	defer db.Close()

	// 解析请求体
	var updateData map[string]interface{}
	if err := json.Unmarshal(c.Request.Body(), &updateData); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "请求参数错误"})
		return
	}

	// 构建 SQL 更新语句
	var setClauses []string
	var values []interface{}

	allowedFields := map[string]bool{
		"nickname": true, "introduction": true, "phone": true,
		"qq": true, "gender": true, "email": true, "birthday": true,
	}

	for key, value := range updateData {
		if allowedFields[key] {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, utils.H{"message": "没有提供要更新的字段"})
		return
	}

	// 追加 WHERE 子句
	values = append(values, userID)
	query := fmt.Sprintf("UPDATE userinfo SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	// 执行更新操作
	result, err := db.Exec(query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "更新失败"})
		return
	}

	// 检查是否有行受影响
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, utils.H{"message": "用户不存在或未修改任何信息"})
		return
	}

	c.JSON(http.StatusOK, utils.H{"message": "用户信息更新成功"})
}

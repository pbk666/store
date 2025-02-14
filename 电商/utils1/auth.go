package utils1

import (
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

var SecretKey = []byte("your-secret-key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成 JWT Token
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID, // 这里存储用户ID或其他相关信息
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Token 过期时间为 24 小时
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// 解析 JWT Token
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法是否符合预期
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SecretKey, nil
	})
}

// RefreshToken 刷新 JWT Token
func RefreshToken(oldToken string) (string, error) {
	// 解析旧的 Token
	token, err := ParseToken(oldToken)
	if err != nil {
		return "", errors.New("无效的 Token")
	}

	// 检查 Token 是否有效
	if !token.Valid {
		return "", errors.New("Token 无效")
	}

	// 获取 Token 中的声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Token 声明无效")
	}

	// 获取用户 ID
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("Token 中缺少用户 ID")
	}

	// 生成新的 Token
	newToken, err := GenerateToken(userID)
	if err != nil {
		return "", errors.New("生成新 Token 失败")
	}

	return newToken, nil
}

// 验证token
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 安全地获取 username 字段
		username, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("token does not contain a valid username")
		}
		return username, nil
	}

	return "", fmt.Errorf("invalid toke")
}

// 提取token
func ExtractToken(c *app.RequestContext) (string, error) {
	auth := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}
	if len(auth) == 0 {
		return "", errors.New("token不能为空")
	}
	return auth, nil
}

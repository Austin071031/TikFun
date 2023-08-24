package utils

import (
	_ "errors"
	"fmt"
	"time"
	//"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

var secretchar = "JhsgqWusd"

// 定义自定义的声明结构，可以根据实际需求添加更多声明
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 密钥，用于签署令牌
var jwtSecret = []byte(secretchar) //自己设置，随机的一串英文字符，一般是随机生成的

func GenerateToken(username string) (string, error) {
	// 设置token过期时间，这里设置为1小时
	expirationTime := time.Now().Add(1 * time.Hour)

	// 创建声明
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签署令牌
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 验证Token
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func VerifyToken(tokenString string) (string, error) {
	// 解析Token并验证其有效性
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 返回用户ID
	return claims.Username, nil
}


// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 常量
const (
	SecretKey = "library-management-jwt-secret-2024"
	ExpireHours = 24
)

// Claims JWT负载结构体
type Claims struct {
	UserId   uint64 `json:"user_id"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Service JWT服务
type Service struct{}

// New 创建JWT服务实例
func New() *Service {
	return &Service{}
}

// GenerateToken 生成JWT Token
func (s *Service) GenerateToken(userId uint64, phone, email, role string) (string, error) {
	claims := Claims{
		UserId: userId,
		Phone:  phone,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireHours * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "library-management-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析并验证JWT Token
func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的Token")
	}

	return claims, nil
}

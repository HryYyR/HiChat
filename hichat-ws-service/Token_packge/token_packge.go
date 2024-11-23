package Token_packge

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"time"
)

// GenerateToken token
func GenerateToken(id int, UUID, name string, expiretime time.Duration) (string, error) {
	uc := models.UserClaim{
		ID:       id,
		UUID:     UUID,
		UserName: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiretime)), // 定义过期时间 单位:分钟
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	ttoken, err := token.SignedString([]byte(config.JwtKey))
	if err != nil {
		fmt.Println("jwt error: ", err)
		return "", err
	}
	return ttoken, nil
}

// DecryptToken 解密token
func DecryptToken(token string) (*models.UserClaim, error) {
	uc := new(models.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(tk *jwt.Token) (any, error) {
		return []byte(config.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, nil
}

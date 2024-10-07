package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	random "crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"hichat_static_server/config"
	"hichat_static_server/models"
	"math/rand"
	"net/smtp"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
)

// Md5 md5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateEcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), random.Reader)
}

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

// MailSendCode 发送验证码邮件
func MailSendCode(mail string, code string) error {
	e := email.NewEmail()
	e.From = "HiChat <2452719312@qq.com>"
	e.To = []string{mail}
	e.Subject = "Code"                              //标题
	e.HTML = []byte("<h1>你的验证码: " + code + "</h1>") //内容
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", config.EmailAccount, config.EmailPassword, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		return err
	}
	return nil
}

// RandCode 生成随机验证码
func RandCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}

// GenerateUUID 生成随机UUID
func GenerateUUID() string {
	u1, _ := uuid.NewV4()
	return u1.String()
}

func TimeSortApplyJoinGroupList(arr []models.ApplyJoinGroup, order string) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].CreatedAt.After(arr[j].CreatedAt)
	})
}

func TimeSortAddUserList(arr []models.ApplyAddUser, order string) {
	switch order {
	case "desc":
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].CreatedAt.After(arr[j].CreatedAt)
		})
	case "asc":
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].CreatedAt.Before(arr[j].CreatedAt)
		})
	}

}

// HandleJsonArgument 处理server参数 json -> user struct
func HandleJsonArgument(c *gin.Context, data *models.Users) error {
	rawbyte, err := c.GetRawData()
	if err != nil {
		return err
	}
	// var data models.ApplyAddUser
	err = json.Unmarshal(rawbyte, &data)
	if err != nil {
		return err
	}
	return nil
}

func GroupTimeSort(arr []models.ApplyJoinGroupResponse, order string) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].CreatedAt.After(arr[j].CreatedAt)
	})
}
func UserTimeSort(arr []models.ApplyAddUser, order string) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].CreatedAt.After(arr[j].CreatedAt)
	})
}

func H(c *gin.Context, status int, msg string, err error) {
	c.JSON(status, gin.H{
		"msg":   msg,
		"error": err,
	})
}

package util

import (
	// "crypto/ecdsa"
	// "crypto/elliptic"
	"crypto/md5"
	"github.com/golang-jwt/jwt/v4"
	"net"

	// random "crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"go-websocket-server/config"
	"go-websocket-server/models"
	"log"
	"math/rand"
	"net/smtp"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jordan-wright/email"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// md5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// func generateEcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
// 	return ecdsa.GenerateKey(elliptic.P256(), random.Reader)
// }

// token
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

// 解密token
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

// 发送验证码邮件
func MailSendCode(mail string, code string) {
	e := email.NewEmail()
	e.From = "HiChat <2452719312@qq.com>"
	e.To = []string{mail}
	e.Subject = "Code"                              //标题
	e.HTML = []byte("<h1>你的验证码: " + code + "</h1>") //内容
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", config.EmailAccount, config.EmailPassword, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		log.Printf("发送邮箱验证码失败!%s\n", err)
	}
}

// 生成随机验证码
func RandCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}

// 生成随机UUID
func GenerateUUID() string {
	u1, _ := uuid.NewV4()
	return u1.String()
}

func GroupTimeSort(arr []models.ApplyJoinGroup, order string) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].CreatedAt.After(arr[j].CreatedAt)
	})
}
func UserTimeSort(arr []models.ApplyAddUser, order string) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].CreatedAt.After(arr[j].CreatedAt)
	})
}

// 处理server参数 json -> struct
func HandleJsonArgument(c *gin.Context, data any) error {
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

func FormatTampTime(tamptime *timestamppb.Timestamp) time.Time {
	return tamptime.AsTime().Local().UTC().Add(time.Hour * -8)
}

func FormatTime(targettime time.Time) time.Time {
	return targettime.Local().UTC().Add(time.Hour * -8)
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && ipnet.IP.String()[:3] != "169" { // IPv4 address
			ip := ipnet.IP.String()
			return ip
		}
	}
	return ""
}

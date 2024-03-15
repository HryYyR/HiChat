package util

import (
	// "crypto/ecdsa"
	// "crypto/elliptic"
	"crypto/md5"
	"net"
	"strconv"

	// random "crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-websocket-server/config"
	"log"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jordan-wright/email"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Md5 md5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// func generateEcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
// 	return ecdsa.GenerateKey(elliptic.P256(), random.Reader)
// }

// MailSendCode 发送验证码邮件
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

//func GroupTimeSort(arr []models.ApplyJoinGroup, order string) {
//	sort.Slice(arr, func(i, j int) bool {
//		return arr[i].CreatedAt.After(arr[j].CreatedAt)
//	})
//}
//func UserTimeSort(arr []models.ApplyAddUser, order string) {
//	sort.Slice(arr, func(i, j int) bool {
//		return arr[i].CreatedAt.After(arr[j].CreatedAt)
//	})
//}

// HandleJsonArgument 处理server参数 json -> struct
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

func H(c *gin.Context, status int, msg string, err error) {
	c.JSON(status, gin.H{
		"msg":   msg,
		"error": err,
	})
}

func IntArrToStrArr(intarr []int) []string {
	var strarr []string
	for _, in := range intarr {
		strarr = append(strarr, strconv.Itoa(in))
	}
	return strarr
}

func StrArrToIntArr(strarr []string) []int {
	var intarr []int
	for _, st := range strarr {
		str, err := strconv.Atoi(st)
		if err != nil {
			return []int{}
		}
		intarr = append(intarr, str)
	}
	return intarr
}

// DeleteSlice 删除指定元素。
func DeleteStrSlice(a []string, elem string) []string {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

// DeleteSlice 删除指定元素。
func DeleteIntSlice(a []int, elem int) []int {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"io"
	"strings"

	// "crypto/ecdsa"
	// "crypto/elliptic"
	"crypto/md5"
	randd "crypto/rand"
	"crypto/rsa"
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

// DeleteStrSlice DeleteSlice 删除指定元素。
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

// DeleteIntSlice DeleteSlice 删除指定元素。
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

// GenerateRsaKey 生成 rsa Key
func GenerateRsaKey() (*rsa.PublicKey, *rsa.PrivateKey) {
	privateKey, _ := rsa.GenerateKey(randd.Reader, 2048)
	publicKey := privateKey.PublicKey
	return &publicKey, privateKey
}

// DecryptRSA 解密 RSA 密文
func DecryptRSA(encryptedKey []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// 使用 RSA 私钥解密加密后的对称密钥
	decryptedKey, err := rsa.DecryptPKCS1v15(randd.Reader, privateKey, encryptedKey)
	if err != nil {
		return nil, err
	}
	return decryptedKey, nil
}

// DecryptAES DecryptData 解密 aes 密文
func DecryptAES(encryptedData, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New(" IV length must equal block size\n")
	}
	if err != nil {
		fmt.Println("解码失败:", err)
		return nil, err
	}

	// CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedData, encryptedData)
	// 去除填充
	decryptedData := unpad(encryptedData)

	//var decryptedDataJson EncryptedMessage
	//err = json.Unmarshal(decryptedData, &decryptedDataJson)
	//fmt.Printf("%v\n", decryptedDataJson)
	//marshal, err := json.Marshal(decryptedDataJson.Message)
	//if err != nil {
	//	return nil, err
	//}
	return decryptedData, nil
}

type EncryptedMessage struct {
	Message string
}
type EncryptedData struct {
	Iv      []byte
	Message []byte
}

// EncryptAESCBC 使用AES CBC模式对数据进行加密
func EncryptAESCBC(data []byte, key []byte) (EncryptedData, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptedData{}, err
	}
	// 创建一个随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(randd.Reader, iv); err != nil {
		return EncryptedData{}, err
	}
	plaintext := pkcs7Pad(data, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext)) //用时2天解决的bug,在调用CryptBlocks时
	// 使用CBC模式加密数据
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext) //必须用一个新的变量来接加密数据,不然两个协程之间回共享变量,导致数据重复
	// 返回加密数据和IV
	return EncryptedData{Message: ciphertext, Iv: iv}, nil
}

// 去除填充
func unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// Pad 对数据进行填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func IntsToBytes(ints []int) ([]byte, error) {
	// 创建一个字节缓冲区
	buf := new(bytes.Buffer)
	// 将 []int 写入缓冲区
	for _, i := range ints {
		err := binary.Write(buf, binary.LittleEndian, int32(i))
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func BytesToInts(bytess []byte) ([]int, error) {
	// 创建一个缓冲区
	buf := bytes.NewBuffer(bytess)
	ints := make([]int, 0, len(bytess)/4)
	for {
		var i int32
		err := binary.Read(buf, binary.LittleEndian, &i)
		if err != nil {
			break
		}
		ints = append(ints, int(i))
	}
	return ints, nil
}

// 清除字符串中多余的逗号
func removeConsecutiveCommas(s string) string {
	var builder strings.Builder
	lastChar := rune(0) // 记录上一个字符，用于判断是否是逗号

	for _, char := range s {
		if char != ',' || lastChar != ',' {
			// 如果当前字符不是逗号，或者上一个字符不是逗号，则添加到builder中
			builder.WriteRune(char)
			lastChar = char
		}
	}

	return builder.String()
}

package test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"go-websocket-server/util"
	"testing"
)

//
//func TestXormTest(t *testing.T) {
//	adb.InitMySQL()
//	// err := adb.Engine.CreateTables(&models.User{})
//	err := adb.SqlStruct.Conn.Sync2(
//		//new(models.UserUnreadMessage),
//		//new(models.UserMessage),
//		//new(models.UserUserRelative),
//		//new(models.ApplyJoinGroup),
//		//new(models.ApplyAddUser),
//		//new(models.GroupUnreadMessage),
//		//new(models.Users),
//		//new(models.Group),
//		new(models.GroupMessage),
//		//new(models.GroupUserRelative),
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func TestRsa(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
		return
	}

	publicKey := privateKey.PublicKey

	fmt.Println(privateKey)
	fmt.Println(publicKey)

	var pubKeyBytes bytes.Buffer

	encoder := gob.NewEncoder(&pubKeyBytes)
	_ = encoder.Encode(publicKey)

	fmt.Println(pubKeyBytes)
	fmt.Println(pubKeyBytes.Bytes())
	//fmt.Println("")
	//fmt.Println(publicKey)
	//var receivedPubKey rsa.PublicKey
	//decoder := gob.NewDecoder(&pubKeyBytes)
	//err = decoder.Decode(&receivedPubKey)
	//fmt.Println("Received public key:", receivedPubKey)
}

func TestAes(t *testing.T) {

	ivBase64 := ""
	keyBase64 := ""

	// 解码密钥和 IV
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		fmt.Println("解码密钥失败:", err)
		t.Fatal("解码密钥失败:", err)
		return
	}

	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		fmt.Println("解码 IV 失败:", err)
		t.Fatal("解码 IV 失败:", err)
		return
	}

	// 待解密的数据，这里假设为 base64 编码的字符串
	encryptedData := "your_encrypted_data_base64_encoded"

	// 解码待解密的数据
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		fmt.Println("解码待解密的数据失败:", err)
		t.Fatal("解码待解密的数据失败:", err)
		return
	}

	// 使用 AES 密钥创建一个解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("创建解密器失败:", err)
		t.Fatal("创建解密器失败:", err)
		return
	}

	// 使用 CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	plaintext = unpad(plaintext)

	// 输出解密后的数据
	fmt.Println("解密后的数据:", string(plaintext))

}

// 反向 PKCS7 填充
func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func TestGetIP(t *testing.T) {
	ip := util.GetIP()
	fmt.Println(ip)
}

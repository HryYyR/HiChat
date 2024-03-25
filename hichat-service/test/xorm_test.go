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
//
//func TestXormsql(t *testing.T) {
//	adb.InitMySQL()
//	// err := adb.Engine.CreateTables(&models.User{})
//	_, err := adb.SqlStruct.Conn.Table("user_unread_message").Insert(&models.UserUnreadMessage{
//		UserName:     "666",
//		UserID:       8,
//		FriendID:     9,
//		UnreadNumber: 0,
//		CreatedAt:    time.Time{},
//		DeletedAt:    time.Time{},
//		UpdatedAt:    time.Time{},
//	})
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//}
//
//func TestGetUserGroupList(t *testing.T) {
//	userdata := models.UserClaim{
//		ID:       2,
//		UUID:     "2028bdd7d36f",
//		UserName: "niko",
//	}
//	adb.InitMySQL()
//
//	grouplist := make(map[int]models.Group, 0)
//
//	var groupidlist []int
//	if err := adb.SqlStruct.Conn.Cols("group_id").Table("group_user_relative").Where("user_id=?", userdata.ID).Find(&groupidlist); err != nil {
//		fmt.Println(err.Error())
//		t.Fatal(err.Error())
//	}
//	fmt.Printf("%+v\n", groupidlist)
//	for _, groupid := range groupidlist {
//		var groupitem models.Group
//		if _, err := adb.SqlStruct.Conn.Table("group").Where("id=?", groupid).Get(&groupitem); err != nil {
//			t.Fatal(err)
//		}
//		grouplist[groupid] = groupitem
//	}
//	fmt.Printf("%+v\n", grouplist)
//}
//
//func TestUserToClient(t *testing.T) {
//	adb.InitMySQL()
//	mockservicecenter := make(map[int]models.UserClient, 0)
//
//	var userdatalist []models.Users
//	err := adb.SqlStruct.Conn.Table("users").Find(&userdatalist)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, userdata := range userdatalist {
//		uuid := util.GenerateUUID()
//		grouplist, err := models.GetUserGroupList(userdata.ID)
//		if err != nil {
//			t.Fatal(err)
//		}
//		client := models.UserClient{
//			ClientID: uuid,
//			UserID:   userdata.ID,
//			UserUUID: userdata.UUID,
//			UserName: userdata.UserName,
//			Send:     make(chan []byte, 256),
//			Status:   false,
//			Groups:   grouplist,
//			// CachingMessages: make(map[int]int, 0),
//		}
//		mockservicecenter[userdata.ID] = client
//	}
//}
//
//func TestSyncMsg(t *testing.T) {
//	adb.InitMySQL()
//	id := 1
//	var unreadmsglist []models.GroupUnreadMessage
//	err := adb.SqlStruct.Conn.Table("group_unread_message").Where("user_id = ?", id).Find(&unreadmsglist)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("%+v\n", unreadmsglist)
//}
//
//func TestGetIP(t *testing.T) {
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		fmt.Println("Error:", err)
//		t.Fatal(err)
//	}
//
//	for _, addr := range addrs {
//		ipnet, ok := addr.(*net.IPNet)
//		if !ok {
//			continue
//		}
//		if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && ipnet.IP.String()[:3] != "169" { // IPv4 address
//			ip := ipnet.IP.String()
//			fmt.Println(ip)
//		}
//	}
//}
//
//func TestRefreshRedis(t *testing.T) {
//	adb.InitMySQL()
//	adb.InitRedis()
//	var groupdata []models.Group
//	err := adb.SqlStruct.Conn.Table("group").Find(&groupdata)
//	fmt.Println(len(groupdata))
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, g := range groupdata {
//		var useridlist []int
//		err = adb.SqlStruct.Conn.Cols("user_id").Table("group_user_relative").Where("group_id=?", g.ID).Find(&useridlist)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		var struserlist []string
//		for _, id := range useridlist {
//			struserlist = append(struserlist, strconv.Itoa(id))
//		}
//		str := strings.Join(struserlist, ",")
//		set := adb.Rediss.HSet("GroupToUserMap", strconv.Itoa(g.ID), str)
//		fmt.Println(set)
//
//	}
//}

//func TestStringsSplitRedis(t *testing.T) {
//	str := "1234"
//	arr := strings.Split(str, ",")
//	fmt.Println(len(arr))
//	for i, s := range arr {
//		fmt.Println(i, s)
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

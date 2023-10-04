package service

import (
	"crypto/md5"
	"fmt"
	adb "hichat-file-service/ADB"
	"hichat-file-service/models"
	"hichat-file-service/util"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	//获取普通文本
	// 获取文件(注意这个地方的file要和html模板中的name一致)
	// fileheader, err := c.FormFile("file")
	file, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("获取数据失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "获取数据失败",
		})
		return
	}
	//文件大小
	filesize := fileheader.Size
	if filesize > 1024*1024*10 {
		// fmt.Println("文件过大", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "文件大小超过了10M",
		})
		return
	}
	//获取文件的后缀名
	extstring := path.Ext(fileheader.Filename)
	fmt.Println(extstring)
	extmap := map[string]string{
		".png":  "png",
		".jpg":  "jpg",
		".jpeg": "jpeg",
		".bmp":  "bmp",
		".tiff": "tiff",
		".gif":  "gif",
	}
	_, ok := extmap[extstring]
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "不支持的文件类型",
		})
		return
	}

	b := make([]byte, fileheader.Size) //新建一个bytes[]保存文件数据
	_, err = file.Read(b)              //写入文件流
	if err != nil {
		fmt.Printf("文件读取失败: %v \n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "文件读取失败",
		})
		return
	}
	hash := fmt.Sprintf("%x", md5.Sum(b)) //通过
	var rp models.UsersFile
	has, err := adb.Ssql.Table("users_file").Where("hash=?", hash).Get(&rp)
	if err != nil {
		fmt.Printf("数据库查询失败: %v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "数据库查询失败",
		})
		return
	}
	if has {
		//判断hash在库中是否存在,如果存在直接返回
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"msg":     "上传成功!",
			"fileurl": rp.Path,
		})
		return
	}
	//根据当前时间鹾生成一个新的文件名
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	//新的文件名
	fileName := fileNameStr + extstring
	//保存上传文件
	filePath := filepath.Join("file", "/", fileName)

	data := models.UsersFile{
		Identity: util.GenerateUUID(),
		Hash:     hash,
		Name:     fileName,
		Ext:      extstring,
		Size:     filesize,
		Path:     "static/" + fileName,
	}
	_, err = adb.Ssql.Table("users_file").Insert(data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "上传数据库失败!",
		})
		return
	}

	c.SaveUploadedFile(fileheader, filePath) //文件保存

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "上传成功!",
		"fileurl": "static/" + fileName,
	})

}

package tool

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// GetRandomStr - 获取随机字符串
func GetRandomStr(length int) string {
	str := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var key string
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		key += string(str[rand.Intn(62)])
	}

	return key
}

// CheckPostForm - 检查请求中是否包含必要参数
func CheckPostForm(params []string, c *gin.Context) error {
	for _, param := range params {
		if c.PostForm(param) == "" {
			return errors.New("required parameter missing")
		}
	}

	return nil
}

// ReturnError - 返回错误信息
func ReturnError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"code": code, "msg": err.Error()})
	return
}

// ReturnRawJSON - 返回原始 JSON 信息
func ReturnRawJSON(c *gin.Context, jsonData []byte) {
	c.Render(
		http.StatusOK, render.Data{
			ContentType: "application/json; charset=utf-8",
			Data:        jsonData})
	return
}

// ReturnJSON - 返回处理后的 JSON 信息
func ReturnJSON(c *gin.Context, jsonData interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": jsonData})
	return
}

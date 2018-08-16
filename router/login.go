package router

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zhousiru/Waver/util/nmrequest"
	"github.com/Zhousiru/Waver/util/tool"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func loginByPhone(c *gin.Context) {
	if err := tool.CheckPostForm([]string{"password"}, c); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(c.PostForm("password")))
	passwordMD5 := hex.EncodeToString(hasher.Sum(nil))

	params := map[string]string{
		"password":      passwordMD5,
		"phone":         c.Param("phone"),
		"rememberLogin": "true"}

	resp, err := nmrequest.DoRequest("https://music.163.com/weapi/login/cellphone", params, nil)

	if err != nil {
		tool.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tool.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	if c.Param("type") == "/raw" {
		tool.ReturnRawJSON(c, jsonData)
	} else {
		code, err := jsonparser.GetInt(jsonData, "code")
		if err != nil {
			tool.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		if code != http.StatusOK {
			tool.ReturnError(c, http.StatusInternalServerError, fmt.Errorf("netease music api returned %d", code))
			return
		}

		avatarURL, err := jsonparser.GetString(jsonData, "profile", "avatarUrl")
		if err != nil {
			tool.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		backgroundURL, err := jsonparser.GetString(jsonData, "profile", "backgroundUrl")
		if err != nil {
			tool.ReturnError(c, http.StatusInternalServerError, err)
			return
		}

		var cookiesStr string
		cookiesArray := resp.Cookies()
		for _, cookie := range cookiesArray {
			cookiesStr += cookie.Raw
		}

		tool.ReturnJSON(c, gin.H{
			"avatarUrl":     avatarURL,
			"backgroundUrl": backgroundURL,
			"cookies":       cookiesStr})
	}

	return
}

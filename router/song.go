package router

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Zhousiru/Waver/util/nmrequest"
	"github.com/Zhousiru/Waver/util/tool"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func getSongFile(c *gin.Context) {
	brRaw := c.Param("br")
	authCookie := c.PostForm("authCookie")

	var br int
	switch brRaw {
	case "128kpbs":
		br = 128000
	case "192kpbs":
		br = 192000
	case "320kpbs":
		br = 320000
	default:
		tool.ReturnError(c, http.StatusBadRequest, errors.New("unknown bit rate"))
		return
	}

	var cookies string
	if authCookie == "" {
		cookies = fmt.Sprintf("_ntes_nuid=%s;", tool.GetRandomStr(32))
		if br > 128000 {
			tool.ReturnError(c, http.StatusBadRequest, errors.New("bit rate can't be greater than 128kpbs"))
			return
		}
	} else {
		cookies = fmt.Sprintf("_ntes_nuid=%s; MUSIC_U=%s", tool.GetRandomStr(32), authCookie)
	}

	params := map[string]string{
		"br":  strconv.Itoa(br),
		"ids": "[" + c.Param("songId") + "]"}

	resp, err := nmrequest.DoRequest("https://music.163.com/weapi/song/enhance/player/url", params, func(req *http.Request) {
		req.Header.Add("Cookie", cookies)
	})
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
		respJSON := make(map[int64]interface{})
		jsonparser.ArrayEach(jsonData, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
			code, err := jsonparser.GetInt(value, "code")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			id, err := jsonparser.GetInt(value, "id")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			if code != http.StatusOK {
				respJSON[id] = gin.H{"code": code}
			} else {
				url, err := jsonparser.GetString(value, "url")
				if err != nil {
					tool.ReturnError(c, http.StatusInternalServerError, err)
					return
				}
				size, err := jsonparser.GetInt(value, "size")
				if err != nil {
					tool.ReturnError(c, http.StatusInternalServerError, err)
					return
				}
				md5, err := jsonparser.GetString(value, "md5")
				if err != nil {
					tool.ReturnError(c, http.StatusInternalServerError, err)
					return
				}
				respJSON[id] = gin.H{
					"url":  url,
					"size": size,
					"md5":  md5}
			}
		}, "data")
		tool.ReturnJSON(c, respJSON)
	}

	return
}

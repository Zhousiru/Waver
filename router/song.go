package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Zhousiru/Waver/util/nmrequest"
	"github.com/Zhousiru/Waver/util/tool"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func getSongFile(c *gin.Context) {
	br := c.Param("br")

	params := map[string]string{
		"br":  br,
		"ids": "[" + c.Param("songId") + "]"}

	resp, err := nmrequest.DoRequest("https://music.163.com/weapi/song/enhance/player/url", params, func(req *http.Request) {
		req.Header.Add("Cookie", fmt.Sprintf("_ntes_nuid=%s;", tool.GetRandomStr(32)))
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
			if code != http.StatusOK {
				return
			}

			id, err := jsonparser.GetInt(value, "id")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
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
			actualBr, err := jsonparser.GetInt(value, "br")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}

			respJSON[id] = gin.H{
				"url":      url,
				"size":     size,
				"actualBr": actualBr,
				"md5":      md5}
		}, "data")
		tool.ReturnJSON(c, respJSON)
	}

	return
}

func getSongDetail(c *gin.Context) {
	songIDRaw := c.Param("songId")
	songIDs := strings.Split(songIDRaw, ",")
	songIDsJSON := "["
	for i, songID := range songIDs {
		songIDsJSON += `{"id":"` + songID + `"}`
		if i != len(songIDs)-1 {
			songIDsJSON += ","
		}
	}
	songIDsJSON += "]"
	songIDsArray := "[" + songIDRaw + "]"

	params := map[string]string{
		"c":  songIDsJSON,
		"id": songIDsArray}

	resp, err := nmrequest.DoRequest("https://music.163.com/weapi/v3/song/detail", params, nil)
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

		respJSON := make(map[int64]interface{})
		jsonparser.ArrayEach(jsonData, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
			id, err := jsonparser.GetInt(value, "id")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			name, err := jsonparser.GetString(value, "name")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			publishTime, err := jsonparser.GetInt(value, "publishTime")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			albumID, err := jsonparser.GetInt(value, "al", "id")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			albumName, err := jsonparser.GetString(value, "al", "name")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}
			albumPicURL, err := jsonparser.GetString(value, "al", "picUrl")
			if err != nil {
				tool.ReturnError(c, http.StatusInternalServerError, err)
				return
			}

			artistsJSON := make([]map[string]interface{}, 0)
			jsonparser.ArrayEach(value, func(artistsValue []byte, _ jsonparser.ValueType, _ int, _ error) {
				artistID, err := jsonparser.GetInt(artistsValue, "id")
				if err != nil {
					tool.ReturnError(c, http.StatusInternalServerError, err)
					return
				}
				artistName, err := jsonparser.GetString(artistsValue, "name")
				if err != nil {
					tool.ReturnError(c, http.StatusInternalServerError, err)
					return
				}
				artistJSON := make(map[string]interface{})
				artistJSON["id"] = artistID
				artistJSON["name"] = artistName

				artistsJSON = append(artistsJSON, artistJSON)
			}, "ar")

			respJSON[id] = gin.H{
				"name": name,
				"album": gin.H{
					"id":     albumID,
					"name":   albumName,
					"picUrl": albumPicURL},
				"artist":      artistsJSON,
				"publishTime": publishTime}
		}, "songs")
		tool.ReturnJSON(c, respJSON)
	}

	return
}

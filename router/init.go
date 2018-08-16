package router

import (
	"errors"
	"net/http"

	"github.com/Zhousiru/Waver/util/tool"

	"github.com/gin-gonic/gin"
)

// Init - 初始化路由
func Init(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		tool.ReturnError(c, http.StatusNotFound, errors.New("not found"))
	})

	root := r.Group("/")
	{
		root.Any("/", func(c *gin.Context) {
			c.String(http.StatusOK, "ヽ(✿ﾟ▽ﾟ)ノ Hello, world!")
		})

		root.Any("/login/phone/:phone/*type", loginByPhone)
		root.Any("/song/:songId/file/:br/*type", getSongFile)
		root.Any("/song/:songId/detail/*type", getSongDetail)
	}
}

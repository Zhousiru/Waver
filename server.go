package main

import (
	"io"
	"os"

	"github.com/Zhousiru/Waver/router"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	logFile, _ := os.Create("waver.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	r := gin.Default()
	router.Init(r)

	r.Run(":80")
}

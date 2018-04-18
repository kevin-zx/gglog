package main

import (
	"github.com/gin-gonic/gin"
	"gglog/router"
)

func main()  {
	webRouter := gin.Default()
	webRouter.LoadHTMLGlob("templates/*")
	router.LoadAll(webRouter)
	webRouter.Static("/public","./public")
	webRouter.Run("0.0.0.0:8081")

}

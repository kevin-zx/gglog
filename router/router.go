package router

import (
	"github.com/gin-gonic/gin"
	"gglog/services"
)

func LoadAll(router *gin.Engine) {
	router.GET("/", services.Index)
	router.GET("/detail/:web_site_id", services.Detail)
	router.POST("/logupload", services.LogUpload)
	router.POST("/websitestatus", services.WebSiteStatus)
	router.POST("/spider", services.Spider)
}
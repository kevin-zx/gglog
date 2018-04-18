package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"fmt"
	"github.com/kevin-zx/go-util/fileUtil"
	"os"
	"gglog/parser"
	"gglog/model"
	"gglog/core"
	"encoding/json"
)

func Index(c *gin.Context) {
	var webSites []model.WebSite
	core.WebLog.Find(&webSites)
	webSitesJsonBytes,err := json.Marshal(&webSites)
	//todo: 这里的err应该可以忽略
	if err != nil {
		panic(err)
	}
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":"123",
			"webSites":string(webSitesJsonBytes[:]),
		},
	)
}

func WebSiteStatus(c *gin.Context)  {
	domain := c.PostForm("domain")
	log.Println(domain)
	website := model.WebSite{}
	db := core.WebLog.First(&website,model.WebSite{Domain:domain})
	//log.Println(string())
	if db.Error != nil {
		c.JSON(http.StatusOK,map[string]string{"status":fmt.Sprintf("%d",0)})
	}
	if website.Status == 3 {
		c.JSON(http.StatusOK,map[string]string{"status":fmt.Sprintf("%d",website.Status),"id":string(website.ID)})
	}
	c.JSON(http.StatusOK,map[string]string{"status":fmt.Sprintf("%d",website.Status)})
}

func LogUpload(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	logFormat := c.PostForm("log_format")
	domain := c.PostForm("domain")
	log.Println(logFormat)
	log.Println(domain)
	log.Println(file.Filename)
	filePath := "./log_upload/"+file.Filename
	// Upload the file to specific dst.
	if fileUtil.CheckFileIsExist(filePath){
		os.Remove(filePath)
	}
	err := c.SaveUploadedFile(file, filePath)
	log.Printf("upload file %s domain is %s",filePath,domain)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"status":"0","error":fmt.Sprintf(err.Error())})
	}else{
		c.JSON(http.StatusOK, map[string]string{"status":"1"})
		webSite := model.WebSite{}
		db := core.WebLog.FirstOrCreate(&webSite,model.WebSite{Domain:domain})
		//把网址状态重新置位1
		webSite.Status = 1
		core.WebLog.Save(&webSite)

		if db.Error != nil {
			panic(db.Error)
		}
		go parser.NginxParse(filePath,logFormat,&webSite)
	}
}

func Spider(c *gin.Context)  {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":"123",
		},
	)
}

func Detail(c *gin.Context)  {
	webSiteId := c.Param("web_site_id")
	var webSite model.WebSite
	core.WebLog.Find(&webSite,"id = ?", webSiteId)
	wd := parser.NewWebLogDetail(webSite)
	c.HTML(
		http.StatusOK,
		"detail.html",
		gin.H{
			"title":"123",
			"webSite":webSite,
			"wdDetail":wd,
		},
	)
}
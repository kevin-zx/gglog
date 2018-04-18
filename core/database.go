package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"github.com/kevin-zx/go-util/mysqlUtil"
)

var (
	WebLog gorm.DB
	Mu mysqlutil.MysqlUtil
)

func init() {
	log.Println("WebLog 初始化")
	//webLog, err := gorm.Open("mysql", "remote:Iknowthat@@!221@tcp(182.254.158.105:3306)/site_log?charset=utf8&parseTime=True&loc=Local")
	//webLog, err := gorm.Open("mysql", "remote:Iknowthat@tcp(115.159.3.51:3306)/site_log?charset=utf8&parseTime=True&loc=Local")
	webLog, err := gorm.Open("mysql", "root:Iknowthat@tcp(localhost:3306)/site_log?charset=utf8&parseTime=True&loc=Local")

	WebLog = *webLog
	if err != nil {
		panic(err)
	}
	Mu.InitMySqlUtilByDb(WebLog.DB())

}

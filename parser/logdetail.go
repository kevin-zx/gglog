package parser

import (
	"gglog/model"
	"gglog/core"
	"strconv"
)

type WebLogDetail struct {
	RecentWeekCrawlerTimes map[string][]DayCrawlerTime `json:"recent_week_crawler_times"`
	webSite model.WebSite
	HomePageWeekCrawlerTimes map[string][]DayCrawlerTime `json:"home_page_week_crawler_times"`
}

// webLog的初始化方法
func NewWebLogDetail(webSite model.WebSite) *WebLogDetail {
	var wd WebLogDetail
	wd.webSite = webSite
	//默认会统计百度蜘蛛的抓取次数
	wd.getRecentWeekCrawlerTimesByEngine(BAIDU_SPIDER_ENGINE)

	return &wd
}

func (wd *WebLogDetail) getHomePageWeekCrawlerTimesByEngine(engine string)  {

}

//func (wd *WebLogDetail)  {
//
//}

func (wd *WebLogDetail) getRecentWeekCrawlerTimesByEngine(engine string) {
	data,err := core.Mu.SelectAll(`SELECT DATE_FORMAT(time_local,'%Y-%m-%d') date,COUNT(1) count FROM nginx_logs WHERE web_site_id = ? AND http_user_agent  LIKE "%`+engine+`%"  AND time_local < DATE_SUB(NOW(),INTERVAL 7 day)GROUP BY DATE_FORMAT(time_local,"%Y-%m-%d")`,wd.webSite.ID)
	if err!=nil {
		panic(err)
	}
	var dcs []DayCrawlerTime
	for _,row := range *data {
		date :=row["date"]
		tStr := row["count"]
		//这里的错误可以忽略
		times,_ := strconv.Atoi(tStr)
		dcs = append(dcs,DayCrawlerTime{Date:date,Times:times})

	}
	if wd.RecentWeekCrawlerTimes==nil {
		wd.RecentWeekCrawlerTimes = make(map[string][]DayCrawlerTime)
	}
	wd.RecentWeekCrawlerTimes[engine] = dcs
}



type DayCrawlerTime struct {
	Date string `json:"date"`
	Times int `json:"times"`
}

package parser

import (
	"gglog/core"
	"io"
	"os"
	"github.com/satyrius/gonx"
	"time"
	"strings"
	"gglog/model"
	"github.com/pkg/errors"
	"log"
	"bufio"
)

func NginxParse(logFile string,format string, webSite *model.WebSite)  {

	err := core.Mu.InitMySqlUtilByDb(core.WebLog.DB())
	if err != nil {
		panic(err)
	}
	logReader, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	defer logReader.Close()
	d :=bufio.NewReader(logReader)
	

	var result [][]interface{}


	// 更改状态为2
	webSite.Status = 2
	core.WebLog.Save(webSite)
	var startTime time.Time
	var uplogManger UpLogManager
	for {
		line,err := d.ReadString('\n')
		if err == io.EOF {
			break
		}
		reader := gonx.NewReader(strings.NewReader(line), format)
		rec, err := reader.Read()
		if err != nil{
			panic(err)
		}
		remoteAddr,_ := rec.Field("remote_addr")
		remoteUser,_ := rec.Field("remote_user")
		timeLocalStr,_ := rec.Field("time_local")
		timeLocalDate,_:=time.Parse("2/Jan/2006:15:04:05 -0700", timeLocalStr)
		timeLocal := timeLocalDate.Format("2006-01-02 15:04:05")
		log.Println(timeLocal)
		request,_ := rec.Field("request")
		requestParts := strings.Split(request," ")
		//第一次会初始化开始时间
		if startTime.IsZero(){
			startTime = timeLocalDate
			uplogManger = *NewUpLogManager(webSite,startTime)
		}
		if !uplogManger.NeedAnalysis(timeLocalDate){
			continue
		}
		//每次都会更新endTime
		requestMethod, requestPath, requestProtocol := "","",""
		if len(requestParts) == 3 {
			requestMethod = requestParts[0]
			requestPath = requestParts[1]
			requestProtocol = requestParts[2]
		}else if len(requestParts) == 2 {
			requestMethod = requestParts[0]
			requestPath = requestParts[1]
		}else {
			requestMethod = requestParts[0]
		}

		status,_ := rec.Field("status")
		bodyBytesSent,_ := rec.Field("body_bytes_sent")
		httpReferer,_ := rec.Field("http_referer")
		httpUserAgent,_ := rec.Field("http_user_agent")
		httpXForwardedFor,_ := rec.Field("http_x_forwarded_for")

		result = append(result, []interface{}{webSite.ID,remoteAddr, remoteUser, timeLocal, requestMethod, requestPath, requestProtocol,status, bodyBytesSent, httpReferer, httpUserAgent, httpXForwardedFor})
		if len(result) >= 10000 {
			err = core.Mu.ExecBatch("INSERT INTO nginx_logs " +
				"(`web_site_id`,`remote_addr`,`remote_user`,`time_local`,`request_method`,`request_path`,`request_protocol`,`status`,`body_bytes_sent`,`http_referer`,`http_user_agent`,`http_x_forwarded_for`,`created_at`,`deleted_at`,`updated_at`)" +
				"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW(),NOW())", result)
			if err != nil{
				panic(err)
			}

			result = [][]interface{}{}
		}
	}

	err = core.Mu.ExecBatch("INSERT INTO nginx_logs " +
		"(`web_site_id`,`remote_addr`,`remote_user`,`time_local`,`request_method`,`request_path`,`request_protocol`,`status`,`body_bytes_sent`,`http_referer`,`http_user_agent`,`http_x_forwarded_for`,`created_at`,`deleted_at`,`updated_at`)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW(),NOW())", result)
	if err != nil{
		panic(err)
	}
	// 更改状态为3
	webSite.Status = 3
	db := core.WebLog.Save(webSite)
	if db.Error != nil {
		panic(db.Error)
	}
	uplogManger.GroupDate()

}



type UpLogManager struct{
	// 这个是可能会涉及到的冲突时间段
	upLogs []model.UploadLog
	Site *model.WebSite
	currentUpLogIndex int
	initFlag bool
	fileStartDate time.Time
	fileEndDate time.Time
}

func NewUpLogManager(site *model.WebSite, fileStartDate time.Time) *UpLogManager {
	um := &UpLogManager{Site:site}
	um.fileStartDate = fileStartDate
	// 筛选出与当前文件时间段有可能重合的上传日志
	//core.WebLog.Order("start_time").Find(&um.upLogs,"(start_time > ? OR end_time >?) AND web_site_id = ?", fileStartDate, fileStartDate,site.ID)
	core.WebLog.Order("start_time").Find(&um.upLogs,"web_site_id = ?", site.ID)
	if len(um.upLogs) == 0 {
		um.currentUpLogIndex = -1
	}
	um.initFlag = true
	return um
}

// 处理一个站点上传任务的时间段聚合逻辑
func (um *UpLogManager) GroupDate() {
	// 如果没有 有可能重叠时间的上传记录则直接存储
	if len(um.upLogs) == 0 {
		ul := model.UploadLog{WebSite:*um.Site,StartTime: um.fileStartDate,EndTime: um.fileEndDate}
		Db := core.WebLog.Save(&ul)
		if Db.Error != nil {
			panic(Db.Error)
		}
		return
	}

	var deleteLogs []model.UploadLog
	//聚合操作的主要执行逻辑
	for _,uplog :=range um.upLogs {
		// 这样的情况是上传的记录是子集这种情况不进行任何存储操作
		if uplog.StartTime.Before(um.fileStartDate) &&uplog.EndTime.After(um.fileEndDate){
			return
		}
		// 如果uplog的时间记录和文件中的时间有重叠则进行的一些操作
		if (uplog.StartTime.After(um.fileStartDate) && uplog.StartTime.Before(um.fileEndDate)) || (uplog.EndTime.After(um.fileStartDate) && uplog.EndTime.Before(um.fileEndDate)){
			deleteLogs = append(deleteLogs, uplog)

			if uplog.StartTime.Before(um.fileStartDate) {
				um.fileStartDate = uplog.StartTime
			}
			if uplog.EndTime.After(um.fileEndDate) {
				um.fileEndDate = uplog.EndTime
			}

		}
	}
	//存储最新的日志时间信息
	updateUpLog :=model.UploadLog{StartTime:um.fileStartDate,EndTime:um.fileEndDate,WebSite:*um.Site}
	Db := core.WebLog.Save(&updateUpLog)
	if Db.Error != nil {
		panic(Db.Error)
	}
	// 删除重叠的时间信息
	for _,dl := range deleteLogs{
		Db = core.WebLog.Delete(&dl)
		if Db.Error != nil {
			panic(Db.Error)
		}
	}


}

// 看当前的记录是否需要被分析，库中是否有这个时间段的记录
func (um *UpLogManager) NeedAnalysis(fileCurDate time.Time) bool {
	//更新结束时间
	um.fileEndDate = fileCurDate

	// 初始化确认, golang没有构造函数确实不简单
	if !um.initFlag {
		panic(errors.New("请使用 NewDateSelector"))
	}

	// 没有限制的时间段则
	if um.currentUpLogIndex < 0{
		return true
	}
	// 小于当前上传记录的开始时间
	if um.upLogs[um.currentUpLogIndex].StartTime.After(fileCurDate){
		return true
	}

	// 文件当前时间 大于开始时间 小于结束时间的 则返回false
	if um.upLogs[um.currentUpLogIndex].StartTime.Before(fileCurDate) && um.upLogs[um.currentUpLogIndex].EndTime.After(fileCurDate) {
		return false
	}

	// 大于结束时间 这个时候过度到下一个uplogs
	if um.upLogs[um.currentUpLogIndex].EndTime.Before(fileCurDate) {
		um.currentUpLogIndex += 1
		if um.currentUpLogIndex == len(um.upLogs) {
			um.currentUpLogIndex = -1
		}
		return um.NeedAnalysis(fileCurDate)
	}
	// 等于也会返回false
	return false
}
package main

import (
	"flag"
	//"fmt"
	"io"
	"os"
	"strings"
	"github.com/satyrius/gonx"
	"github.com/kevin-zx/go-util/mysqlUtil"
	"time"
)

var conf string
var format string
var logFile string

func init() {
	flag.StringVar(&conf, "conf", "dummy", "Nginx config file (e.g. /etc/nginx/nginx.conf)")
	flag.StringVar(&format, "format", "main", "Nginx log_format name")
	flag.StringVar(&logFile, "log", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
}

func main() {
	flag.Parse()
	mu := mysqlutil.MysqlUtil{}
	mu.InitMySqlUtil("182.254.158.105",3306,"remote","Iknowthat@@!221","site_log")
	mu.InitMySqlUtil("localhost",3306,"root","Iknowthat","site_log")
	domain := "www.topqiye.com"
	defer mu.Close()
	logFile = "data/access.log"
	format = `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
	// Read given file or from STDIN
	var logReader io.Reader
	var err error
	file, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	logReader = file
	defer file.Close()

	// Use nginx config file to extract format by the name


	// Read from STDIN and use log_format to parse log records
	reader := gonx.NewReader(logReader, format)

	var result [][]interface{}
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		// Process the record... e.g.
		//fmt.Printf("Parsed entry: %+v\n", rec)
		remoteAddr,_ := rec.Field("remote_addr")
		remoteUser,_ := rec.Field("remote_user")
		timeLocalStr,_ := rec.Field("time_local")
		timeLocalDate,_:=time.Parse("2/Jan/2006:15:04:05 -0700", timeLocalStr)
		timeLocal := timeLocalDate.Format("2006-01-02 15:04:05")
		request,_ := rec.Field("request")
		requestParts := strings.Split(request," ")
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

		result = append(result, []interface{}{domain,remoteAddr, remoteUser, timeLocal, requestMethod, requestPath, requestProtocol,status, bodyBytesSent, httpReferer, httpUserAgent, httpXForwardedFor})
		if len(result) >= 500 {
			err = mu.ExecBatch("INSERT INTO nginx_log " +
				"(`domain`,`remote_addr`,`remote_user`,`time_local`,`request_method`,`request_path`,`request_protocol`,`status`,`body_bytes_sent`,`http_referer`,`http_user_agent`,`http_x_forwarded_for`,`created_at`,`deleted_at`,`updated_at`)" +
				"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW(),NOW())", result)
			if err != nil{
				panic(err)
			}

			result = [][]interface{}{}
		}
	}
	err = mu.ExecBatch("INSERT INTO nginx_log " +
		"(`domain`,`remote_addr`,`remote_user`,`time_local`,`request_method`,`request_path`,`request_protocol`,`status`,`body_bytes_sent`,`http_referer`,`http_user_agent`,`http_x_forwarded_for`,`created_at`,`deleted_at`,`updated_at`)" +
			"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW(),NOW())", result)
	if err != nil{
		panic(err)
	}

}
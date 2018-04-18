package model

import (
	"github.com/jinzhu/gorm"
	"time"
)
type WebSite struct {
	gorm.Model
	Domain string `gorm:"size:255" json:"domain"`
	Status int `json:"status"`
}

type NginxLog struct {
	gorm.Model
	WebSite WebSite `json:"web_site"`
	RemoteAddr        string    `gorm:"size:255" json:"remote_addr"`
	TimeLocal         time.Time `json:"time_local"`
	RequestMethod     string    `json:"request_method"`
	RequestPath       string    `json:"request_path"`
	RequestProtocol   string    `json:"request_protocol"`
	RemoteUser        string    `json:"remote_user"`
	Status            int       `json:"status"`
	BodyBytesSent     int       `json:"body_bytes_sent"`
	HttpReferer       string    `json:"http_referer"`
	HttpUserAgent     string    `json:"http_user_agent"`
	HttpXForwardedFor string    `json:"http_x_forwarded_for"`
}

type UploadLog struct {
	gorm.Model
	WebSite WebSite `json:"web_site"`
	WebSiteId int `json:"web_site_id"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
}

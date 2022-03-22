package com

import (
	log "github.com/sirupsen/logrus"
)

// Log teamkit 公共日志输出
var Log *log.Entry

func init() {
	// 打印文件等详细信息，这个对性能有影响，产环境不开启
	log.SetReportCaller(true)
	// 日志前缀设置
	Log = log.WithField("pkg", "nsq_exec_php")
}

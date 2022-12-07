package pkg

import (
	"time"
)

//将日期格式转换成time.Time
func TimeStr2Time(fmtStr, valueStr, locStr string) time.Time {
	loc := time.Local
	if locStr != "" {
		loc, _ = time.LoadLocation(locStr) // 设置时区
	}
	if fmtStr == "" {
		fmtStr = "2006年1月2日"
	}
	t, _ := time.ParseInLocation(fmtStr, valueStr, loc)
	return t
}

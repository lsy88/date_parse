package test

import (
	"fmt"
	"github.com/lsy88/date_parse/date"
	"testing"
)

func TestDate(t *testing.T) {
	//首先
	date.InitDate()
	//获取下一个节假日信息
	agenday, day := date.NextBigDay("2022年6月1日", date.Date_4)
	fmt.Println(agenday)
	fmt.Println(day)
	
	//新建自定义日程
	date.AddBigDay("生日", "2022年6月2日", "", "今天是我的生日", true)
	countdown, day := date.NextBigDay("2022年6月1日", "")
	fmt.Println(countdown)
	fmt.Println(day)
	//fmt.Println(date.GetBigDayList())
	//
	date.DeleteBigDay("2022年6月2日")
	countdown, day = date.NextBigDay("2022年6月1日", "")
	fmt.Println(countdown)
	fmt.Println(day)
	fmt.Println(date.GetBigDayList())
	date.AddBigDay("出去玩", "2022年4月29日", "", "出去吃饭", false)
	countdown, day = date.NextBigDay("2022年4月10日", "")
	fmt.Println(countdown)
	fmt.Println(day)
	
	fmt.Println(date.IsWeekDay("2022/12/04", ""))
	
	fmt.Println(date.FetchDay("2022/12/11"))
	
	fmt.Println(date.FetchByChName("端午节"))
	
	fmt.Println(date.IsHoliday("2022/06/04"))
}

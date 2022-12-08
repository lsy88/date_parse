package date

import (
	"errors"
	"fmt"
	"github.com/lsy88/date_parse/parse"
	"github.com/lsy88/date_parse/pkg"
	"github.com/lsy88/date_parse/translate"
	"log"
	"sort"
	"time"
)

const (
	Date_1 = "2006/01/02"
	Date_2 = "2006-01-02"
	Date_3 = "2006年01月02日"
	Date_4 = "2006年1月2日"
)

func InitDate() {
	web := parse.GetWeb()
	url, err := parse.GetURL(web)
	if err != nil {
		log.Fatalln(err)
	}
	initHoliday(url)
}

func initHoliday(url string) {
	date, err := parse.Parse(url)
	if err != nil {
		log.Fatalf("parse html failed: %v\n", err)
	}
	Agenda = buildHoliday(true, date...)
}

//重要日程(包括节假日,自定义备忘日程)
type AgendaCollection struct {
	Ch_Name   string `json:"ch_name"`
	En_Name   string `json:"en_name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Remark    string `json:"remark"` //备注
	isTrans   bool
	IsHoliday bool `json:"is_holiday"` //是否是节假日
	isDel     bool
}

//选项模式
type CollectionOption func(d *AgendaCollection)

//构造函数
func NewDateCollection(chName, enName, start, end, remark string, trans, isHoliday bool) *AgendaCollection {
	return &AgendaCollection{
		Ch_Name:   chName,
		EndTime:   end,
		En_Name:   enName,
		StartTime: start,
		Remark:    remark,
		isTrans:   trans,
		IsHoliday: isHoliday,
	}
}

//设置开始值
func WithStart(start string) CollectionOption {
	return func(d *AgendaCollection) {
		d.StartTime = start
	}
}

//设置结束值
func WithEnd(end string) CollectionOption {
	return func(d *AgendaCollection) {
		d.EndTime = end
	}
}

//设置备注内容
func WithRemark(remark string) CollectionOption {
	return func(d *AgendaCollection) {
		d.Remark = remark
	}
}

//设置是否翻译英文名称
func WithIsTrans(isTrans bool) CollectionOption {
	return func(d *AgendaCollection) {
		d.isTrans = isTrans
	}
}

//设置中文名称
func WithChName(chName string) CollectionOption {
	return func(d *AgendaCollection) {
		d.Ch_Name = chName
	}
}

type IAgendaCollection interface {
	Set(...AgendaCollection) error
	Get() []AgendaCollection
	Delete(string)
	swap()
}

//全局接收器
var Agenda AgendaDay

//节假日集合
type AgendaDay struct {
	data []AgendaCollection
}

func newAgendaDay(data []AgendaCollection) *AgendaDay {
	return &AgendaDay{
		data: data,
	}
}

func (h *AgendaDay) Set(data ...AgendaCollection) (err error) {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	h.data = append(h.data, data...)
	return
}

func (h *AgendaDay) Get() []AgendaCollection {
	agendas := make([]AgendaCollection, 0)
	for _, v := range h.data {
		if !v.isDel {
			agendas = append(agendas, v)
		}
	}
	return agendas
}

//根据日期删除自定义的日程
func (h *AgendaDay) Delete(date string) {
	for i, v := range h.data {
		if v.StartTime == date && !v.IsHoliday {
			h.data[i].isDel = true
			break
		}
	}
}

func (h *AgendaDay) swap() {
	sort.SliceStable(h.data, func(i, j int) bool {
		iTime := pkg.TimeStr2Time("", h.data[i].StartTime, "")
		jTime := pkg.TimeStr2Time("", h.data[j].StartTime, "")
		if iTime.After(jTime) {
			return false
		}
		return true
		
	})
}

func buildHoliday(isTrans bool, date ...[]string) (agenday AgendaDay) {
	var year string
	for _, v := range date {
		if len(v[2]) == 4 || len(v[2]) == 5 {
			if len(v[1]) == 15 {
				year = v[1][0:7]
				v[2] = v[1][7:11] + v[2]
				v[1] = v[1][7:]
			} else if len(v[1]) == 9 {
				if v[1][1] >= '0' && v[1][1] <= '9' {
					v[2] = v[1][0:5] + v[2]
				} else {
					v[2] = v[1][0:4] + v[2]
				}
			} else {
				v[2] = v[1][0:4] + v[2]
			}
		}
		if len(v[1]) != 15 {
			v[1] = year + v[1]
		}
		v[2] = year + v[2]
		var enName string
		if isTrans {
			enName = translate.Translate(v[0])
		}
		ag := AgendaCollection{
			Ch_Name:   v[0],
			En_Name:   enName,
			StartTime: v[1],
			EndTime:   v[2],
			isTrans:   isTrans,
			IsHoliday: true,
		}
		agenday.data = append(agenday.data, ag)
	}
	return
}

//查询某日的下一个节假日
func NextBigDay(date string, format string) (agenda AgendaCollection, day string) {
	//fmt.Println(Agenda.data)
	if format == "" {
		format = Date_4
	}
	if date == "" {
		date = time.Now().Format(Date_4)
	}
	nowTime := pkg.TimeStr2Time(format, date, "")
	//fmt.Println(nowTime)
	for i := 0; i < len(Agenda.Get()); i++ {
		b := pkg.TimeStr2Time(format, Agenda.Get()[i].StartTime, "")
		if nowTime.Before(b) {
			return Agenda.Get()[i], fmt.Sprintf("%v", b.Sub(nowTime).Hours()/24) + "天"
		}
	}
	return
}

func build(chName string, start, end string, remark string, isTrans bool) error {
	var enName string
	if isTrans {
		enName = translate.Translate(chName)
	}
	ag := AgendaCollection{
		Ch_Name:   chName,
		En_Name:   enName,
		StartTime: start,
		EndTime:   end,
		Remark:    remark,
		isTrans:   isTrans,
		IsHoliday: false,
	}
	err := Agenda.Set(ag)
	if err != nil {
		return err
	}
	Agenda.swap()
	return nil
}

//添加重要日程
func AddBigDay(chName string, start, end string, remark string, isTrans bool) error {
	return build(chName, start, end, remark, isTrans)
}

//删除自定义日程（节假日不允许删除）
func DeleteBigDay(date string) {
	Agenda.Delete(date)
}

//获取所有日程表
func GetBigDayList() []AgendaCollection {
	return Agenda.Get()
}

//判断某天是不是周日,日期格式:2006/01/02
func IsWeekDay(date string, format string) bool {
	if format == "" {
		format = Date_1
	}
	day := pkg.TimeStr2Time(format, date, "").Weekday()
	if day == time.Saturday || day == time.Sunday {
		return true
	}
	return false
}

//判断某日是不是工作日
func IsWorkDay(date string, format string) bool {
	if IsWeekDay(date, format) {
		return false
	}
	return true
}

//查询某天是周几
func FetchDay(date string) string {
	return pkg.TimeStr2Time(Date_1, date, "").Weekday().String()
}

//查询某一天是不是节假日
func IsHoliday(date string) bool {
	day := pkg.TimeStr2Time(Date_1, date, "")
	fmt.Println(day)
	for _, v := range Agenda.data {
		if v.IsHoliday { //是节假日时进入
			start := pkg.TimeStr2Time("", v.StartTime, "")
			end := pkg.TimeStr2Time("", v.EndTime, "")
			if day.Before(end) && day.After(start) {
				return true
			}
		}
	}
	return false
}

//根据中文名查询节假日
func FetchByChName(name string) interface{} {
	for _, v := range Agenda.Get() {
		if name == v.Ch_Name {
			plan := struct {
				ChName string
				EnName string
				Start  string
				End    string
				Remark string
			}{}
			plan.ChName = v.Ch_Name
			plan.EnName = v.En_Name
			plan.End = v.EndTime
			plan.Start = v.StartTime
			plan.Remark = v.Remark
			return plan
		}
	}
	return "就是普通的一天,啥事儿没有"
}

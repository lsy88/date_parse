package parse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	date1 = "[\u4e00-\u9fa5]、[\u4e00-\u9fa5]{2,3}：</span>([0-9]{4}[\u4e00-\u9fa5])([0-9]{1,2}[\u4e00-\u9fa5])([0-9]{1,2}[\u4e00-\u9fa5])[\u4e00-\u9fa5]([0-9]{0,2}[\u4e00-\u9fa5]{0,1})([0-9]{1,2}[\u4e00-\u9fa5])"
	date2 = "[\u4e00-\u9fa5]、[\u4e00-\u9fa5]{2,3}：</span>([0-9]{1,2}[\u4e00-\u9fa5])([0-9]{1,2}[\u4e00-\u9fa5])[\u4e00-\u9fa5]([0-9]{0,2}[\u4e00-\u9fa5]{0,1})([0-9]{1,2}[\u4e00-\u9fa5])"
)

// 年月日 [\u4e00-\u9fa5]、[\u4e00-\u9fa5]{2,3}：([0-9]{4}[\u4e00-\u9fa5])([0-9]{1,2}[\u4e00-\u9fa5])([0-9]{1,2}[\u4e00-\u9fa5])[\u4e00-\u9fa5]([0-9]{0,2}[\u4e00-\u9fa5]{0,1})([0-9]{1,2}[\u4e00-\u9fa5])
// 月日   [\u4e00-\u9fa5]、[\u4e00-\u9fa5]{2,3}：([0-9]{1,2}[\u4e00-\u9fa5])([0-9]{1,2}[\u4e00-\u9fa5])[\u4e00-\u9fa5]([0-9]{0,2}[\u4e00-\u9fa5]{0,1})([0-9]{1,2}[\u4e00-\u9fa5])

//从网站搜索要爬取的网页地址
func GetWeb() string {
	//获取当前年份
	year := time.Now().Year()
	web_site := fmt.Sprintf(`http://sousuo.gov.cn/s.htm?t=paper&advance=false&n=10&timetype=timezd&mintime=%d-09-01&maxtime=%d-12-31`, year-1, year-1) + `&sort=pubtime&q=%E9%83%A8%E5%88%86%E8%8A%82%E5%81%87%E6%97%A5%E5%AE%89%E6%8E%92%E7%9A%84%E9%80%9A%E7%9F%A5`
	web, err := HttpGet(web_site)
	if err != nil {
		panic(err)
	}
	return web
}

//获取网页的地址
func GetURL(web string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(web))
	if err != nil {
		return "", err
	}
	url, _ := doc.Find(`body > div.content > div > div.gov-right > div.result > ul > li > h3 > a`).Attr("href")
	return url, nil
}

// 爬取指定url页面，返回result
func HttpGet(url string) (result string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	// 设置头部信息
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = err
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	//循环爬取整页数据
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])
	}
	return
}

func Parse(url string) (date [][]string, err error) {
	get, err := HttpGet(url)
	if err != nil {
		return nil, err
	}
	//正则匹配数据
	ret := regexp.MustCompile(date1)
	file := ret.FindAllStringSubmatch(get, -1)
	for _, name := range file {
		date = append(date, []string{name[0]})
	}
	ret = regexp.MustCompile(date2)
	file = ret.FindAllStringSubmatch(get, -1)
	for _, name := range file {
		date = append(date, []string{name[0]})
	}
	for i := range date {
		//数据清洗
		line := date[i][0][6:]
		replace := strings.Replace(line, "</span>", "", -1)
		replace = strings.Replace(replace, "至", " ", -1)
		replace = strings.Replace(replace, "：", " ", -1)
		t := strings.Split(replace, " ")
		date[i] = t
	}
	return date, nil
}

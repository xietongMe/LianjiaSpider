package spider

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

//定义page结构体用来处理json
type Page struct {
	TotalPage int `json:"totalPage"`
	CurPage   int `json:"curPage"`
}

func GetSellingPageSpider(db *gorm.DB, districtName string) int {
	var totalPage int
	c := colly.NewCollector(
		//colly.Async(true),并发
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.SetRequestTimeout(time.Duration(35) * time.Second)
	c.Limit(&colly.LimitRule{DomainGlob: "https://cs.lianjia.com/ershoufang", Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//获取不同地区的总页数
	c.OnHTML(".contentBottom .house-lst-page-box", func(e *colly.HTMLElement) {
		page := Page{}
		err := json.Unmarshal([]byte(e.Attr("page-data")), &page)
		if err != nil {
			log.Fatalln(err)
		}
		totalPage = page.TotalPage
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		c.Visit("https://cs.lianjia.com/ershoufang/" + districtName)
	})
	c.Visit("https://cs.lianjia.com/ershoufang/" + districtName)
	c.Wait()
	return totalPage
}
func GetSoldPageSpider(db *gorm.DB, districtName string) int {
	var totalPage int
	c := colly.NewCollector(
		//colly.Async(true),并发
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.SetRequestTimeout(time.Duration(90) * time.Second)
	c.Limit(&colly.LimitRule{DomainGlob: "https://cs.lianjia.com/chengjiao", Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//获取不同地区的总页数
	c.OnHTML(".contentBottom .house-lst-page-box", func(e *colly.HTMLElement) {
		page := Page{}
		err := json.Unmarshal([]byte(e.Attr("page-data")), &page)
		if err != nil {
			log.Fatalln(err)
		}
		totalPage = page.TotalPage
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		c.Visit("https://cs.lianjia.com/chengjiao/" + districtName)
	})
	c.Visit("https://cs.lianjia.com/chengjiao/" + districtName)
	c.Wait()
	return totalPage
}

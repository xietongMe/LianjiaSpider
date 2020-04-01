package spider

import (
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"xietong.me/LianjiaSpider/model"
)

func GetSellingInfoSpider(db *gorm.DB, districtName string, page int) {
	//num := rand.Intn(300)
	//time.Sleep(time.Duration(num) * time.Second)
	c := colly.NewCollector(
		//colly.Async(true),并发
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "https://cs.lianjia.com/ershoufang", Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//访问所有info 访问前20页采用goroutine
	c.OnHTML(".sellListContent>li", func(e *colly.HTMLElement) {
		re, _ := regexp.Compile(`\d+`)
		houseId := e.Attr("data-lj_action_housedel_id")
		unitPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .priceInfo .unitPrice span").Eq(0).Text())))) //读取页面元素获取单价,正则匹配单价的数字
		nameRegion := e.ChildText("div.info > div.flood > div.positionInfo > a")
		name := strings.Split(nameRegion, " ")[0]
		region := strings.Split(nameRegion, " ")[1] //获取房子ID，可根据ID直接访问房子详情主页
		totalPrice, _ := strconv.Atoi(e.DOM.Find(".info .priceInfo .totalPrice span").Eq(0).Text())
		area, _ := strconv.Atoi(string(re.Find([]byte(strings.Split(e.ChildText("div.info > div.address > div.houseInfo "), " | ")[1]))))
		var wg sync.WaitGroup
		if houseId != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				model.SaveSellingInfo(houseId, name, totalPrice, unitPrice, districtName, region, area, db)
			}()
		}
		//fmt.Println(houseId, name, totalPrice, unitPrice, "District", region, area)
	})
	c.Visit("https://cs.lianjia.com/ershoufang/" + districtName + "/pg" + strconv.Itoa(page))
	c.Wait()

}

# LianjiaSpider

## 目标
我们的目标站点是链家,主要爬取[在售二手房]((https://cs.lianjia.com/ershoufang/))数据以及[已售成交](https://cs.lianjia.com/chengjiao/)数据

本次爬取在售数据的七个字段（房子id、小区名称、房子总价、房子单价、所属行政区、详细区域、面积），
已售数据的八个字段（房子id、小区名称、房子总价、房子单价、所属行政区、具体地区、交易年份、交易月份、面积）

简单分析一下目标源
- 一页共30条
- 分页规则不连续，需要手动获取最大分页数
- 每一项的数据字段排序都是规则且不变


## 开始
我们的爬取步骤如下
- 分析页面，获取对应区域总页数totalPage
- 分析页面，循环爬取所有页面的信息
- 将爬取的信息存入数据库

### 安装
项目开发使用Go版本为1.13.7，建议使用>=1.13.7版本运行软件
```
$ git clone git@github.com:xietongMe/LianjiaSpider.git
$ cd LianjiaSpider
```
修改config/application.yaml中的配置为你本地的数据库配置

创建database并将其填入配置文件中
### 运行
```
$ go run main.go
```

### 代码片段
#### 1、获取总页数-spier/gettingPage.go
由于每个区的分页总数不同，因此需要根据爬取的不同地区的页面来判断总分页数为多少
,其核心部分为
```//获取class属性为contentBottom的div标签下的class属性为house-lst-page-box的div标签
c.OnHTML(".contentBottom .house-lst-page-box", func(e *colly.HTMLElement) {
	page := Page{}
	err := json.Unmarshal([]byte(e.Attr("page-data")), &page) 
	if err != nil {
		log.Fatalln(err)
	}
	totalPage = page.TotalPage
})
```
在这里我们的爬虫框架使用的是colly，colly的使用建议阅读[colly的使用](https://juejin.im/post/5d2fb6845188252a7b1d8a32)

其中.contentBottom .house-lst-page-box是对class属性为这两个值的div标签匹配，而e.Attr("page-data")则是获取匹配到的标签的属性值，由于其属性值为json，因此对其进行解析。

#### 2、获取二手房在售信息-spier/sellingSpider.go
```
c.OnHTML(".sellListContent>li", func(e *colly.HTMLElement) {
	re, _ := regexp.Compile(`\d+`)                                           //正则表达式用来匹配数字
	houseId := e.Attr("data-lj_action_housedel_id")                          //获取房子ID，可根据ID直接访问房子详情主页
	nameRegion := e.ChildText("div.info > div.flood > div.positionInfo > a") //同时获取小区名和详细地区
	name := strings.Split(nameRegion, " ")[0]                                //将同时获取的小区名和详细地区分离，取其中的小区名字
	region := strings.Split(nameRegion, " ")[1]                              //将同时获取的小区名和详细地区分离，取其中的详细地区
	totalPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .priceInfo .totalPrice span").Eq(0).Text()))))              //根据页面元素获取总价，正则匹配数字，转换成int类型
	unitPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .priceInfo .unitPrice span").Eq(0).Text()))))                //读取页面元素获取单价,正则匹配单价的数字，转换成int类型
	area, _ := strconv.Atoi(string(re.Find([]byte(strings.Split(e.ChildText("div.info > div.address > div.houseInfo "), " | ")[1])))) // //读取页面元素获取面积,正则匹配单价的数字，转换成int类型
	if houseId != "" {
		fmt.Println("start save", houseId, page)
		sellingInfo := model.Selling{Id: houseId, Name: name, TotalPrice: totalPrice, UnitPrice: unitPrice, District: districtName, Region: region, Area: area}
		err := db.Save(&sellingInfo).Error
		for ; err != nil; {
			sellingInfo := model.Selling{Id: houseId, Name: name, TotalPrice: totalPrice, UnitPrice: unitPrice, District: districtName, Region: region, Area: area}
			err = db.Save(&sellingInfo).Error
		}
	}
})
```
这里主要还是根据页面结构，来匹配不同的字段，可以通过访问[链家在售信息](https://cs.lianjia.com/ershoufang/)在控制台查看不同属性值对应的标签

#### 3、获取成交信息-spier/soldSpider.go
```
	c.OnHTML(".listContent>li", func(e *colly.HTMLElement) {
		re, _ := regexp.Compile(`\d+`)                                                                             //正则表达式用来匹配数字
		houseId := string(re.Find([]byte(strings.Split(e.ChildAttr("div.info > div.title > a", "href"), "/")[4]))) //获取房子ID，可根据ID直接访问房子详情主页
		name := strings.Split(e.ChildText("div.info > div.title > a"), " ")[0]                                     //获取小区名
		area := 0
		if len(strings.Split(e.ChildText("div.info > div.title > a"), " ")) == 3 {
			area, _ = strconv.Atoi(string(re.Find([]byte(strings.Split(e.ChildText("div.info > div.title > a"), " ")[2])))) //获取总面积
		}
		totalPrice, _ := strconv.Atoi(e.DOM.Find(".info .address .totalPrice span").Eq(0).Text())                      //获取总价
		unitPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .flood .unitPrice span").Eq(0).Text())))) //获取单价
		dealDate := e.DOM.Find(".info .address .dealDate").Eq(0).Text()                                                //获取成交年月日
		soldYear := strings.Split(dealDate, ".")[0]                                                                    //分离出成交年份
		soldMonth := strings.Split(dealDate, ".")[1]                                                                   //分离出成交月
		if houseId != "" {
			fmt.Println("start save", houseId, page)
			soldInfo := model.Sold{Id: houseId, Name: name, TotalPrice: totalPrice, UnitPrice: unitPrice, District: districtName, SoldYear: soldYear, SoldMonth: soldMonth, Area: area}
			err := db.Save(&soldInfo).Error
			for ; err != nil; {
				soldInfo := model.Sold{Id: houseId, Name: name, TotalPrice: totalPrice, UnitPrice: unitPrice, District: districtName, SoldYear: soldYear, SoldMonth: soldMonth, Area: area}
				err = db.Save(&soldInfo).Error
			}
		}
	})
```
这里也是根据成交页面结构，来匹配不同的字段，可以通过访问[链家成交信息](https://cs.lianjia.com/chengjiao/)在控制台查看不同属性值对应的标签
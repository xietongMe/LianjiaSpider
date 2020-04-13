package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
	"xietong.me/LianjiaSpider/common"
	"xietong.me/LianjiaSpider/spider"
)

func main() {
	//初始化配置
	InitConfig()
	db := common.InitDB()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(20)
	defer db.Close()
	// "yuhua","yuelu", "tianxin", "kaifu", "furong", "wangcheng", "ningxiang", "liuyang", "changshaxian"
	district := [9]string{"yuhua", "yuelu", "tianxin", "kaifu", "furong", "wangcheng", "ningxiang", "liuyang", "changshaxian"}
	var wgSelling sync.WaitGroup
	var wgSold sync.WaitGroup
	for _, districtName := range district {
		totalSellingPage := spider.GetSellingPageSpider(db, districtName)
		for page := 1; page < totalSellingPage; page++ {
			wgSelling.Add(1)
			time.Sleep(time.Duration(page) * time.Millisecond)
			go func(page int) {
				fmt.Println("start spider", page)
				defer wgSelling.Done()
				spider.GetSellingInfoSpider(db, districtName, page)
			}(page)
		}
	}
	wgSelling.Wait()

	for _, districtName := range district {
		totalSoldPage := spider.GetSoldPageSpider(db, districtName)
		for page := 1; page < totalSoldPage; page++ {
			wgSold.Add(1)
			time.Sleep(time.Duration(page*20) * time.Millisecond)
			go func(page int) {
				fmt.Println("start spider", page)
				defer wgSold.Done()
				spider.GetSoldInfoSpider(db, districtName, page)
			}(page)
		}
	}
	wgSold.Wait()
}

//初始化配置函数
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

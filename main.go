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
	defer db.Close()

	// "yuhua","yuelu", "tianxin", "kaifu", "furong", "wangcheng", "ningxiang", "liuyang", "changshaxian"
	district := [1]string{"wangcheng"}
	var wg sync.WaitGroup
	for _, districtName := range district {
		for page := 1; page < 27; page++ {
			wg.Add(1)
			go func(page int) {
				fmt.Println("start spider", page)
				defer wg.Done()
				time.Sleep(time.Duration(page) * time.Millisecond)
				spider.GetSellingInfoSpider(db, districtName, page)
			}(page)
		}
	}
	wg.Wait()

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

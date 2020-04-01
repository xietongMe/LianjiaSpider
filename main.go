package main

import (
	"github.com/spf13/viper"
	"os"
	"sync"
	"xietong.me/LianjiaSpider/common"
	"xietong.me/LianjiaSpider/spider"
)

func main() {
	//初始化配置
	InitConfig()
	db := common.InitDB()
	defer db.Close()
	// "yuelu", "tianxin", "kaifu", "furong", "wangcheng", "ningxiang", "liuyang", "changshaxian"
	district := [1]string{"yuhua"}
	var wg sync.WaitGroup
	for _, districtName := range district {
		for page := 1; page < 3; page++ {
			wg.Add(1)
			go func() {
				//num := rand.Intn(200)
				//time.Sleep(time.Duration(num) * time.Second)
				defer wg.Done()
				spider.GetSellingInfoSpider(db, districtName, page)
			}()
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

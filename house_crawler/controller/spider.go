/*
   @Author: StudentCWZ
   @Description:
   @File: spider.go
   @Software: GoLand
   @Project: spider
   @Date: 2022/3/19 17:40
*/

package controller

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"house_crawler/logic"
	"house_crawler/models"
	"strconv"
	"time"
)

func InitSpider() (err error) {
	// 初始化 Result 指针
	r := new(models.Result)
	// 1. 获取城市名称参数
	cityMap := InitCityMap()
	// 2. 遍历城市 map
	for _, city := range cityMap {
		// url 拼接
		initialUrl := fmt.Sprintf("https://%v.ihk.cn/newhouse/houselist/", city)
		// 3. 请求索引页
		html, err := logic.ScrapePage(initialUrl)
		if err != nil {
			zap.L().Error("Request index page failed", zap.Error(err))
			// 跳入下次循环
			continue
		}
		// 4. 获取最大翻页
		maxPage, err := logic.ScrapePageNumber(html)
		if err != nil {
			zap.L().Error("Get max page failed", zap.Error(err))
			// 跳入下次循环
			continue
		}
		fmt.Printf("The max num of pages: %v\n", maxPage)
		// 5. 遍历页数
		for i := 1; i <= maxPage; i++ {
			url := initialUrl + "p" + strconv.Itoa(i)
			// 6. 请求详情页
			body, err := logic.ScrapePage(url)
			if err != nil {
				zap.L().Error("Request detail page failed", zap.Error(err))
				// 跳入下次循环
				continue
			}
			// 休眠
			time.Sleep(15 * time.Second)
			// 7. 解析详情页
			logic.ParseDetailPage(body, r)
		}
	}
	// 条件判断
	if len(r.HouseList) == 0 {
		err = errors.New("spider failed... ")
	}
	fmt.Println("--------------------")
	fmt.Println(len(r.HouseList))
	fmt.Println("--------------------")
	fmt.Println(r.HouseList)
	return
}

/*
@Time : 2022/3/20 01:16
@Author : StudentCWZ
@File : spider
@Software: GoLand
*/

package controller

import (
	"fmt"
	"job_crawler/logic"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func InitSpider() (err error) {
	// 1. 获取区域名称参数
	areaSlice := InitCitySlice()
	// 2. 获取职位名称参数
	keySlice := InitKeySlice()
	// 3. 获取 areaMap
	areaMap := logic.GetAreaCode()
	// 4. 遍历城市名称切片
	for _, area := range areaSlice {
		// 获取 areaCode
		areaCode := areaMap[area]
		for _, key := range keySlice {
			for i := 1; i < 101; i++ {
				url := "https://search.51job.com/list/" + areaCode + ",000000,0000,00,9,99," + key + ",2," +
					strconv.Itoa(i) + ".html"
				html, err := logic.ScrapePage(url, "gbk", "index")
				time.Sleep(15 * time.Second)
				if html == "" {
					zap.L().Debug("页面返回空字符串")
					continue
				}
				if err != nil {
					fmt.Println(err)
					continue
				}
				logic.ParseDetailPage(html)
				fmt.Println("-----------------------------")
			}
		}
	}
	return err
}

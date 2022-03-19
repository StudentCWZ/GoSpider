/*
   @Author: StudentCWZ
   @Description:
   @File: city.go
   @Software: GoLand
   @Project: spider
   @Date: 2022/3/19 17:54
*/

package controller

// InitCityMap 初始化 city 序列
func InitCityMap() map[string]interface{} {
	cityMap := map[string]interface{}{
		"广州": "gz",
		"佛山": "fs",
		"东莞": "dg",
		"江西": "jx",
		"安徽": "ah",
		"天津": "tj",
	}
	return cityMap
}

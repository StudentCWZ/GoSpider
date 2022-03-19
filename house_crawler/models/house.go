/*
   @Author: StudentCWZ
   @Description:
   @File: house.go
   @Software: GoLand
   @Project: spider
   @Date: 2022/3/19 17:20
*/

package models

type Result struct {
	HouseList []string
}

type House struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Price   string `json:"price"`
	Address string `json:"address"`
	Dynamic string `json:"dynamic"`
	Feature string `json:"feature"`
}

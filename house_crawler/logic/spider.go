/*
   @Author: StudentCWZ
   @Description:
   @File: spider.go
   @Software: GoLand
   @Project: spider
   @Date: 2022/3/19 18:09
*/

package logic

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"house_crawler/dao/mysql"
	"house_crawler/header"
	"house_crawler/models"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ScrapePage(url string) (html string, err error) {
	// 输出日志信息
	fmt.Println("Fetch url:", url)
	// 实例化 client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("Create request object failed", zap.Error(err))
		return "", err
	}
	req.Header.Set("User-Agent", header.Init()[rand.Intn(len(header.Init()))])
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("Request web page failed", zap.Error(err))
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.L().Error("Close response body failed", zap.Error(err))
			return
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("Read response body failed", zap.Error(err))
		return "", err
	}
	return string(body), nil
}

func ScrapePageNumber(h string) (int, error) {
	body := strings.Replace(h, "\n", "", -1)
	PageRE := regexp.MustCompile(`<a href="javascript:SetCondition\('p(.*?)'`)
	PageREItems := PageRE.FindAllStringSubmatch(body, -1)
	maxPage, err := strconv.Atoi(PageREItems[len(PageREItems)-1][1])
	if err != nil {
		zap.L().Error("Str convert int failed", zap.Error(err))
		return 0, err
	}
	return maxPage, nil
}

func ParseDetailPage(h string, r *models.Result) {
	body := strings.Replace(h, "\n", "", -1)
	houseRe := regexp.MustCompile(`<strong id=".*?">(.*?)</strong>`)
	rp := regexp.MustCompile(`(<div class="n_conlist" onclick="javaScript:listPageClick.*?)arryMetroInfo = getNewArray`)
	houseTagRe := regexp.MustCompile(`<div class="n_conlistrbrief">.*?<span>(.*?)</span>.*?</div>`)
	housePriceRe := regexp.MustCompile(`<div class="n_conliprice">.*?<strong>(.*?)</strong>.*?</div>`)
	houseAddRe := regexp.MustCompile(`<div class="n_conlistradd">.*?<span>(.*?)</span>.*?</div>`)
	houseDynamicRe := regexp.MustCompile(`<b>楼盘动态：</b>(.*?)</span>`)
	houseFeatureRe := regexp.MustCompile(`var array = arraySort\(arrayStringToChar\("(.*?)",`)
	rpItems := rp.FindAllStringSubmatch(body, -1)

	for _, item := range rpItems {
		houseName := strings.Replace(houseRe.FindStringSubmatch(item[1])[1], " ", "", -1)
		houseTag := strings.Replace(houseTagRe.FindStringSubmatch(item[1])[1], " ", "", -1)
		houseTag = strings.Replace(houseTag, "#183", "", -1)
		houseTag = strings.Replace(houseTag, "&", "", -1)
		houseTag = strings.Replace(houseTag, "#", "", -1)
		houseTag = strings.Replace(houseTag, "\\u0026", "", -1)
		housePrice := housePriceRe.FindStringSubmatch(item[1])[1]
		housePriceInt, _ := strconv.Atoi(housePrice)
		if housePriceInt > 0 && housePriceInt <= 1000 {
			housePrice = housePrice + "万"
		} else if housePriceInt > 1000 {
			housePrice = housePrice + "元/㎡"
		}
		fmt.Printf("house: %v\n", houseName)
		fmt.Printf("houseTag: %v\n", houseTag)
		fmt.Printf("housePrice: %v\n", housePrice)
		houseAdd := strings.Replace(houseAddRe.FindStringSubmatch(item[1])[1], "&#183", "", -1)
		houseAdd = strings.Replace(houseAdd, " ", "", -1)
		fmt.Printf("houseAdd: %v\n", houseAdd)
		var houseDynamic string
		if len(houseDynamicRe.FindStringSubmatch(item[1])) > 1 {
			houseDynamic = strings.Replace(houseDynamicRe.FindStringSubmatch(item[1])[1],
				"#183", "", -1)
			houseDynamic = strings.Replace(houseDynamic, "&", "", -1)
			houseDynamic = strings.Replace(houseDynamic, "#", "", -1)
			houseDynamic = strings.Replace(houseDynamic, "\\u0026", "", -1)
			houseDynamic = strings.Replace(houseDynamic, "quot", "", -1)
			houseDynamic = strings.Replace(houseDynamic, " ", "", -1)
			fmt.Printf("houseDynamic: %v\n", houseDynamic)
		} else {
			houseDynamic = ""
			fmt.Printf("houseDynamic: %v\n", houseDynamic)
		}
		houseFeature := houseFeatureRe.FindStringSubmatch(item[1])[1]
		fmt.Printf("houseFeature: %v\n", houseFeature)
		house := models.House{
			Name:    houseName,
			Tag:     houseTag,
			Price:   housePrice,
			Address: houseAdd,
			Dynamic: houseDynamic,
			Feature: houseFeature,
		}
		mysql.DB.Create(house)
		buf, err := json.Marshal(house)
		if err != nil {
			fmt.Printf("Stuct convert json failed, err: %v\n", err)
			continue
		}
		r.HouseList = append(r.HouseList, string(buf))
		fmt.Printf("%v\n", string(buf))
		fmt.Println("-------------------------")
	}
}

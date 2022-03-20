/*
@Time : 2022/3/20 01:38
@Author : StudentCWZ
@File : spider
@Software: GoLand
*/

package logic

import (
	"fmt"
	"io"
	"io/ioutil"
	"job_crawler/header"
	"job_crawler/settings"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"go.uber.org/zap"
)

func ScrapePage(url string, encode string, size string) (html string, err error) {
	// 输出日志信息
	fmt.Println("Fetch url:", url)
	// 实例化 client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("Create request object failed", zap.Error(err))
		return "", err
	}
	switch size {
	case "index":
		req.Header.Set("User-Agent", header.Init()[rand.Intn(len(header.Init()))])
		req.Header.Set("Host", "search.51job.com")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
	default:
		req.Header.Set("User-Agent", header.Init()[rand.Intn(len(header.Init()))])
	}
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
	switch encode {
	case "gbk":
		html = mahonia.NewDecoder("gbk").ConvertString(string(body))
	default:
		html = string(body)
	}
	return html, nil
}

func GetAreaCode() (areaMap map[string]string) {
	areaMap = make(map[string]string)
	html, err := ScrapePage(settings.Conf.AreaUrl, "gbk", "")
	time.Sleep(5 * time.Second)
	if err != nil {
		zap.L().Error("Requests area code page failed", zap.Error(err))
		return
	}
	body := strings.Replace(html, "\n", "", -1)
	rp := regexp.MustCompile(`(".*?":".*?")`)
	rpItems := rp.FindAllStringSubmatch(body, -1)
	areaCodeRe := regexp.MustCompile(`"(.*?)":".*?"`)
	areaNameRe := regexp.MustCompile(`".*?":"(.*?)"`)
	for _, item := range rpItems {
		areaCode := areaCodeRe.FindStringSubmatch(item[1])[1]
		areaName := areaNameRe.FindStringSubmatch(item[1])[1]
		areaMap[areaName] = areaCode
	}
	return areaMap
}

func GetPageNumber(url string) {
	html, err := ScrapePage(url, "gbk", "")
	if err != nil {
		zap.L().Error("Requests index page failed", zap.Error(err))
	}
	body := strings.Replace(html, "\n", "", -1)
	maxPageRe := regexp.MustCompile(`"total_page":"(.*?)",`)
	maxPage := maxPageRe.FindStringSubmatch(body)
	fmt.Println(maxPage)
}

func ParseDetailPage(html string) {
	body := strings.Replace(html, "\n", "", -1)
	dataRe := regexp.MustCompile(`window.__SEARCH_RESULT__ = ({.*?})</script>`)
	dataSlice := dataRe.FindStringSubmatch(body)
	fmt.Println(dataSlice)
}

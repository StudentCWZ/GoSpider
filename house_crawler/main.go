/*
   @Author: StudentCWZ
   @Description:
   @File: main.go
   @Software: GoLand
   @Project: house_crawler
   @Date: 2022/3/16 22:37
*/

package main

import (
	"fmt"
	"go.uber.org/zap"
	"house_crawler/controller"
	"house_crawler/dao/mysql"
	"house_crawler/logger"
	"house_crawler/models"
	"house_crawler/settings"
)

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("Init settings failed, err: %v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("Init logger failed, err: %v\n", err)
		return
	}
	// 日志初始化成功
	zap.L().Debug("logger init success")
	defer func(l *zap.Logger) {
		err := l.Sync()
		if err != nil {
			return
		}
	}(zap.L())
	// 3. 初始化 MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Error("Init mysql failed", zap.Error(err))
		return
	}
	defer mysql.Close()
	// 禁用复数表名
	mysql.DB.SingularTable(true)
	// 绑定模型
	mysql.DB.AutoMigrate(&models.House{})
	// 4. 爬取数据
	err := controller.InitSpider()
	if err != nil {
		zap.L().Error("House data crawler failed", zap.Error(err))
		return
	}
}

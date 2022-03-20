/*
@Time : 2022/3/20 00:16
@Author : StudentCWZ
@File : main
@Software: GoLand
*/

package main

import (
	"fmt"
	"job_crawler/controller"
	"job_crawler/dao/mysql"
	"job_crawler/logger"
	"job_crawler/models"
	"job_crawler/settings"

	"go.uber.org/zap"
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
	mysql.DB.AutoMigrate(&models.Job{})
	// 4. 爬取数据
	err := controller.InitSpider()
	if err != nil {
		zap.L().Error("House data crawler failed", zap.Error(err))
		return
	}
}

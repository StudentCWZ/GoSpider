/*
   @Author: StudentCWZ
   @Description:
   @File: mysql.go
   @Software: GoLand
   @Project: house_crawler
   @Date: 2022/3/16 22:38
*/

package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"house_crawler/settings"
)

var DB *gorm.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Connect mysql database failed, err: %v\n", err)
		return
	}
	return DB.DB().Ping()
}

func Close() {
	err := DB.Close()
	if err != nil {
		zap.L().Error("Closed mysql database failed", zap.Error(err))
		return
	}
	return
}

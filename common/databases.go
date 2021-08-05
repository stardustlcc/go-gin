package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"godev/model"
	"net/url"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	loc := viper.GetString("datasource.loc")
	db, err := gorm.Open(driverName, "root:@tcp(localhost:3306)/lcc?charset=utf8&parseTime=true&loc="+url.QueryEscape(loc))
	if err != nil {
		fmt.Println("err=", err)
	}
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDb() *gorm.DB  {
	return DB
}

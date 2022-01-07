package connect

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	nameDB := viper.GetString("mysql.name")
	host := viper.GetString("mysql.host")
	post := viper.GetString("mysql.post")
	password := viper.GetString("mysql.password")
	dsn := "root:" + password + "@tcp(" + host + ":" + post + ")/" + nameDB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db
}

package cache

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

func InitMysqlCon() (err error) {
	DB, err = gorm.Open("mysql", ":@(:3306)/?charset=utf8&parseTime=false&loc=Local")
	if err != nil {
		return
	}
	err = DB.DB().Ping()

	DB.DB().SetMaxIdleConns(20)

	DB.DB().SetMaxOpenConns(200)

	DB.DB().SetConnMaxLifetime(30 * time.Second)

	//禁用复数表名
	DB.SingularTable(true)
	return err
}

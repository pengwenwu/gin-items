package model

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"gin-items/lib/setting"
)

var db *gorm.DB

type Model struct {
	ID uint `gorm:"primary_key;"`
	LastDated string `json:"last_dated"`
	Dated string `json:"dated"`
}

func init() {
	var (
		err error
		dbType, dbName, user, password, host, tablePrefix string
	)
	dbType = setting.Config().DB.Type
	dbName = setting.Config().DB.Name
	user = setting.Config().DB.User
	password = setting.Config().DB.PassWord
	host = setting.Config().DB.Host
	tablePrefix = setting.Config().DB.TablePrefix

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}

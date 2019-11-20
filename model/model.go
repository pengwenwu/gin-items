package model

import (
	"database/sql"
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

// 结果集转切片
func Rows2SliceMap(rows *sql.Rows) (list []map[string]string) {
	//字段名称
	columns, _ := rows.Columns()
	//多少个字段
	length := len(columns)
	//每一行字段的值
	values := make([]sql.RawBytes, length)
	//保存的是values的内存地址
	pointer := make([]interface{}, length)
	//
	for i := 0; i < length; i++ {
		pointer[i] = &values[i]
	}
	//
	for rows.Next() {
		//把参数展开，把每一行的值存到指定的内存地址去，循环覆盖，values也就跟着被赋值了
		rows.Scan(pointer...)
		//每一行
		row := make(map[string]string)
		for i := 0; i < length; i++ {
			row[columns[i]] = string(values[i])
		}
		list = append(list, row)
	}
	//
	return
}

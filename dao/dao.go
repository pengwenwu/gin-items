package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"gin-items/library/setting"
)

const updateCommonLimit = 1000

type Dao struct {
	MasterServiceItems *gorm.DB
	SlaveServiceItems *gorm.DB
}

func New() (d *Dao) {
	d = &Dao{
	}
	d.init()
	return
}

func (dao *Dao) init() {
	// 防止清空表，暂时注释
	//execMigration(setting.Config().DB.Master.ServiceItems)

	dao.MasterServiceItems = openDB(setting.Config().DB.Master.ServiceItems)
	dao.SlaveServiceItems = openDB(setting.Config().DB.Slave.ServiceItems)
}

func openDB (conf *setting.Database) *gorm.DB {
	DBLink, err := gorm.Open(conf.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		conf.User,
		conf.PassWord,
		conf.Host,
		conf.Name))
	if err != nil {
		panic(err)
	}
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return conf.TablePrefix + defaultTableName
	}
	DBLink.SingularTable(true)
	if conf.NeedConnectionPool {
		DBLink.DB().SetMaxIdleConns(conf.MaxIdleConnections)
		DBLink.DB().SetMaxOpenConns(conf.MaxOpenConnections)
	}

	return DBLink
}

func execMigration(conf *setting.Database)  {
	db, _ := sql.Open(conf.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&multiStatements=true",
		conf.User,
		conf.PassWord,
		conf.Host,
		conf.Name))
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	m.Steps(2)
}

func (dao *Dao) CloseDB() {
	defer dao.MasterServiceItems.Close()
	defer dao.SlaveServiceItems.Close()
}

// 结果集转切片
func RowsToSliceMap(rows *sql.Rows) (list []map[string]string) {
	//字段名称
	columns, _ := rows.Columns()
	//多少个字段
	length := len(columns)
	//每一行字段的值
	values := make([]sql.RawBytes, length)
	//保存的是values的内存地址
	pointer := make([]interface{}, length)
	for i := 0; i < length; i++ {
		pointer[i] = &values[i]
	}
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
	return
}

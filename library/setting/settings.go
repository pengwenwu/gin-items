package setting

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

type appConfig struct {
	RunMode string `toml:"run_mode"`
	APP     appInfo
	Server  serverInfo
	DB      database `toml:"database"`
	Log		logInfo
}

type appInfo struct {
	Page      int    `toml:"page"`
	PageSize  int    `toml:"page_size"`
	JwtSecret string `toml:"jwt_secret"`
}

type serverInfo struct {
	HttpPort     int `toml:"http_port"`
	ReadTimeout  int `toml:"read_timeout"`
	WriteTimeout int `toml:"write_timeout"`
}

type database struct {
	Type        string
	User        string
	PassWord    string
	Host        string
	Name        string
	TablePrefix string `toml:"table_prefix"`
}

type logInfo struct {
	LogFilePath string `toml:"log_file_path"`
	LogFileName string `toml:"log_file_name"`
}

var (
	cfg  *appConfig
	once sync.Once
)

func init() {
}

// 单例模式
// todo: 动态更新
func Config() *appConfig {
	once.Do(func() {
		filePath, err := filepath.Abs("conf/app.toml")
		if err != nil {
			fmt.Printf("file cannot find: %s\n", filePath)
			panic(err)
		}
		fmt.Printf("parse toml file once. filePath: %s\n", filePath)
		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			log.Fatalf("Fail to parse 'conf/app.toml': %v", err)
			panic(err)
		}
	})
	return cfg
}
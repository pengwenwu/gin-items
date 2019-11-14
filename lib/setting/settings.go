package setting

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
	"sync"
	"time"
)

type AppConfig struct {
	RunMode string `toml:"run_mode"`
	APP     appInfo
	Server  serverInfo
	DB      database `toml:"database"`
}

type appInfo struct {
	PageSize  int `toml:"page_size"`
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

var (
	cfg  *AppConfig
	once sync.Once

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string
)

func init() {
	LoadBase()
	LoadServer()
	LoadApp()
}

// 单例模式
// todo: 动态更新
func Config() *AppConfig {
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

func LoadBase() {
	RunMode = Config().RunMode
}

func LoadServer() {
	HTTPPort = Config().Server.HttpPort
	ReadTimeout = time.Duration(Config().Server.ReadTimeout) * time.Second
	WriteTimeout = time.Duration(Config().Server.WriteTimeout) * time.Second
}

func LoadApp() {
	PageSize = Config().APP.PageSize
	JwtSecret = Config().APP.JwtSecret
}

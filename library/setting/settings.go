package setting

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type appConfig struct {
	RunMode  string   `mapstructure:"run_mode"`
	Server   server   `mapstructure:"server"`
	DB       multiDB  `mapstructure:"database"`
	RabbitMq rabbitMq `mapstructure:"rabbitmq"`
	Log      multiLog `mapstructure:"log"`
}

type server struct {
	HttpPort     int `mapstructure:"http_port"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

type multiDB struct {
	Master multiMasterDB `mapstructure:"master"`
	Slave  multiSlaveDB  `mapstructure:"slave"`
}

type multiMasterDB struct {
	ServiceItems *Database `mapstructure:"service_items"`
	//ShopTrades *Database `toml:"shop_trades"`
}

type multiSlaveDB struct {
	ServiceItems *Database `mapstructure:"service_items"`
}

type Database struct {
	Type               string `mapstructure:"type"`
	User               string `mapstructure:"user"`
	PassWord           string `mapstructure:"password"`
	Host               string `mapstructure:"host"`
	Name               string `mapstructure:"name"`
	TablePrefix        string `mapstructure:"table_prefix"`
	NeedConnectionPool bool   `mapstructure:"need_connection_pool"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
}

type rabbitMq struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`
}

type multiLog struct {
	AccessLog accessLog `mapstructure:"access_log"`
	ErrorLog  errorLog  `mapstructure:"error_log"`
}

type accessLog struct {
	LogConfig `mapstructure:"log_config"`
}

type errorLog struct {
	LogConfig `mapstructure:"log_config"`
}

type LogConfig struct {
	FilePath string `mapstructure:"file_path"`
	FileName string `mapstructure:"file_name"`
}

var (
	cfg  *appConfig
	once sync.Once
)

func init() {
}

// 单例模式
// 热更新
func Config() *appConfig {
	once.Do(func() {
		var confName string
		confName = os.Getenv("APP_ENV")
		if confName == "" {
			confName = "development"
		}
		viper.SetConfigName(fmt.Sprintf("config.%s", confName)) // 指定配置文件名称（不需要带后缀）
		viper.SetConfigType("toml")                             // 指定配置文件类型
		viper.AddConfigPath("./conf/")                          // 指定查找配置文件的路径（这里使用相对路径）
		err := viper.ReadInConfig()                             // 读取配置信息
		if err != nil { // 读取配置信息失败
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		err = viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			// 配置文件发生变更之后会调用的回调函数
			fmt.Println("Config file changed:", e.Name)

			err = viper.Unmarshal(&cfg)
			if err != nil {
				log.Fatalf("unable to decode into struct, %v", err)
			}
		})
	})

	return cfg
}

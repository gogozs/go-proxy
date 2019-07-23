package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

type Config struct {
	Server ServerConfig
	Proxy  []ProxyConfig
}

type ServerConfig struct {
	Port string
}

type ProxyConfig struct {
	Name      string `mapstructure:"name"`
	Location  string `mapstructure:"location"`
	ProxyPass string `mapstructure:"proxy_pass"`
}

var config Config

func GetConfig() *Config {
	return &config
}

func GetConfigPath() string {
	wd := os.Getenv("AGENT_WORKDIR")
	confPath := path.Join(wd, "conf/")
	return confPath
}

func init() {
	// 需要配置项目根目录的环境变量，方便执行test
	confPath := GetConfigPath()
	ginEnv := "conf"
	viper.SetConfigName(ginEnv)   // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(confPath) // 第一个搜索路径
	viper.WatchConfig()           // 监控配置文件热重载
	err := viper.ReadInConfig()   // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

type Config struct {
	Var VarConfig
	Server ServerConfig
	Common CommonConfig
}

type VarConfig map[string]string

type ServerConfig struct {
	Port   string
	Proxy  []ProxyConfig
	Static []StaticConfig
}

type ProxyConfig struct {
	Name      string `mapstructure:"name"`
	Location  string `mapstructure:"location"`
	ProxyPass string `mapstructure:"proxy_pass"`
}

type StaticConfig struct {
	Path  string `mapstructure:"path"`
	Alias string `mapstructure:"alias"`
}

type CommonConfig struct {
	LogFile string
	Level string
}

var config Config
var ginEnv = "proxy"

func GetConfig() *Config {
	return &config
}

func GetConfigPath() string {
	wd := os.Getenv("PROXY_WORKDIR")
	confPath := path.Join(wd, "conf/")
	return confPath
}

// 参数解析
func updateConfig(c *Config) {
	vars := c.Var
	proxyList := &(c.Server.Proxy)
	for i := range *proxyList {
		p := &(*proxyList)[i]
		for k, v := range vars {
			k := "$" + k
			if strings.Contains(p.ProxyPass, k) {
				p.ProxyPass = strings.Replace(p.ProxyPass, k, v, 1)
			}
		}
	}
}

func init() {
	// 需要配置项目根目录的环境变量，方便执行test
	confPath := GetConfigPath()
	viper.SetConfigName(ginEnv)   // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(confPath) // 第一个搜索路径
	viper.WatchConfig()           // 监控配置文件热重载
	err := viper.ReadInConfig()   // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	updateConfig(&config)
	fmt.Println(config)
	if err != nil {
		panic(err)
	}
}

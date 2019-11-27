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
	Tls bool
	LogFile string
	Level string
	ReadTimeOut int
	WriteTimeOut int
}

var config Config
var confEnv = "proxy"

func GetConfig() *Config {
	return &config
}

// get home dir of app, use PROXY_WORKDIR env var if present, else executable dir.
func exeDir() string {
	dir, exists := os.LookupEnv("PROXY_WORKDIR")
	if exists {
		return dir
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := path.Dir(ex)
		return exPath
	}
}

func GetConfigPath() string {
	baseDir := exeDir()
	confPath := path.Join(baseDir, "conf/")
	return confPath
}

// 参数解析
func updateConfig(c *Config) {
	vars := c.Var
	proxyList := &(c.Server.Proxy)
	for i := range *proxyList {
		p := &(*proxyList)[i]
		for k, v := range vars {
			fmt.Println("k=", k)
			fmt.Println("v=", v)
			k := "$" + k
			if strings.Contains(p.ProxyPass, k) {
				p.ProxyPass = strings.Replace(p.ProxyPass, k, v, 1)
			}
		}
	}
}

// 解析环境变量
func initVars(c *Config) {
	for k, v := range c.Var {
		// 环境变量
		if strings.HasPrefix(v, "$") {
			envVal := string([]byte(v)[1:])
			v = os.Getenv(envVal)
		} else if strings.HasPrefix(v, `\$`) {
			v = string([]byte(v)[1:])
		}
		c.Var[k] = v
	}
}

func init() {
	// 需要配置项目根目录的环境变量，方便执行test
	confPath := GetConfigPath()
	viper.SetConfigName(confEnv)   // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(confPath) // 第一个搜索路径
	viper.WatchConfig()           // 监控配置文件热重载
	err := viper.ReadInConfig()   // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	initVars(&config)
	updateConfig(&config)
	fmt.Println(config)
	if err != nil {
		panic(err)
	}
}

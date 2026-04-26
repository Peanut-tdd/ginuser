package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App      App      `json:"app"`
	Database Database `json:"database"`
	Redis    Redis    `json:"redis"`
	Jwt      Jwt      `json:"jwt"`
}

type App struct {
	Name    string
	Version string
	Env     string
	Port    string
	Debug   bool
}

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Username string `json:"username"`
	Charset  string `json:"charset"`
	Dbname   string `json:"dbname"`
}

type Redis struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	PoolSize int    `json:"pool_size"`
}

type Jwt struct {
	Secret        string `json:"secret"`
	AccessExpire  int    `json:"access_expire"`
	RefreshExpire int    `json:"refresh_expire"`
}

var Conf Config

func InitConfig() {
	workDir, _ := os.Getwd()
	fmt.Println("current workdir:", workDir)
	env := os.Getenv("env")
	if env == "" {
		env = "dev"
	}
	fmt.Println("env:", env)

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")

	// 添加多个搜索路径，确保无论在根目录还是 ops 目录都能找到配置
	viper.AddConfigPath("./etc")      // 本地 go run (在项目根目录)
	viper.AddConfigPath("../etc")     // 在子目录下执行测试
	viper.AddConfigPath("/app/etc")   // Docker 容器内部路径

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		return
	}

	// 兼容 Docker Compose 环境：如果检测到在容器内运行，且配置指向 localhost，则自动替换为服务名
	if os.Getenv("GIN_MODE") == "release" {
		if Conf.Database.Host == "localhost" || Conf.Database.Host == "127.0.0.1" {
			Conf.Database.Host = "Mysql" // 对应 docker-compose.yaml 中的服务名
		}
		// Redis 地址处理
		if Conf.Redis.Address == "localhost:6379" || Conf.Redis.Address == "127.0.0.1:6379" {
			Conf.Redis.Address = "Redis:6379"
		}
	}

	b, _ := json.Marshal(Conf)
	fmt.Println(string(b))
}

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
	fmt.Println(workDir)
	env := os.Getenv("env")
	if env == "" {
		env = "dev"
	}
	fmt.Println("env:", env)
	viper.AddConfigPath("./etc")
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		return
	}

	b, _ := json.Marshal(Conf)
	fmt.Println(string(b))

}

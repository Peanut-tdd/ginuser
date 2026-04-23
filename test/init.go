package test

import (
	"gin-user/config"
)

func init() {
	config.InitConfig()
	config.InitLogger()
	config.InitDb()
	config.InitRedis()

}

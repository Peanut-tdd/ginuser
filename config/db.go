package config

import (
	"context"
	"fmt"
	"gorm.io/plugin/opentelemetry/tracing"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	DbConfig := Conf.Database

	dbDns := strings.Join([]string{DbConfig.Username, ":", DbConfig.Password, "@tcp(", DbConfig.Host, ")/", DbConfig.Dbname, "?charset=", DbConfig.Charset, "&parseTime=true"}, "")

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbDns, // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	// 这一行，这种是不带 metric 的，详细使用可以看官方文档
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	DB = db

	fmt.Println("init db success")

}

func NewDBClient(ctx context.Context) *gorm.DB {

	db := DB

	return db.WithContext(ctx)
}

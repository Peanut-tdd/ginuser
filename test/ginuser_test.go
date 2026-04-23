package test

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-user/config"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func TestRedis(t *testing.T) {
	key := "ginuser:tdd"
	data := map[string]interface{}{
		"age":  30,
		"name": "tdd1111",
	}

	// 将 map 序列化为 JSON 字符串，因为 redis 不支持直接存储 map
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	// 这里改为 30 秒
	du := time.Duration(30) * time.Minute

	fmt.Println(du, du.Seconds())

	// 确保 RedisClient 已初始化
	if config.RedisClient == nil {
		t.Fatal("RedisClient 未初始化，请确保运行测试时包含了 init.go (例如使用 go test -v ./test)")
	}

	val, err := config.RedisClient.Set(context.Background(), key, jsonData, du).Result()
	if err != nil {
		fmt.Println("Redis Set Error:", err)
		t.Fail()
	}

	fmt.Println(val, err)

	fmt.Println("=========")

	r, err := config.RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println("Redis Get Error:", err)
	}

	u := map[string]interface{}{}
	json.Unmarshal([]byte(r), &u)

	fmt.Println(u, err)

	// JSON 中的数字在解析到 map[string]interface{} 时，默认会解析为 float64
	// 方式 1: 直接断言为 float64 再转换
	if value, ok := u["age"].(float64); ok {
		fmt.Println("age (from float64):", int(value))
	}

	// 方式 2: 使用 cast 库进行转换 (推荐)
	age := cast.ToInt(u["age"])
	fmt.Println("age (from cast):", age)

	if value, ok := u["name"].(string); ok {
		fmt.Println("name:", value)
	}
}

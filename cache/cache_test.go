package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

type Value struct {
	Id   string `redis:"id"`
	Name string `redis:"name"`
}

func (v *Value) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v *Value) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, v)
}

var rdb *redis.Client

// TestMain 设置 Redis 客户端并运行所有测试
func TestMain(m *testing.M) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "",
		Password: "",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	// 运行测试
	m.Run()
}

func TestLoad(t *testing.T) {
	vlaue := &Value{}
	testCache := NewCache(rdb, "test")
	v, exits := testCache.Load("test1")
	if !exits {
		t.Fatalf("expected value to exist in cache, but it does not")
	}

	if err := mapstructure.Decode(v, vlaue); err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(vlaue)
}

// TestStore 测试缓存的存储和加载
func TestStoreAndLoad(t *testing.T) {
	// 初始化缓存
	testCache := NewCache(rdb, "test")

	// 存储数据
	if err := testCache.Store("test1", &Value{
		Id:   "test1",
		Name: "test1",
	}); err != nil {
		t.Fatalf("failed to store value: %v", err)
	}

	// 加载数据并检查是否存在
	v, exists := testCache.Load("test1")
	if !exists {
		t.Fatalf("expected value to exist in cache, but it does not")
	}

	fmt.Println(">>>>", v)
}

func TestDelete(t *testing.T) {
	// 初始化缓存
	testCache := NewCache(rdb, "test")

	if err := testCache.Delete("test1"); err != nil {
		panic(err)
	}
}

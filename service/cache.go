package service

/*
基于分布式缓存的cache实现
基于redis
*/
import (
	"encoding/json"

	"github.com/go-redis/redis"
)

type ListObject []interface{}
type Object interface{}
type Cache interface {
	//获取缓存对象，如果存在返回true
	Get(key string, result Object) bool
	//设置缓存对象
	Set(key string, value Object)

	//设置缓存对象,并设置过期时间
	SetNx(key string, value Object, extime int64)
	//获取列表
	GetAsList(key string, typeTemplate Object) ListObject
	//向列表中添加对象
	Add(key string, value Object)
	//获取map
	GetAsMap(key string, typeTemplate Object) map[string]Object
	//向map中添加对象
	AddToMap(key string, subkey string, value Object)
	//设置过期时间
	SetEx(key string, extime int64)
}

type RedisCache struct {
	addr   string //地址
	db     int    //数据库
	pwd    string //密码
	client *redis.Client
}

func (this *RedisCache) Init(addr string, pwd string, db int) {
	//防止多次初始化
	if this.client!=nil {
		return
	}

	this.addr = addr
	this.db = db
	if db < 0 || db > 16 {
		this.db = 0
	}
	this.pwd = pwd

	this.client = redis.NewClient(&redis.Options{
		Addr:     this.addr,
		Password: this.pwd, // no password set
		DB:       this.db,  // use default DB
	})

}

//获取缓存对象，如果存在返回true
func (this *RedisCache) Get(key string, result Object) bool {
	value, err := this.client.Get(key).Result()
	if err == nil && value != "" && value != "nil" {
		err = json.Unmarshal([]byte(value), result)

		return true
	}

	return false

}

//	//设置缓存对象
func (this *RedisCache) Set(key string, value Object) {
	if value != nil && key != "" {
		data, _ := json.Marshal(value)
		this.client.Set(key, string(data), 0)
	}

}

//获取列表
func (this *RedisCache) GetAsList(key string, typeTemplate Object) ListObject {
	if typeTemplate != nil && key != "" {
		length, _ := this.client.LLen(key).Result()
		results := make(ListObject, length)
		var index int64
		for index = 0; index < length; index++ {
			value, _ := this.client.LIndex(key, index).Result()
			obj := CreateObjByType(GetRealType(typeTemplate))
			json.Unmarshal([]byte(value), obj)
			results[index] = obj
		}
		return results

	}
	return nil
}

//向列表中添加对象
func (this *RedisCache) Add(key string, value Object) {
	if key != "" && value != nil {
		data, _ := json.Marshal(value)
		this.client.RPush(key, string(data))

	}
}

//获取map
func (this *RedisCache) GetAsMap(key string, typeTemplate Object) map[string]Object {
	if key != "" && typeTemplate != nil {
		maps, _ := this.client.HGetAll(key).Result()
		if maps != nil {
			results := make(map[string]Object, len(maps))
			for key, value := range maps {
				obj := CreateObjByType(GetRealType(typeTemplate))
				json.Unmarshal([]byte(value), obj)
				results[key] = obj
			}
			return results

		}

	}

	return nil

}

//向map中添加对象
func (this *RedisCache) AddToMap(key string, subkey string, value Object) {
	if key != "" && subkey != "" && value != nil {
		data, _ := json.Marshal(value)
		this.client.HSet(key, subkey, string(data))
	}

}

//设置过期时间 second
func (this *RedisCache) SetEx(key string, extime int64) {
	if key != "" && extime > 0 {
		this.client.Expire(key, ToSecond(extime))
	}

}

//设置缓存对象,并设置过期时间
func (this *RedisCache) SetNx(key string, value Object, extime int64) {
	if key != "" && extime > 0 && value != nil {
		data, _ := json.Marshal(value)
		this.client.SetNX(key, string(data), ToSecond(extime))
	}

}



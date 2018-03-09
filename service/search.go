package service

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo"
	"github.com/aosfather/bingo/utils"
	"github.com/go-redis/redis"
	"strconv"
	"crypto/md5"
)

/**
  搜索实现
  通过倒排实现关键信息的实现
  规则：
   1、原始内容使用hashmap存储，对象【ID，Content】，键值 indexname，二级key根据md5 对象转json字符串
   2、针对原始内容带的标签，key value，生成set，名称为 indexname_key_value的形式，set内容放 通过内容md5出来的键值
   3、搜索的时候，根据传递的搜索条件 key，value 数组，对找到的set（形如indexname_key_value ），实行找交集
   4、根据交集的结果二级key，并从indexname的hashmap中获取json内容

*/

type Field struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type targetObject struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}
type sourceObject struct {
	targetObject
	Fields []Field `json:"fields"`
}

type searchEngine struct {
	indexs map[string]*searchIndex
	client *redis.Client
	logger utils.Log
}

func (this *searchEngine) Init(context *bingo.ApplicationContext) {
	db, err := strconv.Atoi(context.GetPropertyFromConfig("bingo.search.db"))
	if err != nil {
		db = 0
	}
	this.client = redis.NewClient(&redis.Options{
		Addr:     context.GetPropertyFromConfig("bingo.search.redis"),
		Password: "", // no password set
		DB:       db,
	})
	this.indexs = make(map[string]*searchIndex)
	this.logger = context.GetLog("bingo_search")
}

func (this *searchEngine) CreateIndex(name string) *searchIndex {
	if name != "" {
		index := this.indexs[name]
		if index == nil {
			index = &searchIndex{name, this}
			this.indexs[name] = index
		}
		return index
	}

	return nil
}

func (this *searchEngine) LoadSource(name string, obj *sourceObject) {

	index := this.CreateIndex(name)
	if index != nil {
		index.LoadObject(obj)
	}

}

func (this *searchEngine) Search(name string, input ...Field) []targetObject {
	if name != "" {
		index := this.indexs[name]
		if index != nil {
			return index.Search(input...)
		}
		this.logger.Info("not found index %s", name)
	}

	return nil
}

type searchIndex struct {
	name   string
	engine *searchEngine
}

//搜索信息
func (this *searchIndex) Search(input ...Field) []targetObject {
	//搜索索引
	var searchkeys []string
	for _, f := range input {
		searchkeys = append(searchkeys, this.buildTheKey(f))
	}
	//取交集
	result := this.engine.client.SInter(searchkeys...)
	targetkeys, err := result.Result()
	if err != nil {
		this.engine.logger.Error("inter key error!%s", err.Error())
		return nil
	}
	if len(targetkeys) > 0 {
		//根据最后的id，从data中取出所有命中的元素
		datas, err1 := this.engine.client.HMGet(this.name, targetkeys...).Result()
		if err1 == nil && len(datas) > 0{

				var targets []targetObject

				for _, v := range datas {
					if v != nil {
						t := targetObject{}
						json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &t)
						targets = append(targets, t)
					}
				}

				return targets

		} else {
			this.engine.logger.Error("get data by index error!%s", err1.Error())
		}

	}

	return nil
}

//刷新索引，加载信息到存储中
func (this *searchIndex) LoadObject(obj *sourceObject) {
	data, _ := json.Marshal(obj)
	key := getMd5str(string(data))
	//1、放入数据到目标集合中
	this.engine.client.HSet(this.name, key, string(data))

	//2、根据field存储到各个对应的索引中

	for _, field := range obj.Fields {
		this.engine.client.SAdd(this.buildTheKey(field), key)
	}

}

func (this *searchIndex) buildTheKey(f Field) string {
	return fmt.Sprintf("%s_%s_%s", this.name, f.Key, f.Value)
}

func getMd5str(value string) string {
	data := []byte(value)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1

}
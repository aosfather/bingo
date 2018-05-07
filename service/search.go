package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aosfather/bingo"
	"github.com/aosfather/bingo/utils"
	"github.com/go-redis/redis"
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

type PageSearchResult struct {
	Id    string `json:"uuid"` //查询的请求id
	Index int64  `json:"page"` //页码
	Data  []TargetObject
}

type TargetObject struct {
	Id   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}
type SourceObject struct {
	TargetObject
	Fields map[string][]string `json:"fields"`
}

type FieldType byte
const (
	FT_TEXT FieldType =11  //文本类型
	FT_NUMBER FieldType =9 //数字
	FT_ENUM   FieldType =8 //枚举
	FT_ID FieldType =7  //id唯一标识
	FT_DATE FieldType =6 //日期
)
//索引的元数据信息
type IndexMeta struct {
	Name string  //索引名称
	Fields []FieldMeta //字段
}

type FieldMeta struct {
	Name string //字段名称
	Type FieldType //类型

}

type SearchEngine struct {
	indexs   map[string]*searchIndex
	client   *redis.Client
	logger   utils.Log
	pageSize int64
	pageLife int64 //分钟
}

func (this *SearchEngine) Init(context *bingo.ApplicationContext) {
	fmt.Println("init .....")
	db, err := strconv.Atoi(context.GetPropertyFromConfig("service.search.db"))
	if err != nil {
		db = 0
	}

	size, err := strconv.Atoi(context.GetPropertyFromConfig("service.search.pagesize"))
	if err != nil {
		size = 20 //默认大小20条
	}
	this.pageSize = int64(size)

	life, err := strconv.Atoi(context.GetPropertyFromConfig("service.search.pagelife"))
	if err != nil {
		life = 10 //默认时间10分钟
	}
	this.pageLife = int64(life)

	this.client = redis.NewClient(&redis.Options{
		Addr:     context.GetPropertyFromConfig("service.search.redis"),
		Password: "", // no password set
		DB:       db,
	})
	fmt.Println(context.GetPropertyFromConfig("service.search.redis"))
	this.indexs = make(map[string]*searchIndex)
	this.logger = context.GetLog("bingo_search")
}

//创建索引
func (this *SearchEngine) CreateIndex(name string) *searchIndex {
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

//清除索引，将整个索引的数据摧毁
func (this *SearchEngine)FlushIndex(name string) {
	//1、删除存放的数据
	this.client.Del(name)

	//2、删除所有的索引key
	keys,err:=this.client.Keys(name+"_*").Result()
	if err!=nil {
		this.logger.Debug("get index keys error:%s",err.Error())
	}

	if keys!=nil && len(keys)>0 {
		this.client.Del(keys...)
	}
}


//加载和刷新数据
func (this *SearchEngine) LoadSource(name string, obj *SourceObject) {
	if name==""||obj==nil {
		return
	}

	index := this.CreateIndex(name)
	if index != nil {
		index.LoadObject(obj)
	}

}

//删除数据
func (this *SearchEngine)RemoveSource(name string, obj *SourceObject){
   if name!=""&&obj!=nil{
	   index := this.CreateIndex(name)
	   if index != nil {
		   index.RemoveObject(obj)
	   }
   }
}

//按页获取数据
func (this *SearchEngine) FetchByPage(request string, page int64) *PageSearchResult {
	if request != "" {
		//获取request的name
		index := strings.Index(request, ":")
		if index < 0 {
			this.logger.Error("pagerequest's index name not found !")
			return nil //找不到对应的索引类型
		}

		name := request[0:index]
		if page <= 0 {
			page = 1
		}
		startIndex := (page - 1) * this.pageSize //+ 1 //从1开始计数
		endIndex := page*this.pageSize - 1

		keys, err := this.client.LRange(request, startIndex, endIndex).Result()
		if err != nil || len(keys) == 0 {
			this.logger.Debug("no content by page!")
			return nil
		}
		go this.client.Expire(request, time.Duration(this.pageLife)*time.Minute) //更新重置失效时间
		return &PageSearchResult{request, page, this.fetch(name, keys...)}
	}

	return nil
}

func (this *SearchEngine) createRequst(key string, keys ...string) {
	//    key:= getSearchRequestUuid(name)
	var datas []interface{}

	for _, v := range keys {
		datas = append(datas, v)
	}

	this.client.LPush(key, datas...)
	this.client.Expire(key, time.Duration(this.pageLife)*time.Minute) //指定x分钟后失效
	//	return key
}

//获取内容
func (this *SearchEngine) fetch(name string, keys ...string) []TargetObject {
	datas, err1 := this.client.HMGet(name, keys...).Result()
	if err1 == nil && len(datas) > 0 {

		var targets []TargetObject

		for _, v := range datas {
			if v != nil {
				t := TargetObject{}
				json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &t)
				targets = append(targets, t)
			}
		}

		return targets

	} else {
		this.logger.Error("get data by index error!%s", err1.Error())
	}

	return nil
}

func (this *SearchEngine) Search(name string, input ...Field) *PageSearchResult {
	if name != "" {
		index := this.indexs[name]
		if index != nil {
			r, data := index.Search(input...)
			return &PageSearchResult{r, 1, data}

		}
		this.logger.Info("not found index %s", name)
	}

	return nil
}

type searchIndex struct {
	name   string
	engine *SearchEngine
}

//搜索信息
func (this *searchIndex) Search(input ...Field) (string, []TargetObject) {
	//生成索引搜索请求号
	requestkey := getSearchRequestUuid(this.name)

	//搜索索引
	var searchkeys []string
	var tmpkeys []string
	for _, f := range input {
		//处理并集
		v := f.Value
		arrays := strings.Split(v, "|")
		if len(arrays) > 1 {
			skey := requestkey + ":" + f.Key
			var subkeys []string
			for _, subkey := range arrays {
				subkeys = append(subkeys, this.buildTheKeyByItem(f.Key, subkey))
			}

			//取并集，将结果存入临时的key中，用于后续取交集用
			this.engine.client.SUnionStore(skey, subkeys...)
			tmpkeys = append(tmpkeys, skey) //放入临时组中用于使用后删除
			searchkeys = append(searchkeys, skey)
		} else {
			searchkeys = append(searchkeys, this.buildTheKey(f))
		}

	}
	defer this.deleteTempkeys(tmpkeys) //删除临时创建的key
	//取交集
	result := this.engine.client.SInter(searchkeys...)
	targetkeys, err := result.Result()
	if err != nil {
		this.engine.logger.Error("inter key error!%s", err.Error())
		return "", nil
	}

	//取出第一页的key
	size := len(targetkeys)
	if size > 0 {
		//生成request
		this.engine.createRequst(requestkey, targetkeys...)

		var query []string
		if size > int(this.engine.pageSize) {
			//写入到列表中
			query = targetkeys[0:this.engine.pageSize]
		} else {
			query = targetkeys
		}

		//根据最后的id，从data中取出第一页的元素的元素
		return requestkey, this.engine.fetch(this.name, query...)
	}

	return "", nil
}

//删除临时创建的key
func (this *searchIndex) deleteTempkeys(keys []string) {
	if keys != nil && len(keys) > 0 {
		this.engine.client.Del(keys...)
	}
}

//刷新索引，加载信息到存储中
func (this *searchIndex) LoadObject(obj *SourceObject) {
	data, _ := json.Marshal(obj)
	key:=obj.Id
	//如果没有指明对象id，则使用md5作为数据的唯一标识
	if key=="" {
		key= getMd5str(string(data))
	}

	//1、放入数据到目标集合中
	this.engine.client.HSet(this.name, key, string(data))

	//2、根据field存储到各个对应的索引中
	for k, v := range obj.Fields {
		for _,vkey:=range v{
			this.engine.client.SAdd(this.buildTheKeyByItem(k, vkey), key)
		}

	}

}

//删除数据：删除数据及索引存的值
func (this *searchIndex)RemoveObject(obj *SourceObject) {
	data, _ := json.Marshal(obj)
	key:=obj.Id
	if key=="" {
		key=getMd5str(string(data))
	}
	//删除数据
	this.engine.client.HDel(this.name,key)
	//删除索引里的记录
	for k,v:=range obj.Fields {
		for _,vkey:=range v {
			this.engine.client.SRem(this.buildTheKeyByItem(k, vkey), key)
		}
	}


}

func (this *searchIndex) buildTheKey(f Field) string {
	return this.buildTheKeyByItem(f.Key, f.Value)
}

func (this *searchIndex) buildTheKeyByItem(key, value string) string {
	return fmt.Sprintf("%s_%s_%s", this.name, key, value)
}

func getMd5str(value string) string {
	data := []byte(value)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1

}

func getSearchRequestUuid(prefix string) string {
	return fmt.Sprintf("%s:%d", prefix, time.Now().UnixNano())
}

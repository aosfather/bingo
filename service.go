package main

import (
	"github.com/aosfather/bingo_mvc/sqltemplate"
	"gopkg.in/yaml.v2"
)

//基本服务
//1、字典
//上传、下载、导入、导出
//-----------------字典服务------------------//
//词条
type DictCatalog struct {
	Code  string
	Label string
	Tip   string
	Items []DictCatalogItem
}

//词条的选择项
type DictCatalogItem struct {
	Code    string            //值
	Label   string            //显示值
	Tip     string            //提示
	Virtual bool              //是否虚拟,表示存在有同名的词条
	Ord     int               //显示次序
	Extends map[string]string //扩展属性
}
type Cache interface {
	Get(c string) (interface{}, bool)
	Set(key string, value interface{})
}

type TypesDBMeta struct {
	Dao       *sqltemplate.MapperDao `Inject:""`
	TypeCache Cache
	DictCache Cache
}

func (this *TypesDBMeta) GetType(code string) *DataType {
	debug("get type from db meta! ", code)
	if data, ok := this.TypeCache.Get(code); ok {
		return data.(*DataType)
	}
	item := &MetaItem{Code: code, Catalog: 0} //数据类型
	if this.Dao.FindByObj(item, "Code", "Catalog") {
		debug("query dictioary from db for ", code)
		dictitem := &DataType{}
		debug(item.Content)

		err := yaml.Unmarshal([]byte(item.Content), dictitem)
		if err != nil {
			debug(err.Error())
		}
		this.TypeCache.Set(code, dictitem)
		debug(dictitem)
		return dictitem
	}
	return nil
}

func (this *TypesDBMeta) GetDictionary(code string) *DictCatalog {
	debug("get dictionary from db meta! ", code)
	if data, ok := this.DictCache.Get(code); ok {
		return data.(*DictCatalog)
	}
	item := &MetaItem{Code: code, Catalog: 1} //字典类型
	if this.Dao.FindByObj(item, "Code", "Catalog") {
		debug("query dictioary from db for ", code)
		dictitem := &DictCatalog{}
		debug(item.Content)
		err := yaml.Unmarshal([]byte(item.Content), dictitem)
		if err != nil {
			debug(err.Error())
		}
		this.DictCache.Set(code, dictitem)
		debug(dictitem)
		return dictitem
	}
	return nil
}

type MetaItem struct {
	Code    string `Table:"bingo_dictionary" Option:"pk"` //脚本唯一编码
	Catalog int    `Option:"pk"`
	Content string
}

//------------------上传下载--------------------//

//-------------------导入导出-------------------//

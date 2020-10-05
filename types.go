package main

import (
	"fmt"
	"github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
)

/**
  基本类型

*/
//-----------------类型管理---------------------//
var types *TypesManager

//校验器列表
var validates map[string]ValidateFunc

/**
  参数取值转换
*/
var converts map[string]ConvertFunc

type ConvertFunc func(expr string, value string) string
type TypeMeta interface {
	GetType(code string) *DataType
	GetDictionary(code string) *DictCatalog
}

//非法字符校验，防止SQL注入
// 正则过滤sql注入的方法
// 参数 : 要匹配的语句
var sqlcheckPattern *regexp.Regexp

func init() {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	sqlcheckPattern = regexp.MustCompile(str)
	//初始化types manager
	types = new(TypesManager)
	validates = make(map[string]ValidateFunc)
	validates["regex"] = validateByRegexp
	validates["dict"] = types.validateByDict
	converts = make(map[string]ConvertFunc)
	converts["Enum"] = types.GetValueByDict
}

func ValidateBySQLCheck(v string, option string) (bool, string) {
	//过滤 ‘
	result := sqlcheckPattern.MatchString(v)
	if result {
		return false, fmt.Sprintf("'%s'中存在有非法字符，怀疑有sql注入", v)
	}
	return true, ""
}

/**
  使用正则表达式进行校验
*/
func validateByRegexp(v string, option string) (bool, string) {

	pattern, _ := regexp.Compile(option)

	if pattern != nil {
		result := pattern.Match([]byte(v))
		if result {
			return true, ""
		}
	}

	return false, "regex校验不同过！"

}

type TypesManager struct {
	//类型
	meta TypeMeta
}

func (this *TypesManager) validateType(typeName string, value string, pname, expr string) error {
	if t := this.meta.GetType(typeName); t != nil {
		//存在
		return t.validate(value, pname, expr)
	} else {
		return fmt.Errorf("指定的类型[%s],不存在", typeName)
	}

	return nil
}

func (this *TypesManager) GetValueByDict(catalog string, v string) string {
	if catalog == "" {
		return v
	}
	if c := this.meta.GetDictionary(catalog); c != nil {
		for _, item := range c.Items {
			if item.Code == v || item.Label == v {
				return item.Code
			}
		}

	}
	return v
}

//校验字典的值
func (this *TypesManager) validateByDict(v string, catalog string) (bool, string) {
	if catalog == "" {
		return true, ""
	}
	if c := this.meta.GetDictionary(catalog); c != nil {
		for _, item := range c.Items {
			if item.Code == v || item.Label == v {
				return true, ""
			}
		}
		return false, fmt.Sprintf("'%s'不是字典[%s]的合法成员", v, catalog)

	}
	return false, fmt.Sprintf("指定的字典[%s],不存在", catalog)
}

//校验函数
type ValidateFunc func(v string, option string) (bool, string)

func (this *ValidateFunc) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	*this = validates[text]
	return nil
}

//通用数据类型
type DataType struct {
	Code      string       //类型名称
	Label     string       //类型
	Validater ValidateFunc `yaml:"validate"` //校验函数
	Option    string       //校验额外设置
}

func (this *DataType) validate(value string, pname, expr string) error {
	if this.Validater != nil {
		b, msg := this.Validater(value, this.Option)
		if !b {
			return fmt.Errorf("'%s'不符合类型[%s]的规则，校验失败！原因是：‘%s’", value, this.Code, msg)
		}
		//校验额外规则
		if expr != "" {
			b, msg = this.Validater(value, expr)
			if !b {
				return fmt.Errorf("'%s'不符合类型[%s]的规则，校验失败！原因是：‘%s’", value, pname, msg)
			}
		}

	}
	return nil
}

func GetDict(name string) DictCatalog {
	if types != nil {
		dic := types.meta.GetDictionary(name)
		if dic != nil {
			return *dic
		}
	}
	return DictCatalog{}
}

type typeConfigs struct {
	Types []DataType
	Enums []DictCatalog
}
type YamlFileTypeMeta struct {
	types      map[string]*DataType
	dictionary map[string]*DictCatalog
}

func (this *YamlFileTypeMeta) Load(f string) {
	this.types = make(map[string]*DataType)
	this.dictionary = make(map[string]*DictCatalog)
	if files.IsFileExist(f) {
		tf := &typeConfigs{}
		data, err := ioutil.ReadFile(f)
		if err == nil {
			err = yaml.Unmarshal(data, tf)
		}
		if err != nil {
			errs("load type meta error:", err.Error())
			return
		}

		for _, item := range tf.Types {
			this.types[item.Code] = &item
			debug(item)
		}

		for _, catalog := range tf.Enums {
			this.dictionary[catalog.Code] = &catalog
		}
	}
}

func (this *YamlFileTypeMeta) GetType(code string) *DataType {
	return this.types[code]
}
func (this *YamlFileTypeMeta) GetDictionary(code string) *DictCatalog {
	return this.dictionary[code]
}

//------------------------------新模型------------------------------------------//

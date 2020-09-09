package main

import (
	"fmt"
	"regexp"
	"strings"
)

//表单
//案例元信息
type FormMeta struct {
	Code        string      `yaml:"code"`
	Author      string      `yaml:"author"`
	Version     string      `yaml:"version"`
	UpdateDate  string      `yaml:"updateDate"`
	Description string      `yaml:"description"`
	Parameters  []Parameter `yaml:"parameters"`
	ScriptType  string      `yaml:"scriptType"`
	Script      string      `yaml:"script"`
}

type Parameter struct {
	Name       string      `yaml:"name"`
	Policy     string      `yaml:"policy"`
	Label      string      `yaml:"label"`
	Type       string      `yaml:"type"`
	Expr       string      `yaml:"expr"` //表达式
	Conditions []Condition `yaml:"link"` //关联条件，当为 Maybe 的时候使用。
}

func (this *Parameter) validate(v string) error {
	//检查是否必填
	if this.Policy == "Must" {
		if v == "" {
			return fmt.Errorf("参数[%s:%s]为必填参数！", this.Name, this.Label)
		}
	}

	//类型校验
	err := types.validateType(this.Type, v, this.Name, this.Expr)
	if err != nil {
		return fmt.Errorf("参数[%s:%s]:%s", this.Name, this.Label, err.Error())
	}
	//文本不适用SQL注入检查
	if this.Type == "Text" {
		return nil
	}
	//sql注入校验
	b, msg := ValidateBySQLCheck(v, "")
	if !b {
		return fmt.Errorf("SQL注入检查结果:%v,%s", b, msg)
	}
	return nil
}

type Condition struct {
	Name   string   `yaml:"name"`
	Fields []string `yaml:"fields"`
}

func (this *Condition) Validate(name string, v string, datas map[string]string) error {
	on := false
	for _, f := range this.Fields {
		expr := strings.Split(f, "=")
		size := len(expr)
		if d, ok := datas[expr[0]]; ok {
			if size == 1 {
				if d != "" {
					on = true
					continue
				}
			} else { //如果输入了值，当取值等于指定的值时候才触发条件是否满足
				if d == expr[1] {
					on = true
					continue
				}
			}

		}
		on = false
		break
	}

	if on {
		if v == "" {
			return fmt.Errorf("参数[%s]为必填参数！", name)
		}

	}
	return nil
}

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

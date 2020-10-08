package main

import (
	"fmt"
	"github.com/aosfather/bingo_mvc/dd"
	"regexp"
	"strings"
)

type FormMetaManager interface {
	GetFormMeta(name string) *FormMeta
}

//表单
//案例元信息
type FormMeta struct {
	Code        string            `yaml:"code"`
	Author      string            `yaml:"author"`
	Version     string            `yaml:"version"`
	UpdateDate  string            `yaml:"updateDate"`
	FormType    string            `yaml:"type"`        //表单类型
	Title       string            `yaml:"title"`       //表单标题
	Description string            `yaml:"description"` //表单说明
	Action      string            `yaml:"action"`      //表单对应的动作
	Parameters  []Parameter       `yaml:"parameters"`  //参数定义
	Response    ResponseProcessor `yaml:"response"`    //返回处理
	ScriptType  string            `yaml:"scriptType"`  //脚本类型
	Extends     map[string]string `yaml:"extends"`
	Script      string            `yaml:"script"`    //脚本内容
	ResultSet   []ResultField     `yaml:"resultset"` //结果集合
	Toolbar     []Tool            `yaml:"toolbar"`   //表头工具
	Tools       []Tool            `yaml:"tools"`     //工具条
	JSscript    string            `yaml:"jsscript"`  //js脚本
}

func (this *FormMeta) ValidateInput(data map[string]interface{}) (error, map[string]interface{}) {
	for _, p := range this.Parameters {
		v := data[p.Name].(string)
		//参数校验
		err := p.validate(v)
		if err != nil {
			errs("validate failed! ", err.Error())
			return err, nil
		}
		//maybe 方式校验
		if p.Policy == "Maybe" {
			for _, c := range p.Conditions {
				err = c.Validate(p.Name, v, data)
				if err != nil {
					errs("maybe validate failed,", err.Error())
					return err, nil
				}
			}
		}

		//参数转换

		data[p.Name] = dd.InputByDataType(p.Type, v)
		debug("convert input<", v, "-->", data[p.Name], ">")

	}
	return nil, data

}

const (
	//直接返回
	PT_DIRECT  ProcessorType = 1 << iota
	PT_DEFAULT               //默认返回方式

)

//工具条定义
type Tool struct {
	Name      string
	Label     string
	Condition []string
}
type ProcessorType byte

func (this *ProcessorType) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	text = strings.ToLower(text)
	switch text {
	case "direct":
		*this = PT_DIRECT
	default:
		*this = PT_DEFAULT
	}

	return nil
}

type ResponseProcessor struct {
	Type    ProcessorType     `yaml:"type"`
	Options map[string]string `yaml:"options"`
}

//结果字段
type ResultField struct {
	Name      string `yaml:"name"`
	Label     string `yaml:"label"`
	InnerName string `yaml:"inner"`
	Type      string `yaml:"type"`
	Expr      string `yaml:"expr"` //表达式
}

//输入参数
type Parameter struct {
	Name       string      `yaml:"name"`
	Policy     string      `yaml:"policy"`
	Label      string      `yaml:"label"`
	Type       string      `yaml:"type"`
	InputTip   string      `yaml:"tip"`      //输入提示
	Verify     string      `yaml:"verify"`   //前端校验规则
	Expr       string      `yaml:"expr"`     //表达式
	Readonly   bool        `yaml:"readonly"` //只读方式
	Conditions []Condition `yaml:"link"`     //关联条件，当为 Maybe 的时候使用。
}

func (this *Parameter) validate(v string) error {
	//检查是否必填
	if this.Policy == "Must" {
		if v == "" {
			return fmt.Errorf("参数[%s:%s]为必填参数！", this.Name, this.Label)
		}
	}

	//类型校验
	err := dd.ValidateByDataType(this.Name, v, this.Type) //types.validateType(this.Type, v, this.Name, this.Expr)
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

func (this *Condition) Validate(name string, v string, datas map[string]interface{}) error {
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

//非法字符校验，防止SQL注入
// 正则过滤sql注入的方法
// 参数 : 要匹配的语句
var sqlcheckPattern *regexp.Regexp

func init() {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	sqlcheckPattern = regexp.MustCompile(str)

}

func ValidateBySQLCheck(v string, option string) (bool, string) {
	//过滤 ‘
	result := sqlcheckPattern.MatchString(v)
	if result {
		return false, fmt.Sprintf("'%s'中存在有非法字符，怀疑有sql注入", v)
	}
	return true, ""
}

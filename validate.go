package bingo

/*
Bean Validation 中内置的 constraint

@Null   被注释的元素必须为 null
@NotNull    被注释的元素必须不为 null
@AssertTrue     被注释的元素必须为 true
@AssertFalse    被注释的元素必须为 false
@Min(value)     被注释的元素必须是一个数字，其值必须大于等于指定的最小值
@Max(value)     被注释的元素必须是一个数字，其值必须小于等于指定的最大值
@DecimalMin(value)  被注释的元素必须是一个数字，其值必须大于等于指定的最小值
@DecimalMax(value)  被注释的元素必须是一个数字，其值必须小于等于指定的最大值
@Size(max=, min=)   被注释的元素的大小必须在指定的范围内
@Digits (integer, fraction)     被注释的元素必须是一个数字，其值必须在可接受的范围内
@Past   被注释的元素必须是一个过去的日期
@Future     被注释的元素必须是一个将来的日期
@Pattern(regex=,flag=)  被注释的元素必须符合指定的正则表达式

Hibernate Validator 附加的 constraint
@NotBlank(message =)   验证字符串非null，且长度必须大于0
@Email  被注释的元素必须是电子邮箱地址
@Length(min=,max=)  被注释的字符串的大小必须在指定的范围内
@NotEmpty   被注释的字符串的必须非空
@Range(min=,max=,message=)  被注释的元素必须在合适的范围内





*/
import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

//格式:Valid:"name(rule);name(rule)"。其中name(rule):exp表达式
const (
	_VALIDATE_TAG = "Valid"
)

type Validater interface {
	Validate(obj interface{}) BingoError
}

type ValidaterFactory interface {
	CreateValidater(exp string) Validater
}

type validateManager struct {
	factory ValidaterFactory
	caches  map[string]*Validater
}

func (this *validateManager) Init(factory ValidaterFactory) {
	if this.factory == nil {
		this.factory = factory
		this.caches = make(map[string]*Validater)
	}
}

func (this *validateManager) getValidater(key string) *Validater {
	if key != "" {
		v := this.caches[key]
		if v == nil {
			validater := this.factory.CreateValidater(key)
			v = &validater
			this.caches[key] = v
		}
		return v
	}
	return nil

}

func (this *validateManager) Validate(obj interface{}) []BingoError {
	if isMap(obj) {
		//TODO 需要看怎么做校验了
		return nil
	}
	objT, objV, err := getStructTypeValue(obj)
	if err != nil {
		return []BingoError{CreateError(501, err.Error())}
	}
	var errors []BingoError
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		tag := f.Tag.Get(_VALIDATE_TAG)
		v := vf.Interface()
		if len(tag) == 0 {
			if vf.Kind() == reflect.Struct || (vf.Kind() == reflect.Ptr && vf.Elem().Kind() == reflect.Struct) {
				errors = append(errors, this.Validate(v)...)
			}

		} else {

			rules := strings.Split(tag, ";")
			errors = append(errors, this.ValidateByrules(v, rules...)...)
		}

	}
	return errors

}

func (this *validateManager) ValidateByrules(obj interface{}, rules ...string) []BingoError {
	var errors []BingoError
	for _, rule := range rules {
		v := this.getValidater(rule)
		err := (*v).Validate(obj)
		if err != nil {
			errors = append(errors, err)
		}

	}
	return errors
}

type defaultValidaterFactory struct {
}

func (this *defaultValidaterFactory) CreateValidater(exp string) Validater {
	vexp := strings.TrimSpace(exp)
	ruleStart := strings.Index(vexp, "(")
	var vname, rule string
	if ruleStart < 0 {
		vname = strings.ToLower(vexp)
		rule = ""
	} else {
		vname = strings.TrimSpace(vexp[:ruleStart])
		vname = strings.ToLower(vname)
		ruleEnd := strings.Index(vexp, ")")
		if ruleEnd < 0 {
			ruleEnd = len(vexp)
		}
		rule = strings.TrimSpace(vexp[ruleStart+1 : ruleEnd])
	}

	var v Validater
	switch vname {
	case "required":
		v = &required{}
	case "numeric":
		v = &numeric{}
	case "min":
		t, _ := strconv.Atoi(rule)
		v = &compare{t, true}
	case "max":
		t, _ := strconv.Atoi(rule)
		v = &compare{t, false}
	case "match":
		t, _ := regexp.Compile(rule)
		v = &match{t, false}
	case "nomatch":
		t, _ := regexp.Compile(rule)
		v = &match{t, true}
	}

	return v
}

type numeric struct {
}

func (this numeric) Validate(obj interface{}) BingoError {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if '9' < v || v < '0' {
				return CreateError(501, "")
			}
		}
		return nil
	}
	if _, ok := obj.(int); ok {
		return nil
	}
	return CreateError(501, "")
}

type match struct {
	regexp *regexp.Regexp
	not    bool
}

func (this match) Validate(obj interface{}) BingoError {
	result := this.regexp.MatchString(fmt.Sprintf("%v", obj))

	if result != this.not {
		return nil
	}

	return CreateError(501, "match failed")

}

type compare struct {
	target int
	min    bool
}

func (this *compare) Validate(obj interface{}) BingoError {
	v := reflect.ValueOf(obj)
	var result int
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		result = v.Len() - this.target

	}
	if str, ok := obj.(string); ok {
		result = utf8.RuneCountInString(str) - this.target
	}

	if num, ok := obj.(int); ok {
		result = num - this.target
	}

	if this.min && result < 0 {

		return CreateError(501, "must max than ")

	}

	if !this.min && result > 0 {
		return CreateError(501, "must min than")
	}

	return nil
}

type required struct {
}

func (this *required) Validate(obj interface{}) BingoError {
	if objIsNotNil(obj) {
		return nil
	}
	return CreateError(501, "the value required!")
}

func objIsNotNil(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(str) > 0
	}
	if _, ok := obj.(bool); ok {
		return true
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if i, ok := obj.(uint); ok {
		return i != 0
	}
	if i, ok := obj.(int8); ok {
		return i != 0
	}
	if i, ok := obj.(uint8); ok {
		return i != 0
	}
	if i, ok := obj.(int16); ok {
		return i != 0
	}
	if i, ok := obj.(uint16); ok {
		return i != 0
	}
	if i, ok := obj.(uint32); ok {
		return i != 0
	}
	if i, ok := obj.(int32); ok {
		return i != 0
	}
	if i, ok := obj.(int64); ok {
		return i != 0
	}
	if i, ok := obj.(uint64); ok {
		return i != 0
	}
	if i, ok := obj.(float32); ok {
		return i != 0
	}
	if i, ok := obj.(float64); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true

}

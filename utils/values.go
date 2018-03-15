package utils

import (
	"reflect"
	"strings"
)

/*
提供属性注入
 */

 const (
 	_TAG_VALUE ="Value" //属性

 )

 type StoreFunction func(key string) string
 type ValueStore interface{
 	GetProperty(key string) string
 }

 type defaultStore struct {
 	function StoreFunction
 }
 func(this *defaultStore)GetProperty(key string) string{
 	if this.function!=nil {
 		return this.function(key)
	}
	return ""
 }

 type ValuesHolder struct{
 	store ValueStore
 }

 func(this *ValuesHolder)InitByFunction(f StoreFunction){
 	s:=defaultStore{f}
 	this.Init(&s)
 }

 func (this *ValuesHolder) Init(s ValueStore) {
 	if s!=nil {
 		this.store=s
	}
 }

 //处理赋值
 func (this *ValuesHolder) ProcessValueTag(v interface{}){
   if(IsStructPtr(v)){//处理struct
	   reflectType:=reflect.TypeOf(v)
	   reflectValue:=reflect.ValueOf(v)
	   for i := 0; i < reflectValue.Elem().NumField(); i++ {
		   field := reflectValue.Elem().Field(i)
		   fieldType := field.Type()
		   tag := reflectType.Elem().Field(i).Tag.Get(_TAG_VALUE)
		   if tag!="" {
		   	   if field.CanSet(){
				   this.setValue(field,fieldType,this.store.GetProperty(tag))
			   }else{
			   	   //查找setmethod进行设置
			   	   fieldName:=reflectType.Elem().Field(i).Name
				   println("use set method,beacause can not set field "+fieldName)
				   this.setValueByMethod(reflectValue,fieldType,fieldName,this.store.GetProperty(tag))
			   }

		   }
	   }


   }else {

   }
 }

 func (this *ValuesHolder)setValueByMethod(v reflect.Value,ft reflect.Type,fieldName string,value string){
	 methodName := "Set" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
	 rm:=v.MethodByName(methodName)
	 if rm.IsValid(){
		 bv:=BingoString(value)
	 	 var rv reflect.Value
		 switch(ft.Kind()) {
		 case reflect.String : rv=reflect.ValueOf(value)
		 case reflect.Bool:
			 vbool,_:=bv.Bool()
			 rv=reflect.ValueOf(vbool)
		 case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			 vint,_:=bv.Int64()
			 rv=reflect.ValueOf(vint)
		 case reflect.Float64, reflect.Float32:
			 vfloat,_:=bv.Float64()
			 rv=reflect.ValueOf(vfloat)
		 }

		 rm.Call([]reflect.Value{rv})
	 }
 }
 //根据类型将字符串转换成指定的值进行设置
func (this *ValuesHolder)setValue(v reflect.Value,t reflect.Type,value string){
	bv:=BingoString(value)
	switch(t.Kind()) {
	case reflect.String : v.SetString(value)
	case reflect.Bool:
		vbool,_:=bv.Bool()
		v.SetBool(vbool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vint,_:=bv.Int64()
		v.SetInt(vint)
	case reflect.Float64, reflect.Float32:
		vfloat,_:=bv.Float64()
		v.SetFloat(vfloat)
	}

}

func IsStructPtr(v interface{}) bool {
	t:=reflect.TypeOf(v)
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}




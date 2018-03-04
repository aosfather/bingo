package utils

import (
	"net/url"
	"reflect"
	"time"
)

const (
	_TAG_FIELD = "Field"
)

func IsMap(obj interface{}) bool {
	objT := reflect.TypeOf(obj)
	if objT.Kind() == reflect.Map {
		return true
	}
	return false
}
func HasFieldofStruct(obj interface{}, fieldName string) bool {

	_, rv, err :=GetStructTypeValue(obj)
	if err != nil {
		return false
	}
	//	rv := reflect.ValueOf(obj)
	val := rv.FieldByName(fieldName)
	return val.IsValid()
}
func GetRealType(obj interface{}) reflect.Type {
	objT := reflect.TypeOf(obj)
	if objT.Kind() == reflect.Ptr {
		return objT.Elem()
	}
	return objT
}

func CreateObjByType(t reflect.Type) interface{} {
	return reflect.New(t).Interface()
}

func GetStructTypeValue(obj interface{}) (reflect.Type, reflect.Value, error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case objT.Kind() == reflect.Struct:
	case objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct:
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		return objT, objV, CreateError(500, "Must set a struct or a struct pointer!")

	}
	return objT, objV, nil
}

func FillStructByForm(values url.Values, target interface{}) error {
	objT, objV, err := GetStructTypeValue(target)
	if err != nil {
		return err
	}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		colName := getFormFieldName(f)
		v := values.Get(colName)
		if v != "" && vf.CanSet() {
			setFieldValue(vf, v)
		}
	}
	return nil
}

func FillStruct(values map[string]interface{}, target interface{}) error {
	objT, objV, err := GetStructTypeValue(target)
	if err != nil {
		return err
	}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		colName := GetColName(f)
		v := values[colName]
		if v != nil && vf.CanSet() {
			setFieldValue(vf, v)
		}
	}

	return nil

}

func getFormFieldName(field reflect.StructField) string {
	//	colName := strings.ToLower(field.Name)
	//	tagField := field.Tag.Get(_TAG_FIELD)
	//	if tagField != "" {
	//		colName = tagField
	//	}
	//	return colName
	return GetColName(field)
}

func GetColName(field reflect.StructField) string {
	colName := field.Tag.Get(_TAG_FIELD)
	if colName == "" {
		colName = BingoString(field.Name).SnakeString()
	}

	return colName
}

func setFieldValue(ind reflect.Value, value interface{}) {
	var val interface{}
	objT := reflect.TypeOf(value)
	if objT.Kind() == reflect.Ptr {

		val = reflect.ValueOf(value).Elem().Interface()
	} else {
		val = value
	}

	switch ind.Kind() {
	case reflect.Bool:
		setBool(ind, val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setInt(ind, val)
	case reflect.String:
		setString(ind, val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setUint(ind, val)
	case reflect.Float64, reflect.Float32:
		setFloat(ind, val)
	case reflect.Struct:
		if value == nil {
			ind.Set(reflect.Zero(ind.Type()))

		} else if _, ok := ind.Interface().(time.Time); ok {
			setTime(ind, val)
		}
	}

}
func setString(ind reflect.Value, value interface{}) {
	if value == nil {
		ind.SetString("")
	} else {
		ind.SetString(ToStr(value))
	}
}

func setBool(ind reflect.Value, value interface{}) {
	if value == nil {
		ind.SetBool(false)
	} else if v, ok := value.(bool); ok {
		ind.SetBool(v)
	} else {
		v, _ := BingoString(ToStr(value)).Bool()
		ind.SetBool(v)
	}
}

func setInt(ind reflect.Value, value interface{}) {
	if value == nil {
		ind.SetInt(0)
	} else {
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ind.SetInt(val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ind.SetInt(int64(val.Uint()))
		default:
			v, _ := BingoString(ToStr(value)).Int64()
			ind.SetInt(v)
		}
	}
}

func setUint(ind reflect.Value, value interface{}) {
	if value == nil {
		ind.SetUint(0)
	} else {
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ind.SetUint(uint64(val.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ind.SetUint(val.Uint())
		default:
			v, _ := BingoString(ToStr(value)).Uint64()
			ind.SetUint(v)
		}
	}
}

func setFloat(ind reflect.Value, value interface{}) {
	if value == nil {
		ind.SetFloat(0)
	} else {
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Float64:
			ind.SetFloat(val.Float())
		default:
			v, _ := BingoString(ToStr(value)).Float64()
			ind.SetFloat(v)
		}
	}
}

func setTime(ind reflect.Value, value interface{}) {

}

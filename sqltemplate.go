/**

SQL生成模板
用于根据struct的定义生成update语句和insert语句
规则
  1、表名
   默认t_stuct的名称(小写，驼峰转 xx_xx)
   可以通过setTablePreFix来指定默认的表前缀
   可以通过tag Table:""指定表名

  2、字段名称
    默认字段名称(小写，驼峰转 xx_xx)
	可以通过 tag Field:""指定字段名

  3、特例
    通过tag Option:"" 来指定。
	可选值：auto、pk、not 分别表示 自动增长、主健、忽略

*/
package bingo

import (
	"reflect"
	"strings"
)

var (
	table_prefix       string = "t_"
	table_insert_cache map[string]string
)

func init() {
	table_insert_cache = make(map[string]string)
}

func SetTablePreFix(pfix string) {
	table_prefix = pfix
}

func GetInsertSql(target interface{}) (string, []interface{}, error) {
	objT, _, err := getStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	key := objT.Name()
	sql := table_insert_cache[key]

	if sql == "" {
		sql, args, err := CreateInserSql(target)
		if err != nil {
			return "", nil, err
		}
		table_insert_cache[key] = sql
		return sql, args, err

	}

	args, err := structValueToArray(target)
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func StructValueToCustomArray(target interface{}, col ...string) ([]interface{}, error) {
	_, objV, err := getStructTypeValue(target)
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, len(col))
	for i, field := range col {
		vf := objV.FieldByName(field)
		if !vf.CanInterface() {
			args[i] = nil
		} else {
			args[i] = objV.FieldByName(field).Interface()
		}

	}
	return args, nil
}
func structValueToArray(target interface{}) ([]interface{}, error) {
	objT, objV, err := getStructTypeValue(target)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0, 0)
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		//对于自增长和明确忽略的字段不做转换
		if isFieldIgnore(f) {
			continue
		}

		args = append(args, vf.Interface())

	}
	return args, nil

}

func CreateInserSql(target interface{}) (string, []interface{}, error) {
	objT, objV, err := getStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlFields string
	var sqlValues string
	args := make([]interface{}, 0, 0)
	fieldIndex := 0
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}
		tagTable := f.Tag.Get("Table")
		if tagTable != "" {
			tagTableName = tagTable
		}

		colName := getColName(f)
		//对于自增长和明确忽略的字段不做转换
		tagOption := f.Tag.Get("Option")
		if tagOption != "" {
			if strings.Index(tagOption, "auto") != -1 || strings.Index(tagOption, "not") != -1 {
				continue
			}
		}

		args = append(args, vf.Interface())
		if fieldIndex > 0 {
			sqlFields += ","
			sqlValues += ","
		}

		sqlFields += colName
		sqlValues += "?"
		fieldIndex++

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = table_prefix + BingoString(objT.Name()).SnakeString()
	}

	return "Insert into " + tagTableName + "(" + sqlFields + ") Values(" + sqlValues + ")", args, nil
}

func CreateUpdateSql(target interface{}, col ...string) (string, []interface{}, error) {
	objT, objV, err := getStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlFields string
	var sqlwheres string
	args := make([]interface{}, 0, 0)
	argsWhere := make([]interface{}, 0, 0)
	fieldIndex := 0
	whereFields := 0
	if len(col) != 0 {
		for _, fieldName := range col {
			f, b := objT.FieldByName(fieldName)
			vf := objV.FieldByName(fieldName)
			if !b || !vf.CanInterface() {
				continue
			}
			colName := getColName(f)
			if fieldIndex > 0 {
				sqlFields += ","
			}
			sqlFields += colName + "=?"
			args = append(args, vf.Interface())
			fieldIndex++
		}
	}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}
		tagTable := f.Tag.Get("Table")
		if tagTable != "" {
			tagTableName = tagTable
		}

		colName := getColName(f)
		if len(col) == 0 {

			if fieldIndex > 0 {
				sqlFields += ","
			}
			sqlFields += colName + "=?"
			args = append(args, vf.Interface())
			fieldIndex++
		}

		//对于自增长和明确忽略的字段不做转换
		tagOption := f.Tag.Get("Option")
		if tagOption != "" {
			if strings.Index(tagOption, "pk") != -1 {
				//where的处理
				if whereFields > 0 {
					sqlwheres += " and "
				}
				sqlwheres += colName + " =?"
				whereFields++
				argsWhere = append(argsWhere, vf.Interface())
			}
		}

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = table_prefix + BingoString(objT.Name()).SnakeString()
	}

	args = append(args, argsWhere...)

	return "update " + tagTableName + " set " + sqlFields + " where " + sqlwheres, args, nil
}

func CreateQuerySql(target interface{}, col ...string) (string, []interface{}, error) {
	objT, objV, err := getStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlwheres string
	argsWhere := make([]interface{}, 0, 0)
	whereFields := 0
	if len(col) != 0 {
		for _, fieldName := range col {
			f, b := objT.FieldByName(fieldName)
			vf := objV.FieldByName(fieldName)
			if !b || !vf.CanInterface() {
				continue
			}
			colName := getColName(f)
			if whereFields > 0 {
				sqlwheres += " and "
			}
			sqlwheres += colName + "=?"
			argsWhere = append(argsWhere, vf.Interface())
			whereFields++
		}
	}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}
		tagTable := f.Tag.Get("Table")
		if tagTable != "" {
			tagTableName = tagTable
		}
		if len(col) == 0 {
			colName := getColName(f)
			//对于标识为pk的字段做为条件
			tagOption := f.Tag.Get("Option")
			if tagOption != "" {
				if strings.Index(tagOption, "pk") != -1 {
					//where的处理
					if whereFields > 0 {
						sqlwheres += " and "
					}
					sqlwheres += colName + " =?"
					argsWhere = append(argsWhere, vf.Interface())
					whereFields++
				}
			}
		}

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = table_prefix + BingoString(objT.Name()).SnakeString()
	}

	return "select * from " + tagTableName + " where " + sqlwheres, argsWhere, nil
}

func isFieldIgnore(field reflect.StructField) bool {
	//对于自增长和明确忽略的字段不做转换
	tagOption := field.Tag.Get("Option")
	if tagOption != "" {
		if strings.Index(tagOption, "auto") != -1 || strings.Index(tagOption, "not") != -1 {
			return true
		}
	}
	return false
}

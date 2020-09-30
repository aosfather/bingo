package main

import (
	"github.com/aosfather/bingo_mvc/sqltemplate"
	"github.com/aosfather/bingo_utils/lua"
	l "github.com/yuin/gopher-lua"
)

func CreateDataBasseLib(db *sqltemplate.DataSource) map[string]l.LGFunction {
	luadb := LuaDataBase{DB: db}
	luadb.Init()
	libs := make(map[string]l.LGFunction)
	libs["query"] = luadb.lua_query_all
	libs["queryByPage"] = luadb.lua_query_page
	libs["insert"] = luadb.lua_insert
	libs["update"] = luadb.lua_update
	libs["find"] = luadb.lua_find
	libs["delete"] = luadb.lua_delete
	return libs

}

/**
  给lua提供的数据库相关操作
*/
type LuaDataBase struct {
	DB     *sqltemplate.DataSource
	caches map[string]*sqltemplate.MapperDao
}

func (this *LuaDataBase) Init() {
	this.caches = make(map[string]*sqltemplate.MapperDao)
}

func (this *LuaDataBase) getMapperDao(name string) *sqltemplate.MapperDao {
	if this.DB == nil {
		panic("not set Database source!")
	}
	mapper := this.caches[name]
	if mapper == nil {
		mapper = this.DB.GetMapperDao(name)
		if mapper != nil {
			this.caches[name] = mapper
		}
	}
	return mapper
}

/**
  获取单条数据
  find(name,id,table) (table,bool)
  name:
*/
func (this *LuaDataBase) lua_find(l *l.LState) int {
	parames := l.Get(-1)
	l.Pop(1)
	id := l.Get(-1).String()
	l.Pop(1)
	name := l.Get(-1).String()
	l.Pop(1)
	mapper := this.getMapperDao(name)
	if mapper != nil {
		p := lua.ToGoMap(parames)
		result := mapper.Find(&p, id)
		l.Push(lua.ToLuaTable2(l, p))
		l.Push(lua.ToLuaValue(result))
	} else {
		l.Push(lua.ToLuaValue("not found the mapper name:" + name))
		l.Push(lua.ToLuaValue(false))
	}

	return 2
}

/**
  查询获取所有数据
  query(name,id,table) (table,error)
  name:
*/
func (this *LuaDataBase) lua_query_all(l *l.LState) int {
	parames := l.Get(-1)
	l.Pop(1)
	id := l.Get(-1).String()
	l.Pop(1)
	name := l.Get(-1).String()
	l.Pop(1)
	var error string
	mapper := this.getMapperDao(name)
	if mapper == nil {
		error = "not found the mapper name:" + name
	}
	result := mapper.QueryAll(lua.ToGoMap(parames), id)
	l.Push(lua.ArrayToLuaTable(l, result))
	l.Push(lua.ToLuaValue(error))
	return 2
}

/**
  查询指定page的数据数据
  query(name,id,table,pageNo,size) (table,error)
  name:
*/
func (this *LuaDataBase) lua_query_page(l *l.LState) int {
	pageSize := l.ToInt(-1)
	l.Pop(1)
	pageNo := l.ToInt(-1)
	l.Pop(1)
	parames := l.Get(-1)
	l.Pop(1)
	id := l.Get(-1).String()
	l.Pop(1)
	name := l.Get(-1).String()
	l.Pop(1)
	var error string
	mapper := this.getMapperDao(name)
	if mapper == nil {
		error = "not found the mapper name:" + name
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNo < 0 {
		pageNo = 0
	}
	result := mapper.Query(lua.ToGoMap(parames), sqltemplate.Page{Size: pageSize, Index: pageNo}, id)
	l.Push(lua.ArrayToLuaTable(l, result))
	l.Push(lua.ToLuaValue(error))
	return 2
}

/**
  插入数据
  insert(name,id,table) (intid,bool,msg)
  name:
*/
func (this *LuaDataBase) lua_insert(l *l.LState) int {
	parames := l.Get(-1)
	l.Pop(1)
	id := l.Get(-1).String()
	l.Pop(1)
	name := l.Get(-1).String()
	l.Pop(1)
	mapper := this.getMapperDao(name)
	if mapper != nil {
		p := lua.ToGoMap(parames)
		result, e := mapper.Insert(p, id)
		if e != nil {
			l.Push(lua.ToLuaValue(-1))
			l.Push(lua.ToLuaValue(false))
			l.Push(lua.ToLuaValue(e.Error()))
			errs("lua insert db error:", e.Error())
		} else {
			l.Push(lua.ToLuaValue(result))
			l.Push(lua.ToLuaValue(true))
			l.Push(lua.ToLuaValue("ok"))
		}

	} else {
		l.Push(lua.ToLuaValue(-1))
		l.Push(lua.ToLuaValue(false))
		l.Push(lua.ToLuaValue("not found the mapper name:" + name))
		errs("lua insert db error:", "not found the mapper name:", name)
	}

	return 3
}

/**
  更新数据
  update(name,id,table) (int,msg)
  name:
*/
func (this *LuaDataBase) lua_update(l *l.LState) int {
	parames := l.Get(-1)
	l.Pop(1)
	id := l.Get(-1).String()
	l.Pop(1)
	name := l.Get(-1).String()
	l.Pop(1)
	mapper := this.getMapperDao(name)
	if mapper != nil {
		p := lua.ToGoMap(parames)
		result, e := mapper.Update(p, id)
		if e != nil {
			l.Push(lua.ToLuaValue(-1))
			l.Push(lua.ToLuaValue(e.Error()))
			errs("lua update db error:", e.Error())
		} else {
			l.Push(lua.ToLuaValue(result))
			l.Push(lua.ToLuaValue("ok"))
		}

	} else {
		l.Push(lua.ToLuaValue(-1))
		l.Push(lua.ToLuaValue("not found the mapper name:" + name))
		errs("lua update db error:", "not found the mapper name:", name)
	}

	return 2
}

/**
  删除数据
  delete(name,id,table) (int,msg)
  name:
*/
func (this *LuaDataBase) lua_delete(l *l.LState) int {
	parames := l.Get(-1)
	l.Pop(1)
	id := l.Get(-1).String()
	l.Pop(1)
	name := l.Get(-1).String()
	l.Pop(1)
	mapper := this.getMapperDao(name)
	if mapper != nil {
		p := lua.ToGoMap(parames)
		result, e := mapper.Delete(p, id)
		if e != nil {
			l.Push(lua.ToLuaValue(-1))
			l.Push(lua.ToLuaValue(e.Error()))
			errs("lua delete db error:", e.Error())
		} else {
			l.Push(lua.ToLuaValue(result))
			l.Push(lua.ToLuaValue("ok"))
		}

	} else {
		l.Push(lua.ToLuaValue(-1))
		l.Push(lua.ToLuaValue("not found the mapper name:" + name))
		errs("lua delete db error:", "not found the mapper name:", name)
	}

	return 2
}

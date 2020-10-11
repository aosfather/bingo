package main

import (
	lua "github.com/yuin/gopher-lua"
	"os"
)

/**
  系统库
*/
var syslibs map[string]lua.LGFunction

func init() {
	syslibs = make(map[string]lua.LGFunction)
	syslibs["getenv"] = lua_env_get
	syslibs["topassword"] = lua_env_topassword

}

//获取系统环境变量
func lua_env_get(l *lua.LState) int {
	name := l.ToString(1)

	if name == "" {
		l.Push(lua.LString("not found name"))
		l.Push(lua.LFalse)

	} else {
		value := os.Getenv(name)
		l.Push(lua.LString(value))
		l.Push(lua.LTrue)
	}

	return 2
}

func lua_env_topassword(l *lua.LState) int {
	name := l.ToString(1)

	if name == "" {
		l.Push(lua.LString(""))
	} else {
		value := topasswords(name)
		l.Push(lua.LString(value))
	}

	return 1
}

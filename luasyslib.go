package main

import (
	lua "github.com/yuin/gopher-lua"
	"os"
)

/**
  系统库
*/

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

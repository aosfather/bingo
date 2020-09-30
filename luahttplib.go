package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo_utils/lua"
	l "github.com/yuin/gopher-lua"
	"net/url"
	"strings"
)

/**
  提供给lua用的apis
*/
var httplibs map[string]l.LGFunction

func init() {
	httplibs = make(map[string]l.LGFunction)
	httplibs["get"] = lua_http_get
	httplibs["post"] = lua_http_post

}

/**
  http请求，http.get(url,headers)
  返回 response 和 错误信息
*/
func lua_http_get(l *l.LState) int {
	_headers := l.Get(-1)
	l.Pop(1)
	url := l.Get(-1).String()
	l.Pop(1)
	buffer := new(bytes.Buffer)
	err := doHttpRequest("GET", url, "", buffer, http_headers(_headers))
	if err != nil {
		l.Push(lua.ToLuaValue(""))
		l.Push(lua.ToLuaValue(err.Error()))
	} else {
		l.Push(lua.ToLuaValue(buffer.String()))
		l.Push(lua.ToLuaValue(""))
	}
	return 2
}

/**
  http请求 http.post(url,headers,body)
  返回 response 和 错误信息
*/
func lua_http_post(l *l.LState) int {
	body := l.Get(-1)
	l.Pop(1)

	_headers := l.Get(-1)
	l.Pop(1)

	url := l.Get(-1).String()
	l.Pop(1)

	buffer := new(bytes.Buffer)
	headers := http_headers(_headers)
	err := doHttpRequest("POST", url, http_tobody(body, headers), buffer, headers)
	if err != nil {
		l.Push(lua.ToLuaValue(""))
		l.Push(lua.ToLuaValue(err.Error()))
	} else {
		l.Push(lua.ToLuaValue(buffer.String()))
		l.Push(lua.ToLuaValue(""))
	}
	return 2
}

func http_tobody(v l.LValue, headers map[string]string) string {
	switch v.Type() {
	case l.LTTable:
		contenttype, ok := headers["Content-Type"]
		m := lua.ToGoMap(v)
		//使用json方式则转换table格式为json
		if ok && strings.Index(contenttype, "application/json") >= 0 {

			return toJson(m)
		}
		//如果没有设置类型，默认设置成form-urlencoded格式
		if !ok {
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
		//默认http www form方式
		buffer := new(bytes.Buffer)
		first := true
		for key, value := range m {
			if first {
				first = false
			} else {
				buffer.WriteString("&")
			}
			buffer.WriteString(url.QueryEscape(key))
			buffer.WriteString("=")
			buffer.WriteString(url.QueryEscape(fmt.Sprintf("%v", value)))
		}
		return buffer.String()

	default:
		return v.String()
	}

}

func toJson(m interface{}) string {
	result, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(result)
}

func http_headers(v l.LValue) map[string]string {
	goheaders := lua.ToGoMap(v)
	headers := make(map[string]string)
	for k, v := range goheaders {
		headers[k] = fmt.Sprintf("%v", v)
	}
	return headers
}

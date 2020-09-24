package main

import (
	"bytes"
	"github.com/aosfather/bingo_utils/lua"
	l "github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"
)

/**
   表单action
  支持的类型：SQL、GET、POST、LUA
*/

//表单action
type FormActions struct {
	pool *lua.LuaPool
}

func (this *FormActions) Init() {
	//设置lua引擎的lib脚本查找路径。
	lua.SetLuaPath(".\\libs\\lua")

	//加载自定义库
	libs := make(map[string]l.LGFunction)
	//加载lib库

	this.pool = lua.NewLuaPool(100, "bingo", libs)

}

func (this *FormActions) getLuaScript(code string, content string) *lua.LuaScript {
	script := &lua.LuaScript{Log: this.lualog}
	script.SetPool(this.pool)
	script.Load(code, content)
	return script
}

func (this *FormActions) lualog(str string) {
	debug("lua:", str)
}

func (this *FormActions) Execute(meta *FormMeta, parameter map[string]interface{}) (interface{}, error) {
	switch meta.ScriptType {
	case "GET":
		headers, body := this.processHttpScript(meta, parameter)
		return this.doGet(meta.Extends["url"], headers, strings.TrimSpace(body))
	case "POST":
		headers, body := this.processHttpScript(meta, parameter)
		return this.doPost(meta.Extends["url"], headers, strings.TrimSpace(body))
	case "LUA":
		return this.processLuaScript(meta.Code, meta.Script, parameter)
	case "SQL":
	}

	return nil, nil
}

//处理lua脚本，将lua脚本作为整个的运行逻辑
func (this *FormActions) processLuaScript(name, script string, parameter map[string]interface{}) (interface{}, error) {
	luascript := this.getLuaScript(name, script)
	beforeRun := func(l *l.LState) {
		l.SetGlobal("_inputs", lua.SetReadOnly(l, lua.ToLuaTable2(l, parameter)))
	}
	afterRun := func(l *l.LState) {

	}
	return luascript.Call(beforeRun, afterRun)
}

//处理http请求
func (this *FormActions) processHttpScript(meta *FormMeta, parameter map[string]interface{}) (map[string]string, string) {
	t := template.New(meta.Code)
	_, err := t.Parse(meta.Script)
	if err != nil {
		errs("parse template error!", err.Error())
	} else {
		headers := make(map[string]string)
		buffer := new(bytes.Buffer)
		err = t.Execute(buffer, parameter)
		if err != nil {
			errs("execute template error!", err.Error())
		}
		return headers, buffer.String()
	}

	return nil, ""

}

func (this *FormActions) doGet(url string, headers map[string]string, body string) (string, error) {
	debug("header:", headers)
	buffer := new(bytes.Buffer)
	err := doHttpRequest("GET", url, body, buffer, headers)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (this *FormActions) doPost(url string, headers map[string]string, body string) (string, error) {
	if _, ok := headers["Content-Type"]; !ok {
		if strings.Index(body, "{") >= 0 {
			headers["Content-Type"] = "application/json;charset=utf-8"
		} else {
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	debug("header:", headers)
	buffer := new(bytes.Buffer)
	err := doHttpRequest("POST", url, body, buffer, headers)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil

}

func (this *FormActions) doSQL() {

}

//网络访问超时设置
const _ClientTimeout = 20 * time.Second

func doHttpRequest(method string, url string, content string, writer io.Writer, headers map[string]string) error {
	//post
	c := &http.Client{Timeout: _ClientTimeout}
	times := 1

DO_POST:
	req, err := http.NewRequest(method, url, strings.NewReader(string(content)))
	if err != nil {
		return err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	debug("values:", content)
	resp, err := c.Do(req)
	if err != nil {
		errs(err.Error())
		if times < 3 {
			times++
			time.Sleep(time.Second)
			debug("try the ", times, " times!")
			goto DO_POST
		} else {
			return err
		}
	}

	defer resp.Body.Close()
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		errs(err.Error())
		return err
	}
	return nil
}

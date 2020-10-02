package main

import (
	"fmt"
	"io"
)

var verifys = make(map[string]string)
var renders = make(map[string]FormItemEditorRender)

//item渲染，返回校验脚本和额外渲染脚本
type FormItemEditorRender func(input Parameter, w io.Writer) (string, string)

func init() {
	renders["String"] = textRender
	renders["MobileNo"] = phoneRender
	renders["Email"] = emailRender
	renders["Date"] = dateRender
	renders["DateTime"] = datetimeRender
	renders["Enum"] = enumRender
}

func AddVerify(key, value string) {
	if key != "" && value != "" {
		verifys[key] = value
	}
}
func textRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required"))
	}
	script := ""
	if input.Verify != "" {
		if input.Policy == "Must" {
			w.Write([]byte("|"))
		}
		w.Write([]byte(input.Verify))
		script = fmt.Sprintf(",%s:%s", input.Verify, verifys[input.Verify])
	}
	tip := input.InputTip
	if tip == "" {
		tip = "请输入"
	}

	w.Write([]byte(fmt.Sprintf(`" autocomplete="off" placeholder="%s" class="layui-input">`, tip)))
	return script, ""
}

func textAreaRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(`<textarea `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required"))
	}

	tip := input.InputTip
	if tip == "" {
		tip = "请输入"
	}
	w.Write([]byte(fmt.Sprintf(`" autocomplete="off" placeholder="%s" class="layui-textarea" rows="5"></textarea>`, tip)))
	return "", ""
}

func phoneRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`phone" autocomplete="off" placeholder="请输入手机号" class="layui-input">`))
	return "", ""
}

func emailRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`email" autocomplete="off" placeholder="请输入邮箱" class="layui-input">`))
	return "", ""
}

func dateRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(`name="%s" id="%s"`, input.Name, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`date" autocomplete="off" placeholder="yyyy-MM-dd" class="layui-input">`))
	return "", fmt.Sprintf("laydate.render({elem: '#%s'});", input.Name)
}

func datetimeRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(`name="%s" id="%s"`, input.Name, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`datetime" autocomplete="off" placeholder="yyyy-MM-dd HH:mm:ss" class="layui-input">`))
	return "", fmt.Sprintf("laydate.render({elem: '#%s',type: 'datetime'});", input.Name)
}

func enumRender(input Parameter, w io.Writer) (string, string) {
	w.Write([]byte(fmt.Sprintf(`<select name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required"))
	}
	w.Write([]byte(`" >`))
	//输出字典
	dict := GetDict(input.Expr)
	if dict.Code != "" {
		for _, item := range dict.Items {
			w.Write([]byte(fmt.Sprintf(`<option value="%s">%s</option>`, item.Code, item.Label)))
		}
	}

	w.Write([]byte(" </select>"))

	return "", ""
}

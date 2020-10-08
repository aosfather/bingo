package main

import (
	"fmt"
	"github.com/aosfather/bingo_mvc/dd"
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
	renders["text"] = textAreaRender
}

func AddVerify(key, value string) {
	if key != "" && value != "" {
		verifys[key] = value
	}
}

func writeAttrible(att string, value string, w io.Writer) {
	w.Write([]byte(att))
	w.Write([]byte(`="`))
	w.Write([]byte(value))
	w.Write([]byte(`" `))
}

func writeString(attstr string, w io.Writer) {
	w.Write([]byte(attstr))
}

//预制的组件
func preTextinput(input Parameter, verify string, plachholder string, w io.Writer) {
	writeString(`<input type="text" `, w)
	writeAttrible("name", input.Name, w)
	writeAttrible("id", input.Name, w)
	if input.Policy == "Must" {
		verify = "required|" + verify
	}
	writeAttrible("lay-verify", verify, w)
	writeString(` autocomplete="off" `, w)
	writeAttrible("placeholder", plachholder, w)
	if input.Readonly {
		writeString(` readonly="true"`, w)
	}
	writeString(` class="layui-input">`, w)

}
func textRender(input Parameter, w io.Writer) (string, string) {
	writeString(`<input type="text" `, w)
	writeAttrible("name", input.Name, w)
	verify := ""
	if input.Policy == "Must" {
		verify = "required"
	}
	script := ""
	if input.Verify != "" {
		if input.Policy == "Must" {
			verify = verify + "|"
		}
		verify = verify + input.Verify
		script = fmt.Sprintf(",%s:%s", input.Verify, verifys[input.Verify])
	}
	writeAttrible("lay-verify", verify, w)
	tip := input.InputTip
	if tip == "" {
		tip = "请输入"
	}
	writeString(`autocomplete="off"`, w)
	writeAttrible("placeholder", tip, w)
	if input.Readonly {
		writeString(` readonly="true"`, w)
	}

	writeString(` class="layui-input">`, w)
	return script, ""
}

func textAreaRender(input Parameter, w io.Writer) (string, string) {
	writeString(`<textarea `, w)
	writeAttrible("name", input.Name, w)
	verify := ""
	if input.Policy == "Must" {
		verify = "required"
	}
	writeAttrible("lay-verify", verify, w)
	tip := input.InputTip
	if tip == "" {
		tip = "请输入"
	}
	writeString(`autocomplete="off"`, w)
	writeAttrible("placeholder", tip, w)
	if input.Readonly {
		writeString(` readonly="true"`, w)
	}

	writeString(` class="layui-textarea" rows="5"></textarea>`, w)
	return "", ""
}

func phoneRender(input Parameter, w io.Writer) (string, string) {
	preTextinput(input, "phone", "请输入手机号", w)
	return "", ""
}

func emailRender(input Parameter, w io.Writer) (string, string) {
	preTextinput(input, "email", "请输入邮箱", w)
	return "", ""
}

func dateRender(input Parameter, w io.Writer) (string, string) {
	preTextinput(input, "date", "yyyy-MM-dd", w)
	return "", fmt.Sprintf("laydate.render({elem: '#%s'});", input.Name)
}

func datetimeRender(input Parameter, w io.Writer) (string, string) {
	preTextinput(input, "datetime", "yyyy-MM-dd HH:mm:ss", w)
	return "", fmt.Sprintf("laydate.render({elem: '#%s',type: 'datetime'});", input.Name)
}

func enumRender(input Parameter, w io.Writer) (string, string) {
	writeString(`<select `, w)
	writeAttrible("name", input.Name, w)
	verify := ""
	if input.Policy == "Must" {
		verify = "required"
	}
	writeAttrible("lay-verify", verify, w)
	if input.Readonly {
		writeString(` readonly="true"`, w)
	}
	writeString(`>`, w)

	//输出字典
	dict := dd.GetDict(input.Expr)
	if dict.Code != "" {
		for _, item := range dict.Items {
			w.Write([]byte(fmt.Sprintf(`<option value="%s">%s</option>`, item.Code, item.Label)))
		}
	}
	writeString(` </select>`, w)

	return "", ""
}

package main

import (
	"fmt"
	"io"
)

var renders = make(map[string]FormItemEditorRender)

type FormItemEditorRender func(input Parameter, w io.Writer) string

func init() {
	renders["String"] = textRender
	renders["MobileNo"] = phoneRender
	renders["Email"] = emailRender
	renders["Date"] = dateRender
	renders["Enum"] = enumRender
}

func textRender(input Parameter, w io.Writer) string {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required"))
	}
	w.Write([]byte(`" autocomplete="off" placeholder="请输入" class="layui-input">`))
	return ""
}

func textAreaRender(input Parameter, w io.Writer) string {
	w.Write([]byte(`<textarea `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required"))
	}
	w.Write([]byte(`" autocomplete="off" placeholder="请输入" class="layui-textarea" rows="5"></textarea>`))
	return ""
}

func phoneRender(input Parameter, w io.Writer) string {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`phone" autocomplete="off" placeholder="请输入手机号" class="layui-input">`))
	return ""
}

func emailRender(input Parameter, w io.Writer) string {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`email" autocomplete="off" placeholder="请输入手机号" class="layui-input">`))
	return ""
}

func dateRender(input Parameter, w io.Writer) string {
	w.Write([]byte(`<input type="text" `))
	w.Write([]byte(fmt.Sprintf(` name="%s" `, input.Name)))
	w.Write([]byte(fmt.Sprintf(`name="%s" id="%s"`, input.Name, input.Name)))
	w.Write([]byte(`lay-verify="`))
	if input.Policy == "Must" {
		w.Write([]byte("required|"))
	}
	w.Write([]byte(`date" autocomplete="off" placeholder="yyyy-MM-dd" class="layui-input">`))
	return fmt.Sprintf("laydate.render({elem: '#%s'});", input.Name)
}

func enumRender(input Parameter, w io.Writer) string {
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

	return ""
}

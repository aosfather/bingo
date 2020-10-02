package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

//绘制引擎
type RenderEngine interface {
	Render(meta *FormMeta) ([]string, string, string)
	GetTemplate() string
	GetKeys() []string
}

//表单引擎
type FormEngine struct {
}

func (this *FormEngine) Render(meta *FormMeta) ([]string, string, string) {
	//表单体
	buffer := new(bytes.Buffer)
	scriptBuffer := new(bytes.Buffer)
	extendscriptBuffer := new(bytes.Buffer)
	for _, input := range meta.Parameters {
		s, es := renderItem(input, buffer, false)
		scriptBuffer.WriteString(s)
		extendscriptBuffer.WriteString(es)
	}

	return []string{buffer.String()}, scriptBuffer.String(), extendscriptBuffer.String()
}

func renderItem(input Parameter, w io.Writer, inline bool) (string, string) {
	//字段名
	var label string
	if inline {
		label = `<div class="layui-inline">
                <label class="layui-form-label" style="width:120px;">%s</label>
                 <div class="layui-input-inline"> `
	} else {
		label = ` <div class="layui-form-item">
                <label class="layui-form-label">%s</label>
                <div class="layui-input-block">`
	}

	w.Write([]byte(fmt.Sprintf(label, input.Label)))
	//字段输入控件
	fr := renders[input.Type]
	var script, extendscript string
	if fr != nil {
		script, extendscript = fr(input, w)
		debug("script:", script)
	} else {
		fr = renders["String"]
		script, extendscript = fr(input, w)
		debug("script:", script)
	}
	//字段结束
	w.Write([]byte("</div></div>"))
	return script, extendscript
}

func (this *FormEngine) GetTemplate() string {
	return "form"
}

func (this *FormEngine) GetKeys() []string {
	return []string{"FORM_FIELDS"}
}

//查询表单引擎
type QueryFormEngine struct {
}

func (this *QueryFormEngine) GetTemplate() string {
	return "query_form"
}

func (this *QueryFormEngine) GetKeys() []string {
	return []string{"FORM_FIELDS", "FORM_GRID"}
}

func (this *QueryFormEngine) Render(meta *FormMeta) ([]string, string, string) {

	//渲染表单查询体
	parameterBuffer := new(bytes.Buffer)
	script, exscript := this.renderQueryParameters(meta, parameterBuffer)
	//渲染表格
	gridBuffer := new(bytes.Buffer)
	this.renderQueryGrid(meta, gridBuffer)
	return []string{parameterBuffer.String(), gridBuffer.String()}, script, exscript
}

//渲染查询条件
func (this *QueryFormEngine) renderQueryParameters(meta *FormMeta, writer io.Writer) (string, string) {
	scriptBuffer := new(bytes.Buffer)
	extendscriptBuffer := new(bytes.Buffer)
	for index, input := range meta.Parameters {
		if index%4 == 0 {
			writer.Write([]byte("<div class=\"layui-form-item\">"))
		}
		s, es := renderItem(input, writer, true)
		scriptBuffer.WriteString(s)
		extendscriptBuffer.WriteString(es)
		if index%4 == 3 {
			writer.Write([]byte("</div>"))
		}
	}

	if (len(meta.Parameters)-1)%4 != 3 {
		writer.Write([]byte("</div>"))
	}
	return scriptBuffer.String(), extendscriptBuffer.String()
}

//渲染表格
func (this *QueryFormEngine) renderQueryGrid(meta *FormMeta, writer io.Writer) {
	for _, rs := range meta.ResultSet {
		writer.Write([]byte(fmt.Sprintf("<th lay-data=\"{field:'%s'}\">%s</th>", rs.Name, rs.Label)))
	}
	//如果有设置动作就处理
	if meta.Tools != nil && len(meta.Tools) > 0 {
		writer.Write([]byte("<th lay-data=\"{fixed: 'right', toolbar: '#tabletools', width:150,align:'center'}\"></th>"))

		writer.Write([]byte("<div><script type=\"text/html\" id=\"tabletools\">"))
		for _, tool := range meta.Tools {
			var conditons []string
			if len(tool.Condition) > 0 {
				for _, c := range tool.Condition {
					conditons = append(conditons, "d."+c)
				}
			}
			if len(conditons) > 0 {
				writer.Write([]byte(fmt.Sprintf("{{# if(%s){ }}", strings.Join(conditons, "||"))))
			}
			writer.Write([]byte(fmt.Sprintf("<button class=\"layui-btn layui-btn-sm\" lay-event=\"%s\">%s</button>", tool.Name, tool.Label)))
			if len(conditons) > 0 {
				writer.Write([]byte("{{# } }}"))
			}
		}
		writer.Write([]byte("</script></div>"))
	}

}

package main

import (
	"bytes"
	"fmt"
	"io"
)

//绘制引擎
type RenderEngine interface {
	Render(meta *FormMeta) ([]string, string)
	GetTemplate() string
	GetKeys() []string
}

//表单引擎
type FormEngine struct {
}

func (this *FormEngine) Render(meta *FormMeta) ([]string, string) {
	//表单体
	buffer := new(bytes.Buffer)
	scriptBuffer := new(bytes.Buffer)
	for _, input := range meta.Parameters {
		scriptBuffer.WriteString(renderItem(input, buffer, false))
	}

	return []string{buffer.String()}, scriptBuffer.String()
}

func renderItem(input Parameter, w io.Writer, inline bool) string {
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
	var script string
	if fr != nil {
		script = fr(input, w)
		debug("script:", script)
	} else {
		fr = renders["String"]
		script = fr(input, w)
		debug("script:", script)
	}
	//字段结束
	w.Write([]byte("</div></div>"))
	return script
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

func (this *QueryFormEngine) Render(meta *FormMeta) ([]string, string) {

	//渲染表单查询体
	parameterBuffer := new(bytes.Buffer)
	script := this.renderQueryParameters(meta, parameterBuffer)
	//渲染表格
	gridBuffer := new(bytes.Buffer)
	this.renderQueryGrid(meta, gridBuffer)
	return []string{parameterBuffer.String(), gridBuffer.String()}, script
}

//渲染查询条件
func (this *QueryFormEngine) renderQueryParameters(meta *FormMeta, writer io.Writer) string {
	scriptBuffer := new(bytes.Buffer)
	for index, input := range meta.Parameters {
		if index%4 == 0 {
			writer.Write([]byte("<div class=\"layui-form-item\">"))
		}

		scriptBuffer.WriteString(renderItem(input, writer, true))
		if index%4 == 3 {
			writer.Write([]byte("</div>"))
		}
	}

	if (len(meta.Parameters)-1)%4 != 3 {
		writer.Write([]byte("</div>"))
	}
	return scriptBuffer.String()
}

//渲染表格
func (this *QueryFormEngine) renderQueryGrid(meta *FormMeta, writer io.Writer) {
	for _, rs := range meta.ResultSet {
		writer.Write([]byte(fmt.Sprintf("<th lay-data=\"{field:'%s'}\">%s</th>", rs.Name, rs.Label)))
	}
}

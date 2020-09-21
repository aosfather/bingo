package main

import (
	"bytes"
	"fmt"
	"io"
)

//绘制引擎
type RenderEngine interface {
	Render(meta *FormMeta) []string
	GetTemplate() string
	GetKeys() []string
}

//表单引擎
type FormEngine struct {
}

func (this *FormEngine) Render(meta *FormMeta) []string {
	//表单体
	buffer := new(bytes.Buffer)
	for _, input := range meta.Parameters {
		renderItem(input, buffer)
	}

	return []string{buffer.String()}
}

func renderItem(input Parameter, w io.Writer) {
	//字段名
	label := ` <div class="layui-form-item">
                <label class="layui-form-label">%s</label>
                <div class="layui-input-block">`
	w.Write([]byte(fmt.Sprintf(label, input.Label)))
	//字段输入控件
	fr := renders[input.Type]
	if fr != nil {
		script := fr(input, w)
		debug("script:", script)
	} else {
		fr = renders["String"]
		script := fr(input, w)
		debug("script:", script)
	}
	//字段结束
	w.Write([]byte("</div></div>"))
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

func (this *QueryFormEngine) Render(meta *FormMeta) []string {

	//渲染表单查询体
	parameterBuffer := new(bytes.Buffer)
	this.renderQueryParameters(meta, parameterBuffer)
	//渲染表格

	return []string{parameterBuffer.String()}
}

//渲染查询条件
func (this *QueryFormEngine) renderQueryParameters(meta *FormMeta, writer io.Writer) {
	for index, input := range meta.Parameters {
		if index%4 == 0 {
			writer.Write([]byte("<div class=\"layui-form-item\">"))
		}

		renderItem(input, writer)
		if index%4 == 3 {
			writer.Write([]byte("</div>"))
		}
	}

	if (len(meta.Parameters)-1)%4 != 3 {
		writer.Write([]byte("</div>"))
	}
}

//渲染表格
func (this *QueryFormEngine) renderQueryGrid(meta *FormMeta, writer io.Writer) {

}

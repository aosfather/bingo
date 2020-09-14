package main

import (
	"fmt"
	"io"
)

//绘制引擎
type RenderEngine interface {
	Render(meta *FormMeta, writer io.Writer)
	GetTemplate() string
}

//表单引擎
type FormEngine struct {
}

func (this *FormEngine) Render(meta *FormMeta, writer io.Writer) {
	//表单体
	for _, input := range meta.Parameters {
		this.renderItem(input, writer)
	}
}

func (this *FormEngine) renderItem(input Parameter, w io.Writer) {
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

//查询表单引擎
type QueryFormEngine struct {
}

package main

import "io"

//绘制引擎
type RenderEngine interface {
	Render(meta *FormMeta, writer io.Writer)
	GetTemplate() string
}

//表单引擎
type FormEngine struct {
}

func (this *FormEngine) Render(meta *FormMeta, writer io.Writer) {

}

func (this *FormEngine) GetTemplate() string {
	return "formTemplate"
}

//查询表单引擎
type QueryFormEngine struct {
}

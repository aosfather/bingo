package main

import "io"

/**
  界面渲染
  通过模板进行渲染
  输入：input，script
  输出：html片段
  通过yaml配置关系
*/
type TemplateItemEditorRender struct {
}

func (this *TemplateItemEditorRender) Render(input Parameter, w io.Writer) (string, string) {
	return "", ""
}

package bingo

import (
	"html/template"
	"io"
)

/*
 模板的实现
 实现特性
  1、片段库的定义
  2、进行缓存常用的模板对象进行加速
  3、监视模板文件的变化，当模板文件变化后，摧毁缓存的模板对象
*/

type TemplateEngine struct {
	RootPath                string //模板根路径
	SubTemplatePath         string //子模板及片段定义的目录
	Suffix                  string //模板文件后缀
	CacheSize               int    //缓存模板个数
	ErrorTemplate           string //错误模板
	useDefaultErrorTemplate bool   //是否使用默认错误模板
}

func (this *TemplateEngine) Render(w io.Writer, templateName string, data interface{}) {
	//use cache

	//执行template
	e := this.writeTemplate(w, templateName, data)
	if e != nil {
		this.writeError(w, e)
	}
}

func (this *TemplateEngine) writeTemplate(w io.Writer, templateName string, data interface{}) BingoError {
	var templateError BingoError = nil
	templateFile := this.getRealPath(templateName)
	if isFileExist(templateFile) {
		tmpl, err := template.New(templateName).ParseFiles(templateFile)
		if err != nil {
			templateError = CreateError(500, "template load error:"+err.Error())

		} else {
			err = tmpl.Execute(w, data)
			if err != nil {
				templateError = CreateError(500, "template render error:"+err.Error())

			}
		}
	} else {
		templateError = CreateError(404, "View file '"+templateName+"' not found")
	}
	return templateError
}

func (this *TemplateEngine) writeError(w io.Writer, err BingoError) {
	if this.useDefaultErrorTemplate {
		tmpl, _ := template.New("error").Parse("<html><body><h1>{{.Code}}</h1><h3>{{.Error}}</h3></body></html>")
		tmpl.Execute(w, err)
	} else {
		e := this.writeTemplate(w, this.ErrorTemplate, err)
		if e != nil {
			this.useDefaultErrorTemplate = true
			this.writeError(w, err)
		}
	}

}

func (this *TemplateEngine) getRealPath(templateName string) string {
	return this.RootPath + "/" + templateName
}

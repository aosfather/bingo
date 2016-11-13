package bingo

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*
默认返回转换器
1、根据返回类型来进行转换
2、ModelView-> 走template转换
3、其它类型->走json
4、文件流的支持？
5、xml的支持?
6、图片?
*/
type defaultResponseConverter struct {
	templateDir string
}

func (this *defaultResponseConverter) setTemplateDir(dir string) {
	this.templateDir = dir
}

func (this *defaultResponseConverter) Convert(writer http.ResponseWriter, obj interface{}) {
	if mv, ok := obj.(ModelView); ok {
		writer.Header().Add(_CONTENT_TYPE, _CONTENT_HTML)
		this.writeWithTemplate(writer, mv.View, mv.Model)
	} else if rv, ok := obj.(StaticView); ok { //静态资源处理
		writeUseFile(writer, rv)

	} else {
		writer.Header().Add(_CONTENT_TYPE, _CONTENT_JSON)
		writeUseJson(writer, obj)
	}
}

func (this *defaultResponseConverter) writeWithTemplate(writer http.ResponseWriter, templateName string, obj interface{}) {
	var tmpl *template.Template
	var err error
	templatefile := this.templateDir + "/" + templateName
	if isFileExist(templatefile) {
		tmpl, err = template.New(templateName).ParseFiles(templatefile)
		if err != nil {
			panic(err)
		}
	} else {
		if _, ok := obj.(BingoError); !ok {
			obj = CreateError(500, "View file '"+templateName+"' not found")
			this.writeWithTemplate(writer, "error", obj)
			return
		}

		tmpl, err = template.New(templateName).Parse("<html><body><h1>{{.Code}}</h1><h3>{{.Error}}</h3></body></html>")
	}
	err = tmpl.Execute(writer, obj)
	if err != nil {
		panic(err)
	}
}
func writeUseFile(writer http.ResponseWriter, rv StaticView) {
	writer.Header().Add(_CONTENT_TYPE, rv.Media)
	writer.Header().Add("Cache-Control", "max-age=2592000")
	//	writer.Header().Add("Content-Disposition", "attachment;fileName="+rv.Name)

	defer rv.Reader.Close()
	length, err := io.Copy(writer, rv.Reader)

	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
	}

	writer.Header().Add("Content-Length", strconv.Itoa(int(length)))

}

func writeUseJson(writer http.ResponseWriter, obj interface{}) {
	result, err := json.Marshal(obj)
	if err == nil {
		writer.Write(result)
	}
}

func writeUseTemplate(writer http.ResponseWriter, templateName, content string, obj interface{}) {
	tmpl, err := template.New(templateName).Parse(content)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(writer, obj)
	if err != nil {
		panic(err)
	}
}

//解析输入
func parseRequest(request *http.Request, target interface{}) {
	//静态资源的处理
	if sr, ok := target.(*StaticResource); ok {
		sr.Type = request.Header.Get(_CONTENT_TYPE)
		sr.Uri = request.RequestURI
		return
	}

	contentType := request.Header.Get(_CONTENT_TYPE)
	if _CONTENT_TYPE_JSON == contentType || _CONTENT_JSON == contentType { //处理为json的输入
		input, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err == nil {
			parameters := make(map[string]interface{})
			json.Unmarshal(input, &parameters)
			fillStruct(parameters, target)
		}

	} else { //标准form的处理
		if request.Form == nil {
			request.ParseForm()
			fillStructByForm(request.Form, target)
		}
	}

}

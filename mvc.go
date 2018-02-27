package bingo

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
)

const (
	_CONTENT_TYPE      = "Content-Type"
	_CONTENT_TYPE_JSON = "application/json"
	_CONTENT_JSON      = "application/json;charset=utf-8"
	_CONTENT_HTML      = "text/html"
	_CONTENT_XML       = "application/xml;charset=utf-8"
)

//简单的返回结果。用于rest api方式的返回
type SimpleResult struct {
	Action    string
	Success   bool
	ErrorCode int
	Msg       string
}

type ModelView struct {
	View  string
	Model interface{}
}
type StaticResource struct {
	Type string
	Uri  string
}

type RedirectEntity struct {
	Url     string
	Code    int
	Cookies []*http.Cookie
}
type MutiStruct interface {
	GetData() interface{}
	GetDataType() string
}

type FileHandler interface {
	io.Reader
	io.Closer
}

type StaticView struct {
	Name   string      //资源名称
	Media  string      //资源类型
	Length int         //资源长度
	Reader FileHandler //资源内容
}

type routerRule struct {
	url           string
	convert       *ResponseConverter
	methodHandler HttpMethodHandler
}

type Context interface {
	GetSqlSession() *TxSession
	GetCookie(key string) string
}

//返回结果转换器，用于输出返回结果
type ResponseConverter interface {
	Convert(writer http.ResponseWriter, obj interface{})
}

type HttpMethodHandler interface {
	GetSelf() interface{}
	GetParameType(method string) interface{}
	Get(c Context, p interface{}) (interface{}, BingoError)
	Post(c Context, p interface{}) (interface{}, BingoError)
	Put(c Context, p interface{}) (interface{}, BingoError)
	Delete(c Context, p interface{}) (interface{}, BingoError)
}

type HandlerInterceptor interface {
	PreHandle(writer http.ResponseWriter, request *http.Request, handler *routerRule) bool
	PostHandle(writer http.ResponseWriter, request *http.Request, handler *routerRule, mv *ModelView) BingoError
	AfterCompletion(writer http.ResponseWriter, request *http.Request, handler *routerRule, err BingoError) BingoError
}

type Controller struct {
}

func (this *Controller) GetSelf() interface{} {
	return this
}

func (this *Controller) GetParameType(method string) interface{} {
	return this

}
func (this *Controller) Get(c Context, p interface{}) (interface{}, BingoError) {
	return nil, CreateError(Code_NOT_ALLOWED, "method not allowed!")

}
func (this *Controller) Post(c Context, p interface{}) (interface{}, BingoError) {
	return nil, CreateError(Code_NOT_ALLOWED, "method not allowed!")
}
func (this *Controller) Put(c Context, p interface{}) (interface{}, BingoError) {
	return nil, CreateError(Code_NOT_ALLOWED, "method not allowed!")
}
func (this *Controller) Delete(c Context, p interface{}) (interface{}, BingoError) {
	return nil, CreateError(Code_NOT_ALLOWED, "method not allowed!")
}

type SimpleController struct {
	Controller
}

func (this *SimpleController) Post(c Context, p interface{}) (interface{}, BingoError) {
	return this.Get(c, p)
}
func (this *SimpleController) Put(c Context, p interface{}) (interface{}, BingoError) {
	return this.Get(c, p)
}
func (this *SimpleController) Delete(c Context, p interface{}) (interface{}, BingoError) {
	return this.Get(c, p)
}

type staticController struct {
	Controller
	staticDir string
	log Log
}

func (this *staticController) GetParameType(method string) interface{} {
	return &StaticResource{}

}
func (this *staticController) Get(c Context, p interface{}) (interface{}, BingoError) {
	if resource, ok := p.(*StaticResource); ok {
		this.log.Debug("static resource %s,%s",resource.Type,resource.Uri)
		var view StaticView
		var fileDir string
		fileDir, view.Name, view.Media = parseUri(resource.Uri)

		var filePath string = this.staticDir
		if filePath != "" {
			filePath = filePath + "/"
		}
		if fileDir != "" {
			filePath = filePath + fileDir + "/"
		}
		fileRealPath := filePath + view.Name
		fmt.Print(fileRealPath)

		if isFileExist(fileRealPath) {
			fi, err := os.Open(fileRealPath)
			if err != nil {
				this.log.Debug(err.Error())
			} else {
				view.Reader = fi
				return view, nil
			}

		}

	}
	return nil, CreateError(Code_NOT_FOUND, "bingo! The uri not found!")

}

func parseUri(uri string) (dir string, name string, media string) {
	fixIndex := strings.LastIndex(uri, ".")
	lastUrlIndex := strings.LastIndex(uri, "/")
	dir = ""
	if lastUrlIndex > 0 {
		dir = string([]byte(uri)[1:lastUrlIndex])
		dir = strings.Replace(dir, "../", "_", -1)
	}

	if lastUrlIndex < 0 {
		lastUrlIndex = 0
	}

	if fixIndex < 0 {
		fixIndex = len(uri)
	}
	var fileSufix string
	querySufixIndex := strings.LastIndex(uri, "?")
	if querySufixIndex > 0 && fixIndex < querySufixIndex {
		fileSufix = string([]byte(uri)[fixIndex:querySufixIndex])
		name = string([]byte(uri)[lastUrlIndex+1 : querySufixIndex])
	} else {
		fileSufix = string([]byte(uri)[fixIndex:])
		name = string([]byte(uri)[lastUrlIndex+1:])
	}
	fmt.Println(fileSufix)
	return dir, name, getMedia(fileSufix)

}

func getMedia(fileFix string) string {
	media := mime.TypeByExtension(fileFix)
	if media == "" {

	}
	return media
}

package mvc

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/sql"
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

type RouterRule struct {
	url           string
	convert       *ResponseConverter
	methodHandler HttpMethodHandler
}

func (this *RouterRule)Init(url string,handle HttpMethodHandler){
	this.url=url
	this.methodHandler=handle

}

type BeanFactory interface {
	GetLog(module string)utils.Log
	GetService(name string) interface{}
	GetSession() *sql.TxSession
}

type Context interface {
	GetSqlSession() *sql.TxSession
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
	PreHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule) bool
	PostHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule, mv *ModelView) BingoError
	AfterCompletion(writer http.ResponseWriter, request *http.Request, handler *RouterRule, err BingoError) BingoError
}

type HttpController interface{
	Init()
	GetUrl() string
	SetBeanFactory(f BeanFactory)
}

type Controller struct {
	factory BeanFactory
}

func (this *Controller) Init(){

}
func (this *Controller)GetUrl() string{
	return ""
}
func (this *Controller)SetBeanFactory(f BeanFactory){
	this.factory=f
}
func (this *Controller) GetBeanFactory() BeanFactory {
	return this.factory
}
func (this *Controller) GetSelf() interface{} {
	return this
}

func (this *Controller) GetParameType(method string) interface{} {
	return this

}
func (this *Controller) Get(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(utils.Code_NOT_ALLOWED, "method not allowed!")

}
func (this *Controller) Post(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(utils.Code_NOT_ALLOWED, "method not allowed!")
}
func (this *Controller) Put(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(utils.Code_NOT_ALLOWED, "method not allowed!")
}
func (this *Controller) Delete(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(utils.Code_NOT_ALLOWED, "method not allowed!")
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
	log utils.Log
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

		if utils.IsFileExist(fileRealPath) {
			fi, err := os.Open(fileRealPath)
			if err != nil {
				this.log.Debug(err.Error())
			} else {
				view.Reader = fi
				return view, nil
			}

		}

	}
	return nil, utils.CreateError(utils.Code_NOT_FOUND, "bingo! The uri not found!")

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

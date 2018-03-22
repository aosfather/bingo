package mvc

import (
	"net/http"
	"strings"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/sql"
)

type CustomHandlerInterceptor interface {
	PreHandle(writer http.ResponseWriter, request *http.Request) bool
	PostHandle(writer http.ResponseWriter, request *http.Request, mv *ModelView) BingoError
	AfterCompletion(writer http.ResponseWriter, request *http.Request, err BingoError) BingoError
}

type defaultHandlerInterceptor struct {
	interceptors []CustomHandlerInterceptor
}

func (this *defaultHandlerInterceptor) addInterceptor(interceptor CustomHandlerInterceptor) {
	if this.interceptors == nil {
		this.interceptors = []CustomHandlerInterceptor{interceptor}
	} else {
		this.interceptors = append(this.interceptors, interceptor)
	}
}

func (this *defaultHandlerInterceptor) PreHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule) bool {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			if !h.PreHandle(writer, request) {
				return false
			}
		}
	}
	return true
}
func (this *defaultHandlerInterceptor) PostHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule, mv *ModelView) BingoError {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			err := h.PostHandle(writer, request, mv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (this *defaultHandlerInterceptor) AfterCompletion(writer http.ResponseWriter, request *http.Request, handler *RouterRule, err BingoError) BingoError {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			e := h.AfterCompletion(writer, request, err)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

type ContextImp struct {
	txsession *sql.TxSession
	request   *http.Request
}

func (this *ContextImp) GetSqlSession() *sql.TxSession {
	return this.txsession
}

func (this *ContextImp) GetCookie(key string) string {
	if this.request != nil {
		cookie, _ := this.request.Cookie(key)
		if cookie != nil {
			return cookie.Value
		}

	}
	return ""
}

func (this *ContextImp) begin() {
	if this.txsession != nil {
		this.txsession.Begin()
	}
}

func (this *ContextImp) rollback() {
	if this.txsession != nil {
		this.txsession.Rollback()
	}
}
func (this *ContextImp) commit() {
	if this.txsession != nil {
		this.txsession.Commit()
	}
}

type DefaultRouter struct {
	defaultConvert    defaultResponseConverter
	routerMap         map[string]*RouterRule
	interceptor       defaultHandlerInterceptor
	staticHandler     HttpMethodHandler
	//validates         ValidateManager
	sqlsessionfactory *sql.SessionFactory
	logger utils.Log
}

func (this *DefaultRouter) Init(sqlfactory *sql.SessionFactory) {
	this.sqlsessionfactory = sqlfactory
	//this.validates.Init(&DefaultValidaterFactory{})
}
//func (this *DefaultRouter)Validate(obj interface{}) []BingoError{
//	if this.validates.factory==nil {
//		this.validates.Init(&DefaultValidaterFactory{})
//	}
//
//	return this.validates.Validate(obj)
//}

func (this *DefaultRouter)AddInterceptor(h CustomHandlerInterceptor){
	if h!=nil {
		this.interceptor.addInterceptor(h)
	}

}

func (this *DefaultRouter)SetStaticControl(path string,l utils.Log){
	this.staticHandler=&staticController{staticDir: path,log:l}

}
func(this *DefaultRouter)SetLog(l utils.Log){
	this.logger=l
}

func (this *DefaultRouter) SetTemplateRoot(dir string) {
	if utils.IsFileExist(dir) {
		this.defaultConvert.setTemplateDir(dir)
	}

}
func (this *DefaultRouter) AddRouter(rule *RouterRule) {
	if this.routerMap == nil {
		this.routerMap = make(map[string]*RouterRule)
	}
	if rule != nil {
		this.routerMap[rule.url] = rule
	}

}

func (this *DefaultRouter) doConvert(writer http.ResponseWriter, rule *RouterRule, req *http.Request, obj interface{}) {
	if err, ok := obj.(BingoError); ok {
		writer.WriteHeader(err.Code())
		obj = ModelView{"error", err}
	}

	if rule != nil && rule.convert != nil {
		(*rule.convert).Convert(writer, obj)
	} else {
		this.defaultConvert.Convert(writer, obj, req)
	}
}

func (this *DefaultRouter) match(uri string) *RouterRule {
	paramIndex := strings.Index(uri, "?")
	if paramIndex == -1 {
		return this.routerMap[uri]
	} else {
		realuri := strings.TrimSpace((uri[:paramIndex]))
		return this.routerMap[realuri]
	}
}

func (this *DefaultRouter) doMethod(request *http.Request, handler HttpMethodHandler) (interface{}, BingoError) {
	method := request.Method
	param := handler.GetParameType(method)
	parseRequest(this.logger,request, param)
	errors := Validate(param)
	if errors != nil && len(errors) > 0 {
		var errorText string
		for _, err := range errors {
			errorText += err.Error() + ";"
		}
		return nil, utils.CreateError(400, errorText)
	}

	var context ContextImp
	context.request = request
	if this.sqlsessionfactory != nil {
		context.txsession = this.sqlsessionfactory.GetSession()
	}

	//GET 方式不启动事务控制
	if method == utils.Method_GET {
		return handler.Get(&context, param)
	}

	//其余方式，启动事务控制
	context.begin()
	var result interface{}
	var err BingoError
	switch method {
	case utils.Method_POST:
		result, err = handler.Post(&context, param)
	case utils.Method_PUT:
		result, err = handler.Put(&context, param)
	case utils.Method_DELETE:
		result, err = handler.Delete(&context, param)
	default:
		result, err = nil, utils.CreateError(405, "method not found!")
	}

	if err == nil {
		context.commit()
	}
	defer context.rollback()

	return result, err

}

func (this *DefaultRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	uri := request.RequestURI
	rule := this.match(uri)
	//handler前拦截器处理
	if !this.interceptor.PreHandle(writer, request, rule) {
		return
	}

	var handler HttpMethodHandler
	if rule != nil {
		handler = rule.methodHandler
	} else {
		handler = this.staticHandler
	}

	//执行handler处理
	obj, err := this.doMethod(request, handler)

	//handler处理完后，拦截器进行额外补充处理
	if mv, ok := obj.(ModelView); ok {
		this.interceptor.PostHandle(writer, request, rule, &mv)
	} else {
		this.interceptor.PostHandle(writer, request, rule, nil)
	}

	//进行结果输出
	if err != nil {
		this.doConvert(writer, rule, request, err)
	} else {
		this.doConvert(writer, rule, request, obj)
	}

	//请求处理完后拦截器进行处理
	this.interceptor.AfterCompletion(writer, request, rule, err)

}

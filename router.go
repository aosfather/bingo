package bingo

import (
	"net/http"
	"strings"
)

type defaultHandlerInterceptor struct {
	interceptors []HandlerInterceptor
}

func (this *defaultHandlerInterceptor) addInterceptor(interceptor HandlerInterceptor) {
	if this.interceptors == nil {
		this.interceptors = []HandlerInterceptor{interceptor}
	} else {
		this.interceptors = append(this.interceptors, interceptor)
	}
}

func (this *defaultHandlerInterceptor) PreHandle(writer http.ResponseWriter, request *http.Request, handler *routerRule) bool {
	return true
}
func (this *defaultHandlerInterceptor) PostHandle(writer http.ResponseWriter, request *http.Request, handler *routerRule, mv *ModelView) BingoError {
	return nil
}
func (this *defaultHandlerInterceptor) AfterCompletion(writer http.ResponseWriter, request *http.Request, handler *routerRule, err BingoError) BingoError {
	return nil
}

type ContextImp struct {
	txsession *TxSession
}

func (this *ContextImp) GetSqlSession() *TxSession {
	return this.txsession
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

type defaultRouter struct {
	defaultConvert    defaultResponseConverter
	routerMap         map[string]*routerRule
	interceptor       defaultHandlerInterceptor
	staticHandler     HttpMethodHandler
	validates         validateManager
	sqlsessionfactory *SessionFactory
}

func (this *defaultRouter) Init(sqlfactory *SessionFactory) {
	this.sqlsessionfactory = sqlfactory
	this.validates.Init(&defaultValidaterFactory{})
}
func (this *defaultRouter) setTemplateRoot(dir string) {
	if isFileExist(dir) {
		this.defaultConvert.setTemplateDir(dir)
	}

}
func (this *defaultRouter) addRouter(rule *routerRule) {
	if this.routerMap == nil {
		this.routerMap = make(map[string]*routerRule)
	}
	if rule != nil {
		this.routerMap[rule.url] = rule
	}

}

func (this *defaultRouter) doConvert(writer http.ResponseWriter, rule *routerRule, req *http.Request, obj interface{}) {
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

func (this *defaultRouter) match(uri string) *routerRule {
	paramIndex := strings.Index(uri, "?")
	if paramIndex == -1 {
		return this.routerMap[uri]
	} else {
		realuri := strings.TrimSpace((uri[:paramIndex]))
		return this.routerMap[realuri]
	}
}

func (this *defaultRouter) doMethod(request *http.Request, handler HttpMethodHandler) (interface{}, BingoError) {
	method := request.Method
	param := handler.GetParameType(method)
	parseRequest(request, param)
	errors := this.validates.Validate(param)
	if errors != nil && len(errors) > 0 {
		var errorText string
		for _, err := range errors {
			errorText += err.Error() + ";"
		}
		return nil, CreateError(400, errorText)
	}

	var context ContextImp
	if this.sqlsessionfactory != nil {
		context.txsession = this.sqlsessionfactory.GetSession()
	}

	//GET 方式不启动事务控制
	if method == Method_GET {
		return handler.Get(&context, param)
	}

	//其余方式，启动事务控制
	context.begin()
	var result interface{}
	var err BingoError
	switch method {
	case Method_POST:
		result, err = handler.Post(&context, param)
	case Method_PUT:
		result, err = handler.Put(&context, param)
	case Method_DELETE:
		result, err = handler.Delete(&context, param)
	default:
		result, err = nil, CreateError(405, "method not found!")
	}

	if err == nil {
		context.commit()
	}
	defer context.rollback()

	return result, err

}

func (this *defaultRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

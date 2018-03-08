package bingo

import (
	"net/http"
	"strconv"
)

const (
	_DEFAULT_ADMIN_PORT=18990
)
type adminService struct {
	application *TApplication
}

func (this *adminService) run(){
	port, err := strconv.Atoi(this.application.context.getProperty("bingo.admin.port"))
	if err != nil {
		port=_DEFAULT_ADMIN_PORT
	}
	go http.ListenAndServe(":"+strconv.Itoa(port), this)
}

func (this *adminService)ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("admin"))
}

//执行shutdown操作
func (this *adminService)doShutDown(writer http.ResponseWriter) {

	//执行自定义
	if this.application.onShutdown!=nil {
		this.application.onShutdown()
	}
}

type OnDestoryHandler func() bool //shutdown的handler，用于处理关闭服务的自定义动作
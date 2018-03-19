package utils

import (
	"log"
	"runtime"
	"fmt"
	"time"
)

const(
	LEVEL_DEBUG=1
	LEVEL_INFO=2
	LEVEL_WARN=3
	LEVEL_ERROR=4
)
var LEVEL_NAMES=[]string{"debug","info","warn","error"}
type LogConfig struct {
	IsDebug bool
	FileName string
}

type logRecord struct {
	format string
	objs []interface{}
}

type LogFactory struct {
	loglevel int
	logfile *RollingFile
	l *log.Logger
	queue *Queue
	running bool
}

func (this *LogFactory)SetConfig(config LogConfig) {
	if config.FileName=="" {
		log.Println("not set log file error! set default file out.log")
		config.FileName="bingo_out.log"
	}
	r:=RollingFile{}
	r.Filename=config.FileName
	this.logfile=&r
	this.l = log.New(this.logfile, "", log.LstdFlags)

	this.queue=NewQueue()
	this.running=true
	go this.outThread()

}

func (this *LogFactory)outThread(){
	for{
		if v,ok:=this.queue.pop().(*logRecord);ok {
			this.l.Printf(v.format, v.objs...)
		}else {
			time.Sleep(time.Millisecond)
			if this.running==false {//当设置为不要running的时候，退出循环
				return
			}
		}
	}
}

func (this *LogFactory)GetLog(name string)Log {
	return &logImp{name,this}
}

func (this *LogFactory)Close(){
	this.running=false
	if this .logfile!=nil {
		this.logfile.Close()
	}
}

func (this *LogFactory)Write(level int,prefix string,fmt string,obj ...interface{}) {
	   if level >= this.loglevel { //判断loglevel是否大于指定的level，如果大于则输出，否则直接抛弃
		   content := this.formatHeader(prefix, level) + fmt + "\n"
		   this.queue.push(&logRecord{content,obj})

	   }
}

func (this *LogFactory) formatHeader(prefix string,level int) string {
	_, thefile, line, _ := runtime.Caller(3)
	short := thefile
	for i := len(thefile) - 1; i > 0; i-- {
		if thefile[i] == '/' {
			short = thefile[i+1:]
			break
		}
	}
	return fmt.Sprintf("[%s][%s][%s(%d)]- ",LEVEL_NAMES[level-1],prefix,short,line)
}

type Log interface {
	Info(msg string,obj ...interface{})
	Debug(msg string,obj ...interface{})
	Error(msg string,obj ...interface{})
	Warning(msg string,obj ...interface{})
}

type logImp struct {
	module string
	factory *LogFactory
}
func (this *logImp)Info(msg string,obj ...interface{}){
   this.factory.Write(LEVEL_INFO,this.module,msg,obj...)
}

func (this *logImp)Debug(msg string,obj ...interface{}){
	this.factory.Write(LEVEL_DEBUG,this.module,msg,obj...)
}
func (this *logImp)Error(msg string,obj ...interface{}){
	this.factory.Write(LEVEL_ERROR,this.module,msg,obj...)
}
func (this *logImp)Warning(msg string,obj ...interface{}){
	this.factory.Write(LEVEL_WARN,this.module,msg,obj...)
}
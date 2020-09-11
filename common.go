package main

import (
	"fmt"
	"github.com/aosfather/bingo_utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"time"
)

func init() {
	log.SetFlags(log.Lmsgprefix)
	bingo_utils.SetLogDebugFunc(debug)
	bingo_utils.SetLogErrFunc(errs)
}

func debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("DEBUG", msg)
}

func info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("INFO", msg)
}

func errs(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("ERROR", msg)
}

func errsf(formate string, v ...interface{}) {
	msg := fmt.Sprintf(formate, v...)
	_log("ERROR", msg)
}

func _log(level string, msg string) {
	now := time.Now().Format(bingo_utils.FORMAT_DATETIME_LOG)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	} else {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	log.Print(fmt.Sprintf("%s [%s] [%s:%d] [%s]\n", now, level, file, line, msg))
}

//通用返回结果
type Result struct {
	Code  int           `json:"code"`
	Msg   string        `json:"msg"`
	Count int           `json:"count"`
	Data  []interface{} `json:"data"`
}

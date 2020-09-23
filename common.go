package main

import (
	"fmt"
	"github.com/aosfather/bingo_utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"time"
)

const _noprefix = 64

func init() {
	log.SetFlags(_noprefix)
	bingo_utils.SetLogDebugFunc(_debug)
	bingo_utils.SetLogErrFunc(_errs)
}

func _errs(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("ERROR", msg, 3)
}

func _debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("DEBUG", msg, 3)
}

func debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("DEBUG", msg, 2)
}

func info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("INFO", msg, 2)
}

func errs(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("ERROR", msg, 2)
}

func errsf(formate string, v ...interface{}) {
	msg := fmt.Sprintf(formate, v...)
	_log("ERROR", msg, 2)
}

func _log(level string, msg string, skip int) {
	now := time.Now().Format(bingo_utils.FORMAT_DATETIME_LOG)
	_, file, line, ok := runtime.Caller(skip)
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

package main

import (
	mvc "github.com/aosfather/bingo_mvc"
	// "github.com/aosfather/bingo_dao"
)

type ApiType byte

const (
	AT_ModelView = 11 //
	AT_Servlet   = 12
	AT_Task      = 13
)

//接口定义
type Api struct {
	Name        string
	Url         string
	Description string
	Type        ApiType
	Style       mvc.StyleType
	Methods     []mvc.HttpMethodType
	Request     []Paramter
}

// table
type Paramter struct {
}

//应用
type App struct {
}

type Metas struct {
	Types    string
	Elements string
	Dict     string
	Struct   string
	Table    string
}

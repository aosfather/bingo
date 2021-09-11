package main

import (
	"testing"
)

//元数据表定义
func TestLoadFromYaml(t *testing.T) {
	m := Meta{}
	LoadFromYaml("app/meta.yaml", &m)
	t.Log(m)
	e := Elements{}
	LoadElementsFromYaml("app/elements.yaml", &e)
	t.Log(e.Elements[0])
	t.Log(e.Elements[0].Type())

	tables := Tables{}
	tables.LoadTablesFromYaml("app/tables.yaml")
	t.Log(tables)
}

package main

import (
	"testing"
)

func TestLoadFromYaml(t *testing.T) {
	m := Meta{}
	LoadFromYaml("app/meta.yaml", &m)
	t.Log(m)
	e := Elements{}
	LoadElementsFromYaml("app/elements.yaml", &e)
	t.Log(e.Elements[0])
	t.Log(e.Elements[0].Type())

	tables := Tables{}
	LoadTablesFromYaml("app/tables.yaml", &tables)
	t.Log(tables)
}

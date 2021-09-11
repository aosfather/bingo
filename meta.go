package main

import (
	"github.com/aosfather/bingo_dao"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

//基本元素定义。类型和字典
type Meta struct {
	Types []bingo_dao.DataType `yaml:"types"`
	Dicts []bingo_dao.DictCatalog
}

func LoadFromYaml(name string, m *Meta) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Println(err.Error())
	}
	err = yaml.Unmarshal(data, m)
	if err != nil {
		log.Println(err.Error())
	}
	for _, t := range m.Types {
		bingo_dao.GetTypes().AddType(&t)
	}
}

type Elements struct {
	Elements []bingo_dao.DataElement
}

func LoadElementsFromYaml(name string, m *Elements) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Println(err.Error())
	}
	err = yaml.Unmarshal(data, m)
	if err != nil {
		log.Println(err.Error())
	}

	for _, t := range m.Elements {
		bingo_dao.GetTypes().AddElement(&t)
	}
}

/**
  表及结构的定义
  包括表列表和结构列表
*/
type Tables struct {
	Tables  []bingo_dao.Table
	Structs []bingo_dao.DataStruct
}

func (this *Tables) LoadTablesFromYaml(name string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Println(err.Error())
	}
	err = yaml.Unmarshal(data, this)
	if err != nil {
		log.Println(err.Error())
	}

	//for _,t:=range m.Elements {
	//	bingo_dao.GetTypes().AddElement(&t)
	//}
}

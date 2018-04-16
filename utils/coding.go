package utils

import "strings"

//JSON的万能转换，将其它类型转换成string
type AsString struct {
	Value string
}

func (this *AsString) UnmarshalJSON(data []byte) error {
	if len(data)==0 {
		return nil
	}

	if data[0]=='"'&&len(data)>1 {
		end:=len(data)
		data=data[1:end-1]
	}

	str:=string(data)
	if str!="null" {
		this.Value=str
	}

	return nil
}

//作为普通字符串输出
func (this *AsString) MarshalJSON()([]byte,error) {
	return []byte("\""+this.Value+"\""),nil
}



const (
	SPLIT_CHAR = ","
    QUATO = "\""
	)

//标签
type Tags struct{
	Values []string
}

func (this *Tags) FromString(str string) {
	if str!=""{
		this.Values=strings.Split(str, SPLIT_CHAR)
	}
}

func (this *Tags) UnmarshalJSON(data []byte) error {
	if len(data)==0 {
		return nil
	}

	if data[0]=='"'&&len(data)>1 {
		end:=len(data)
		data=data[1:end-1]
	}

	str:=string(data)
	if str!="null" {

		this.Values=strings.Split(str, SPLIT_CHAR)
	}

	return nil
}

//作为普通字符串输出
func (this *Tags) MarshalJSON()([]byte,error) {

	return []byte(QUATO +this.ToString()+ QUATO),nil
}

func (this *Tags)ToString() string {
	return strings.Join(this.Values,SPLIT_CHAR)
}
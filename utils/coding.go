package utils

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

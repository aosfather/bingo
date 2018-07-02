package utils

import (
	"testing"
)

func TestUTF82GB2312(t *testing.T) {
	InitDict("../data/pinyin.dict")
	if ToFirstPyLetter("中国田径") != "zgtj" {
		t.Fail()
	}

	t.Log(ToPinyin("处理大字符串", " "))
	t.Log(ToPinyin("处理大字符串", ""))

}

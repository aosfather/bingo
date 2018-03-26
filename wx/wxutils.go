package wx

import (
	"time"
	"sort"
	"crypto/sha1"
	"io"
	"fmt"
	"net/http"
	"io/ioutil"
	"math/rand"
)

var (
	CHARS []byte = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	charlen := len(CHARS)
	for i := 0; i < length; i++ {
		result = append(result, CHARS[r.Intn(charlen)])
	}
	return string(result)
}

//GetCurrTs return current timestamps
func GetCurrTs() int64 {
	return time.Now().Unix()
}

//Signature sha1签名
func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

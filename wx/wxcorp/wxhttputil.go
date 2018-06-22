package wxcorp

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
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

func PostToWx(theUrl string, data interface{}, result interface{}) error {
	content, _ := json.Marshal(data)
	resp, err := http.Post(theUrl, "application/json;charset=utf-8", strings.NewReader(string(content)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	json.Unmarshal(body, result)
	return nil
}

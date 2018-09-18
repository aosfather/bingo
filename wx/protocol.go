package wx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/aosfather/bingo/wx/wxcorp"
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

//---------------协议算法------------------//
// 生成4个字节的网络字节序
func NumberToBytesOrder(number int) []byte {
	var orderBytes []byte
	orderBytes = make([]byte, 4, 4)
	orderBytes[3] = byte(number & 0xFF)
	orderBytes[2] = byte(number >> 8 & 0xFF)
	orderBytes[1] = byte(number >> 16 & 0xFF)
	orderBytes[0] = byte(number >> 24 & 0xFF)

	return orderBytes
}

// 还原4个字节的网络字节序
func BytesOrderToNumber(orderBytes []byte) int {
	var number int = 0

	for i := 0; i < 4; i++ {
		number <<= 8
		number |= int(orderBytes[i] & 0xff)
	}
	return number
}

//---------------------------PKCS7-----------------------------//
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//-------------------AES----------------------//

type AES struct {
	key       []byte
	block     cipher.Block
	blockSize int
}

func (this *AES) Init(key []byte) {
	this.key = key
	block, err := aes.NewCipher(this.key)
	if err != nil {
		return
	}
	this.block = block
	this.blockSize = this.block.BlockSize()
}

func (this *AES) encrypt(sourceText []byte) []byte {
	sourceText = PKCS7Padding(sourceText, this.blockSize)

	blockModel := cipher.NewCBCEncrypter(this.block, this.key[:this.blockSize])

	ciphertext := make([]byte, len(sourceText))

	blockModel.CryptBlocks(ciphertext, sourceText)
	return ciphertext
}

func (this *AES) decrypt(encryptedText []byte) []byte {
	blockMode := cipher.NewCBCDecrypter(this.block, this.key[:this.blockSize])
	origData := make([]byte, len(encryptedText))
	blockMode.CryptBlocks(origData, encryptedText)
	origData = PKCS7UnPadding(origData, this.blockSize)
	return origData
}

//Encrypt
type BaseEncrypt struct {
	token          string
	appId          string
	encodingAESKey string
	theAES         AES
}

func (this *BaseEncrypt) Init(token, appid, aeskey string) {
	this.token = token
	this.appId = appid
	if len(aeskey) != 43 { //密码长度不对失败
		panic(aeskey)
	}
	this.encodingAESKey = aeskey
	asekey, _ := base64.URLEncoding.DecodeString(this.encodingAESKey + "=")
	this.theAES = AES{}
	this.theAES.Init(asekey)

}

func (this *BaseEncrypt) VerifyURL(msg_signature, echostr string, paramters ...string) bool {
	p := []string{this.token}
	p = append(p, paramters...)
	signature := Signature(p...)
	fmt.Println("signature:" + signature + "|" + msg_signature)
	if signature != msg_signature {
		return false

	}
	return true

}

//加密消息
func (this *BaseEncrypt) EncryptMsg(replyMsg, nonce, timestamp string) (int, string) {
	msg := WxEncryptOutputMessage{}
	msg.Encrypt = this.encrypt(replyMsg)
	msg.Nonce = nonce
	msg.MsgSignature = Signature(this.token, timestamp, nonce, msg.Encrypt)
	msg.TimeStamp = timestamp
	outxml, _ := xml.Marshal(msg)
	return 0, string(outxml)
}

//加密
func (this *BaseEncrypt) encrypt(text string) string {
	var byteGroup bytes.Buffer
	//格式 16位的random字符+文本长度（4位的网络字节序列）+文本+企业id
	randStr := RandomStr(16)
	byteGroup.Write([]byte(randStr))
	byteGroup.Write(NumberToBytesOrder(len(text)))
	byteGroup.Write([]byte(text))

	byteGroup.Write([]byte(this.appId))

	encryptedText := this.theAES.encrypt(byteGroup.Bytes())
	//转base64
	return base64.StdEncoding.EncodeToString(encryptedText)
}

//机密消息
func (this *BaseEncrypt) DecryptInputMsg(msg_signature, timestamp, nonce string, postdata WxEncryptInputMessage) (int, string) {
	signature := Signature(this.token, timestamp, nonce, postdata.Encrypt)
	if msg_signature != signature {
		return 40001, "signature validate failed!"
	}

	msg := this.decrypt(postdata.Encrypt, postdata.ToUserName)
	if msg == "" {
		return 40005, ""
	}
	return 0, msg
}

func (this *BaseEncrypt) decrypt(encryptstr string, targetcorpid string) string {
	//1、解析base64
	str, _ := base64.StdEncoding.DecodeString(encryptstr)
	sourcebytes := this.theAES.decrypt(str)
	//取字节数组长度
	orderBytes := sourcebytes[16:20]
	msgLength := BytesOrderToNumber(orderBytes)

	//取文本
	text := sourcebytes[20 : 20+msgLength]
	//取应用id
	appid := sourcebytes[20+msgLength:]
	fmt.Println("the id:[" + string(appid) + "]")
	corp := strings.TrimSpace(string(appid))
	fmt.Printf("%s:", this.appId)
	if this.appId == corp {
		fmt.Println("yes!")
		return string(text)
	}

	return ""

}

//输入的加密消息
type WxEncryptInputMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	AppId      string
	Encrypt    string
}

//输出的加密消息
type WxEncryptOutputMessage struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string
	MsgSignature string
	TimeStamp    string
	Nonce        string
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

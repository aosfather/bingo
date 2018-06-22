package wxcorp

/**
企业微信加密信息工具类

auth：xiongxiaopeng@qianbaoplus.com
by 钱包行云
2017.10.27

**/
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	//	"math/rand"
	"sort"
	"strings"
	//	"time"
)

/**
40001	签名验证错误
40002	xml解析失败
40003	sha加密生成签名失败
40004	AESKey 非法
40005	corpid 校验错误
40006	AES 加密失败
40007	AES 解密失败
40008	解密后得到的buffer非法
*/
const (
	Block_size  = 32
	BASE_STRING = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

//输出给企业微信服务端的消息
type CorpOutputMessage struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string
	MsgSignature string
	TimeStamp    string
	Nonce        string
}

//接收到企业微信服务端的消息
type CorpInputputMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	AgentID    string
	Encrypt    string
}

type CorpEncrypt struct {
	token          string
	corpid         string
	suiteid        string
	encodingAESKey string
	theAES         AES
}

func (this *CorpEncrypt) Init(token, corpid, suiteid, aeskey string) {
	this.token = token
	this.corpid = corpid
	this.suiteid = suiteid
	if len(aeskey) != 43 { //密码长度不对失败
		panic(aeskey)
	}
	this.encodingAESKey = aeskey
	asekey, _ := base64.URLEncoding.DecodeString(this.encodingAESKey + "=")
	this.theAES = AES{}
	this.theAES.Init(asekey)
}

/**
	 #验证URL
         #@param sMsgSignature: 签名串，对应URL参数的msg_signature
         #@param sTimeStamp: 时间戳，对应URL参数的timestamp
         #@param sNonce: 随机串，对应URL参数的nonce
         #@param sEchoStr: 随机串，对应URL参数的echostr
         #@param sReplyEchoStr: 解密之后的echostr，当return返回0时有效
         #@return：成功0，失败返回对应的错误码

*/
func (this *CorpEncrypt) VerifyURL(msg_signature, timestamp, nonce, echostr string) (int, string) {

	signature := makeSignature(this.token, timestamp, nonce, echostr)
	fmt.Println("signature:" + signature + "|" + msg_signature)
	if signature != msg_signature {
		return 40001, "signature validate failed"

	}
	return 0, this.decrypt(echostr)

}

func (this *CorpEncrypt) DecryptMsg(msg_signature, timestamp, nonce, postdata string) (int, string) {
	input := CorpInputputMessage{}
	err := xml.Unmarshal([]byte(postdata), &input)
	if err != nil {
		return 40002, err.Error()
	}

	return this.DecryptInputMsg(msg_signature, timestamp, nonce, input)
}

func (this *CorpEncrypt) DecryptInputMsg(msg_signature, timestamp, nonce string, postdata CorpInputputMessage) (int, string) {

	signature := makeSignature(this.token, timestamp, nonce, postdata.Encrypt)
	if msg_signature != signature {
		return 40001, "signature validate failed!"
	}

	msg := this.decrypt(postdata.Encrypt)
	if msg == "" {
		return 40005, ""
	}
	return 0, msg
}

func (this *CorpEncrypt) EncryptMsg(replyMsg, nonce, timestamp string) (int, string) {
	msg := CorpOutputMessage{}
	msg.Encrypt = this.encrypt(replyMsg)
	msg.Nonce = nonce
	msg.MsgSignature = makeSignature(this.token, timestamp, nonce, msg.Encrypt)
	msg.TimeStamp = timestamp
	outxml, _ := xml.Marshal(msg)
	return 0, string(outxml)
}

//加密
func (this *CorpEncrypt) encrypt(text string) string {
	var byteGroup bytes.Buffer
	//格式 16位的random字符+文本长度（4位的网络字节序列）+文本+企业id
	randStr := this.getRandomString()
	byteGroup.Write([]byte(randStr))
	byteGroup.Write(numberToBytesOrder(len(text)))
	byteGroup.Write([]byte(text))

	byteGroup.Write([]byte(this.corpid))

	encryptedText := this.theAES.encrypt(byteGroup.Bytes())
	//转base64
	return base64.StdEncoding.EncodeToString(encryptedText)
}

//解密
func (this *CorpEncrypt) decrypt(encryptstr string) string {
	//1、解析base64
	str, _ := base64.StdEncoding.DecodeString(encryptstr)
	sourcebytes := this.theAES.decrypt(str)
	//取字节数组长度
	orderBytes := sourcebytes[16:20]
	msgLength := bytesOrderToNumber(orderBytes)

	//取文本
	text := sourcebytes[20 : 20+msgLength]
	//取企业id
	corpid := sourcebytes[20+msgLength:]
	fmt.Println("the corpid:[" + string(corpid) + "]")
	corp := strings.TrimSpace(string(corpid))
	fmt.Printf("%s,%s:", this.corpid, this.suiteid)
	if this.corpid == corp || this.suiteid == corp {
		fmt.Println("yes!")
		return string(text)

	}

	return ""

}

func (this *CorpEncrypt) getRandomString() string {
	return RandomStr(16)
	//	seed := rand.NewSource(time.Now().UnixNano())
	//	r := rand.New(seed)
	//	lenth := len(BASE_STRING)
	//	buffer := bytes.NewBufferString("")
	//	for i := 0; i < 16; i++ {
	//		buffer.WriteString(string(BASE_STRING[r.Intn(lenth)]))
	//	}

	//	return buffer.String()
}

//-------------网络字节序列----------------//

// 生成4个字节的网络字节序
func numberToBytesOrder(number int) []byte {
	var orderBytes []byte
	orderBytes = make([]byte, 4, 4)
	orderBytes[3] = byte(number & 0xFF)
	orderBytes[2] = byte(number >> 8 & 0xFF)
	orderBytes[1] = byte(number >> 16 & 0xFF)
	orderBytes[0] = byte(number >> 24 & 0xFF)

	return orderBytes
}

// 还原4个字节的网络字节序
func bytesOrderToNumber(orderBytes []byte) int {
	var number int = 0

	for i := 0; i < 4; i++ {
		number <<= 8
		number |= int(orderBytes[i] & 0xff)
	}
	return number
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

//-------------------------------SHA1------------------------------//
func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func makeSignature(token string, timestamp, nonce string, msg string) string {
	sl := []string{token, timestamp, nonce, msg}
	sort.Strings(sl)
	return str2sha1(strings.Join(sl, ""))
}

func MakeSignatureForJs(token string, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	return str2sha1(strings.Join(sl, ""))
}

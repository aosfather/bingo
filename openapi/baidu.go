package openapi

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aosfather/comb"
	"github.com/aosfather/bingo/utils"
)

const (
	BAIDU_TOKEN_URL          = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"
	BAIDU_OCR_URL            = "https://aip.baidubce.com/rest/2.0/ocr/v1/%s"
	BAIDU_idcardUrl          = "https://aip.baidubce.com/rest/2.0/ocr/v1/idcard"
	BAIDU_bankcardUrl        = "https://aip.baidubce.com/rest/2.0/ocr/v1/bankcard"
	BAIDU_generalUrl         = "https://aip.baidubce.com/rest/2.0/ocr/v1/general"
	BAIDU_basicGeneralUrl    = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
	BAIDU_webImageUrl        = "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage"
	BAIDU_enhancedGeneralUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_enhanced"
	BAIDU_drivingLicenseUrl  = "https://aip.baidubce.com/rest/2.0/ocr/v1/driving_license"
	BAIDU_vehicleLicenseUrl  = "https://aip.baidubce.com/rest/2.0/ocr/v1/vehicle_license"
	BAIDU_tableRequestUrl    = "https://aip.baidubce.com/rest/2.0/solution/v1/form_ocr/request"
	BAIDU_tableResultUrl     = "https://aip.baidubce.com/rest/2.0/solution/v1/form_ocr/get_request_result"
	BAIDU_licensePlateUrl    = "https://aip.baidubce.com/rest/2.0/ocr/v1/license_plate"
	BAIDU_accurateUrl        = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate"
	BAIDU_basicAccurateUrl   = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
	BAIDU_receiptUrl         = "https://aip.baidubce.com/rest/2.0/ocr/v1/receipt"
	BAIDU_businessLicenseUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/business_license"
	MAX_SIZE                 = 4 * 1024 * 1024
)

var (
	IMAGE_FORMATS = []string{"JPEG", "BMP", "PNG"}
)

type BaiduErrorMessage struct {
	Code string `json:"error"`
	Msg  string `json:"error_description"`
}
type BaiduAccessToken struct {
	BaiduErrorMessage
	Expires       int64  `json:"expires_in"` //过期时间，单位秒。一般为1个月)
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
	AccessToken   string `json:"access_token"` //要获取的Access Token；
	RefreshToken  string `json:"refresh_token"`
}
type Parameter map[string]string

func BuildParameter() Parameter {
	return make(map[string]string)
}

type BaiduSdk struct {
	AppId     string
	AppKey    string
	AppSecret string
	Token     *BaiduAccessToken
}

func (this *BaiduSdk) getAccessToken() {
	theUrl := fmt.Sprintf(BAIDU_TOKEN_URL, this.AppKey, this.AppSecret)
	token := BaiduAccessToken{}
	utils.Post(theUrl, nil, &token)
	//	if token.Code != "" {
	this.Token = &token

	//	}

}

func (this *BaiduSdk) Call(theUrl string, p baiduQuery, result interface{}) {
	if this.Token == nil {
		this.getAccessToken()
	}
	if p != nil {
		parames := p.toParameter()
		parames["access_token"] = this.Token.AccessToken
		err := utils.PostByForm(theUrl, parames, result)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

//从文件中加载图片内容，对于不符合的格式，和大小进行检查
func (this *BaiduOcr) LoadImageToByte(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}

	img, imgType, imgerr := image.DecodeConfig(file)
	if imgerr != nil {
		return nil, imgerr
	}
	imgType = strings.ToUpper(imgType)

	//1、检测文件类型
	rightFormat := false
	for _, v := range IMAGE_FORMATS {
		if imgType == v {
			rightFormat = true
			break
		}
	}

	if !rightFormat {
		return nil, fmt.Errorf("%s%s", "图像格式错误！只支持", IMAGE_FORMATS)
	}

	//2、检测图像大小
	if img.Width < 15 || img.Width > 4096 || img.Height < 15 || img.Height > 4096 {
		return nil, fmt.Errorf("%s", "图像大小不合适！最短边至少15px，最长边最大4096px")
	}
	file.Close()

	file, err = os.Open(path)
	//3、检测转码后的大小
	content, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return nil, err2
	}

	text := bytesTOBaiduBase64(content)
	if len(text) >= MAX_SIZE {
		return nil, fmt.Errorf("%s", "图像文件编码后过大超过4M了")

	}

	defer file.Close()

	return content, nil
}

func bytesTOBaiduBase64(content []byte) string {

	//	return base64.URLEncoding.EncodeToString(content)
	return base64.StdEncoding.EncodeToString(content)
}

type baiduQuery interface {
	toParameter() Parameter
}
type IdCardQuery struct {
	Detect bool //true、false	是否检测图像朝向，默认不检测
	Front  bool //true 表示front身份证正面；back：身份证背面
	Image  []byte
	Risk   bool //是否开启身份证风险检测
}

func (this *IdCardQuery) toParameter() Parameter {
	p := BuildParameter()
	if this.Detect {
		p["detect_direction"] = "true"
	}
	if this.Front {
		p["id_card_side"] = "front"
	} else {
		p["id_card_side"] = "back"
	}

	if this.Risk {
		p["detect_risk"] = "true"
	}

	p["image"] = bytesTOBaiduBase64(this.Image)

	return p
}

type baseResult struct {
	Id int64 `json:"log_id"` //唯一的log id，用于问题定位
	BaiduErrorMessage
}

type BaseImageQuery struct {
	Image []byte
}

func (this *BaseImageQuery) toParameter() Parameter {
	p := BuildParameter()
	p["image"] = bytesTOBaiduBase64(this.Image)

	return p
}

//身份证识别结果
type IdCardResult struct {
	baseResult
	Direction int32  `json:"direction"`    //图像方向，当detect_direction=true时存在。-1:未定义，0:正向，1: 逆时针90度，2:逆时针180度，3:逆时针270度
	Status    string `json:"image_status"` //normal-识别正常reversed_side-未摆正身份证non_idcard-上传的图片中不包含身份证blurred-身份证模糊over_exposure-身份证关键字段反光或过曝unknown-未知状态
	Risk      string `json:"idcard_type"`  //则返回该字段识别身份证类型: normal-正常身份证；copy-复印件；temporary-临时身份证；screen-翻拍；unknow-其他未知情况
	EditTool  string `json:"edit_tool"`    //如果检测身份证被编辑过，该字段指定编辑软件名称，如:Adobe Photoshop CC 2014 (Macintosh),如果没有被编辑过则返回值无此参数

	Result map[string]RecognizeItem `json:"words_result"`     //识别结果
	Count  int32                    `json:"words_result_num"` // 识别结果数，表示words_result的元素个数
}

type RecognizeItem struct {
	Location Image_location `json:"location"`

	Text string `json:"words"`
}
type Image_location struct {
	Left   int64 `json:"left"`
	Top    int64 `json:"top"`
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

//银行卡识别结果
type BankCardResult struct {
	baseResult
	Result BankCard `json:"result"`
}

type BankCard struct {
	CardNo   string //银行卡卡号
	CardType int32  //银行卡类型，0:不能识别; 1: 借记卡; 2: 信用卡
	BankName string //银行名，不能识别时为空
}

//驾驶证识别结果

//行驶证识别结果

//车牌识别
type CarNoQuery struct {
	Image []byte
	Multi bool //是否是多张一起检测
}

func (this *CarNoQuery) toParameter() Parameter {
	p := BuildParameter()
	if this.Multi {
		p["multi_detect"] = "true"
	}
	p["image"] = bytesTOBaiduBase64(this.Image)

	return p
}

type CarNoMutiResult struct {
	baseResult
	Result []CarNo `json:"words_result"`
}

type CarNoResult struct {
	baseResult
	Result CarNo `json:"words_result"`
	Array  []CarNo
}
type CarNo struct {
	Color  string `json:"color"`  //车牌颜色
	Number string `json:"number"` //车牌号码
}

//营业执照
type BusinessLicenseResult struct {
	baseResult
	Result map[string]RecognizeItem `json:"words_result"`     //识别结果
	Count  int32                    `json:"words_result_num"` // 识别结果数，表示words_result的元素个数
}

type BaiduOcr struct {
	sdk *BaiduSdk
}

func (this *BaiduOcr) Init(sdk *BaiduSdk) {
	this.sdk = sdk
}

//识别身份证
func (this *BaiduOcr) RecognizedIdCard(card *IdCardQuery) IdCardResult {
	result := IdCardResult{}
	this.sdk.Call(BAIDU_idcardUrl, card, &result)
	return result
}

//识别车牌号码
func (this *BaiduOcr) RecognizedCarNo(card *CarNoQuery) CarNoResult {
	result := CarNoResult{}
	if card.Multi {
		muti := CarNoMutiResult{}
		this.sdk.Call(BAIDU_licensePlateUrl, card, &muti)
		result.Array = muti.Result
		if len(muti.Result) > 0 {
			result.Result = muti.Result[0]
		}

	} else {
		this.sdk.Call(BAIDU_licensePlateUrl, card, &result)
		result.Array = []CarNo{result.Result}
	}

	return result
}

//识别银行卡
func (this *BaiduOcr) RecognizedBankCard(card *BaseImageQuery) BankCardResult {
	result := BankCardResult{}
	this.sdk.Call(BAIDU_bankcardUrl, card, &result)
	return result
}

//识别营业执照
func (this *BaiduOcr) RecognizedBusiness(card *BaseImageQuery) BusinessLicenseResult {
	result := BusinessLicenseResult{}
	this.sdk.Call(BAIDU_businessLicenseUrl, card, &result)
	return result
}

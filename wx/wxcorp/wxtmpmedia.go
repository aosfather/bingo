package wxcorp

//上传下载临时素材
const (
	UPLOAD_MEDIA_API   = "https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	DOWNLOAD_MEDIA_API = "https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
)

type wxUploadResult struct {
	baseMessage
	Type       string `json:"type"`
	MediaId    string `json:"media_id"`
	CreateTime string `json:"created_at"`
}

type WxMediaManager struct {
}

func (this *WxMediaManager) Upload(filename string, content []byte) error {

	return nil

}

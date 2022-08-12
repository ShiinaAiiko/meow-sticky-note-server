package controllersV1

import (
	"strings"
	"time"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	"github.com/ShiinaAiiko/meow-sticky-note/server/protos"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/response"
	"github.com/cherrai/nyanyago-utils/cipher"
	"github.com/cherrai/nyanyago-utils/nstrings"
	"github.com/cherrai/nyanyago-utils/saass"
	"github.com/cherrai/nyanyago-utils/validation"
	sso "github.com/cherrai/saki-sso-go"
	"github.com/gin-gonic/gin"
)

var ()

type FileController struct {
}

func (fc *FileController) GetUploadToken(c *gin.Context) {
	// 1、请求体
	var res response.ResponseProtobufType
	res.Code = 200

	// 2、获取参数
	data := new(protos.GetUploadToken_Request)
	var err error
	if err = protos.DecodeBase64(c.GetString("data"), data); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	// 3、验证参数
	if err = validation.ValidateStruct(
		data,
		validation.Parameter(&data.FileInfo, validation.Required()),
	); err != nil {
		res.Errors(err)
		res.Code = 10002
		res.Call(c)
		return
	}

	if err = validation.ValidateStruct(
		data.FileInfo,
		validation.Parameter(&data.FileInfo.Name, validation.Type("string"), validation.Required()),
		validation.Parameter(&data.FileInfo.Size, validation.Type("int64"), validation.Required()),
		validation.Parameter(&data.FileInfo.Type, validation.Type("string"), validation.Required()),
		validation.Parameter(&data.FileInfo.Suffix, validation.Type("string"), validation.Required()),
		validation.Parameter(&data.FileInfo.LastModified, validation.Type("int64"), validation.Required()),
		validation.Parameter(&data.FileInfo.Hash, validation.Type("string"), validation.Required()),
	); err != nil {
		res.Errors(err)
		res.Code = 10002
		res.Call(c)
		return
	}

	authorId := c.MustGet("userInfo").(*sso.UserInfo).Uid

	// 4、获取Token
	chunkSize := int64(128 * 1024)

	if data.FileInfo.Size < 1024*1024 {
		chunkSize = 128 * 1024
	}
	if data.FileInfo.Size > 1024*1024 {
		chunkSize = 256 * 1024
	}
	if data.FileInfo.Size > 15*1024*1024 {
		chunkSize = 512 * 1024
	}

	ut, err := conf.SAaSS.CreateChuunkUploadToken(&saass.CreateUploadTokenOptions{
		FileInfo: &saass.FileInfo{
			Name:         data.FileInfo.Name,
			Size:         data.FileInfo.Size,
			Type:         data.FileInfo.Type,
			Suffix:       data.FileInfo.Suffix,
			LastModified: data.FileInfo.LastModified,
			Hash:         data.FileInfo.Hash,
		},
		Path:           "/" + cipher.MD5(nstrings.ToString(authorId)) + "/files/" + time.Now().Format("2006/01/02") + "/",
		FileName:       strings.ToLower(cipher.MD5(data.FileInfo.Hash+nstrings.ToString(data.FileInfo.Size)+nstrings.ToString(time.Now().Unix()))) + data.FileInfo.Suffix,
		ChunkSize:      chunkSize,
		VisitCount:     -1,
		ExpirationTime: -1,
		// Type:           "File",
		FileConflict: "Replace",
		OnProgress: func(progress saass.Progress) {
			// log.Info("progress", progress)
		},
		OnSuccess: func(urls saass.Urls) {
			// log.Info("urls", urls)
		},
		OnError: func(err error) {
			// log.Info("err", err)
		},
	})
	if err != nil {
		res.Errors(err)
		res.Code = 10019
		res.Call(c)
		return
	}
	urls := protos.Urls{
		DomainUrl:     ut.Urls.DomainUrl,
		EncryptionUrl: ut.Urls.EncryptionUrl,
		Url:           ut.Urls.Url,
	}
	// log.Info("ChunkSize", ut)
	// log.Info("ChunkSize", ut.ChunkSize, chunkSize)
	protoData := &protos.GetUploadToken_Response{
		Urls:           &urls,
		ApiUrl:         ut.ApiUrl,
		Token:          ut.Token,
		ChunkSize:      ut.ChunkSize,
		UploadedOffset: ut.UploadedOffset,
	}

	res.Data = protos.Encode(protoData)

	res.Call(c)
}

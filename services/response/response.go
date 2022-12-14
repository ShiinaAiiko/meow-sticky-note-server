package response

import (
	"encoding/json"
	"time"

	"github.com/ShiinaAiiko/meow-sticky-note/server/protos"
	"github.com/cherrai/nyanyago-utils/cipher"
	"github.com/cherrai/nyanyago-utils/nrand"
	"github.com/cherrai/nyanyago-utils/nsocketio"
	"github.com/cherrai/nyanyago-utils/nstrings"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// var (
// 	log = nlog.New()
// )

type ResponseProtobufType struct {
	protos.ResponseType
}

func (res *ResponseProtobufType) Call(c *gin.Context) {
	var r ResponseType
	r.Code = res.Code
	r.Data = res.Data
	r.Msg = res.Msg
	r.CnMsg = res.CnMsg
	r.Error = res.Error
	r.RequestId = res.RequestId
	r.RequestTime = res.RequestTime
	r.Platform = res.Platform
	c.Set("protobuf", r.GetResponse())
}

func (res *ResponseProtobufType) Errors(err error) {
	if err != nil {
		res.Error = err.Error()
	}
}

// func (res *ResponseProtobufType) ProtoEncode(r *response.ResponseType) string {

// 	return protos.Encode(
// 		&protos.ResponseType{
// 			Code:        r.Code,
// 			Data:        r.Data.(string),
// 			Msg:         r.Msg,
// 			CnMsg:       r.CnMsg,
// 			Error:       r.Error,
// 			RequestId:   r.RequestId,
// 			RequestTime: r.RequestTime,
// 			Platform:    r.Platform,
// 			Author:      r.Author,
// 		},
// 	)
// }

func (res *ResponseProtobufType) GetResponse() interface{} {
	var r ResponseType
	r.Code = res.Code
	r.Data = res.Data
	r.Msg = res.Msg
	r.CnMsg = res.CnMsg
	r.Error = res.Error
	r.RequestId = res.RequestId
	r.RequestTime = res.RequestTime
	r.Platform = res.Platform
	return r.GetResponse()
}
func (res *ResponseProtobufType) CallSocketIo(c *nsocketio.EventInstance) {
	var r ResponseType
	r.Code = res.Code
	r.Data = res.Data
	r.Msg = res.Msg
	r.CnMsg = res.CnMsg
	r.Error = res.Error
	r.RequestId = res.RequestId
	r.RequestTime = res.RequestTime
	r.Platform = res.Platform
	// fmt.Println("r.GetResponse()", r.GetResponse())
	c.Set("protobuf", r.GetResponse())
}
func (res *ResponseProtobufType) ResponseProtoEncode() string {

	st := new(protos.ResponseType)
	responseData := res.GetResponse().(*ResponseType)

	copier.Copy(st, responseData)
	return protos.Encode(
		st,
	)
}

func (res *ResponseProtobufType) Encryption(userAesKey string, getReponseData interface{}) string {
	if getReponseData == nil {
		return ""
	}

	// fmt.Println("getReponseProtobufData", getReponseProtobufData)
	// ??????????????????????????????AesKey?????????
	// ?????????????????????????????????AesKey????????????key???

	aes := cipher.AES{
		Key:  "",
		Mode: "CFB",
	}
	if userAesKey == "" {
		aes.Key = cipher.MD5(nstrings.ToString(nrand.GetRandomNum(18)))
	} else {
		aes.Key = userAesKey
	}
	// log.Info("userAesKey", aes.Key, userAesKey, userAesKey != "")

	getResponseStr, _ := json.Marshal(getReponseData)
	bodyStr, _ := aes.Encrypt(string(getResponseStr), aes.Key)
	if userAesKey != "" {
		return protos.Encode(
			&protos.ResponseEncryptDataType{
				Data: bodyStr.HexEncodeToString(),
			},
		)
	}
	return protos.Encode(
		&protos.ResponseEncryptDataType{
			Data: bodyStr.HexEncodeToString(),
			Key:  aes.Key,
		},
	)
}

type ResponseType struct {
	// Code 200, 10004
	Code        int64  `json:"code,omitempty"`
	Msg         string `json:"msg,omitempty"`
	CnMsg       string `json:"cnMsg,omitempty"`
	Error       string `json:"error,omitempty"`
	RequestId   string `json:"requestId,omitempty"`
	RequestTime int64  `json:"requestTime,omitempty"`
	Author      string `json:"author,omitempty"`
	Platform    string `json:"platform,omitempty"`
	// RequestTime int64                  `json:"requestTime"`
	// Author      string                 `json:"author"`
	Data interface{} `json:"data,omitempty"`
}

type H map[string]interface{}
type Any *anypb.Any

func (res *ResponseType) Call(c *gin.Context) {

	// Log.Info("setResponse", res.GetResponse())
	c.Set("body", res.GetResponse())
	// fmt.Println("setResponse")
	// c.JSON(http.StatusOK, res.GetResponse())
}

func (res *ResponseType) GetResponse() *ResponseType {
	msg := res.Msg
	cnMsg := res.CnMsg
	if res.Msg == "" {
		res.Msg = "Request success."
	}
	if res.CnMsg == "" {
		res.CnMsg = "????????????"
	}
	if res.Platform == "" {
		res.Platform = "Meow Sticky Note"
	}
	if res.Author == "" {
		res.Author = "Shiina Aiiko."
	}
	res.RequestTime = time.Now().Unix()

	switch res.Code {
	case 200:

	case 10304:
		res.Msg = "Anonymous user login failed."
		res.CnMsg = "??????????????????????????????"

	case 10303:
		res.Msg = "Room join failed."
		res.CnMsg = "??????????????????"

	case 10302:
		res.Msg = "Invite code verification failed."
		res.CnMsg = "?????????????????????"

	case 10301:
		res.Msg = "Failed to create invitation code."
		res.CnMsg = "?????????????????????"

	case 10201:
		res.Msg = "Failed to send message."
		res.CnMsg = "??????????????????"

	case 10106:
		res.Msg = "Failed to reject friend."
		res.CnMsg = "??????????????????"

	case 10105:
		res.Msg = "Non-friends."
		res.CnMsg = "????????????"

	case 10104:
		res.Msg = "Failed to delete friend."
		res.CnMsg = "??????????????????"

	case 10103:
		res.Msg = "Failed to add friend."
		res.CnMsg = "??????????????????"

	case 10102:
		res.Msg = "Friend verification failed."
		res.CnMsg = "???????????????????????????"

	case 10101:
		res.Msg = "Failed to update friend log status."
		res.CnMsg = "??????????????????????????????"

		// case 10005:
		// 	res.Msg = "Invalid request."
		// 	res.CnMsg = "????????????"
		//

	case 10022:
		res.Msg = "Delete failed."
		res.CnMsg = "????????????"
	case 10021:
		res.Msg = "Content does not exist."
		res.CnMsg = "???????????????"

	case 10020:
		res.Msg = "Failed to get urls."
		res.CnMsg = "??????Urls??????"

	case 10019:
		res.Msg = "Create file token failed."
		res.CnMsg = "????????????Token??????"

	case 10018:
		res.Msg = "Chunksize is inconsistent."
		res.CnMsg = "??????????????????"

	case 10017:
		res.Msg = "The hash value is inconsistent."
		res.CnMsg = "Hash????????????"

	case 10016:
		res.Msg = "File upload error."
		res.CnMsg = "??????????????????"

	case 10015:
		res.Msg = "Failed to verify token."
		res.CnMsg = "Token????????????"

	case 10014:
		res.Msg = "App does not exist."
		res.CnMsg = "???????????????"

	case 10013:
		res.Msg = "Route does not exist."
		res.CnMsg = "???????????????"

	case 10012:
		res.Msg = "Already executed."
		res.CnMsg = "???????????????"

	case 10011:
		res.Msg = "Update failed."
		res.CnMsg = "????????????"

	case 10010:
		res.Msg = "Insufficient Privilege."
		res.CnMsg = "????????????."

	case 10009:
		res.Msg = "Decryption failed."
		res.CnMsg = "????????????."

	case 10008:
		res.Msg = "Encryption key error."
		res.CnMsg = "????????????."

	case 10007:
		res.Msg = "Encryption key generation failed."
		res.CnMsg = "????????????????????????"

	case 10006:
		res.Msg = "No more."
		res.CnMsg = "?????????????????????"

	case 10005:
		res.Msg = "Repeat request."
		res.CnMsg = "????????????"

	case 10004:
		res.Msg = "Login error."
		res.CnMsg = "??????????????????"

	case 10001:
		res.Msg = "Request error."
		res.CnMsg = "????????????"

	case 10002:
		res.Msg = "Parameter error."
		res.CnMsg = "????????????"

	default:
		res.Msg = "Request error."
		res.CnMsg = "????????????"

	}
	if res.Code == 0 {
		res.Code = 10001
	}

	if msg != "" {
		res.Msg = msg
	}
	if cnMsg != "" {
		res.CnMsg = cnMsg
	}

	return res
}

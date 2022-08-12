package socketioMiddleware

import (
	"github.com/cherrai/nyanyago-utils/nsocketio"
)

// 解密Request
func Decryption() nsocketio.HandlerFunc {
	return func(c *nsocketio.EventInstance) (err error) {
		// var res response.ResponseProtobufType
		// res.Code = 10008
		// enData := c.GetParamsString("data")
		// // requestId := c.GetParamsString("requestId")
		// // fmt.Println("en requestId", requestId)
		// // c.Set("requestId", requestId)
		// var userAesKey *encryption.UserAESKey
		// aes := cipher.AES{
		// 	Key:  "",
		// 	Mode: "CFB",
		// }
		// data := new(protos.RequestEncryptDataType)
		// // fmt.Println("DecryptionSoc中间件", enData)
		// dataBase64, dataBase64Err := base64.StdEncoding.DecodeString(enData)
		// if dataBase64Err != nil {
		// 	res.Data = "[Encryption]" + dataBase64Err.Error()
		// 	res.Code = 10008
		// 	res.CallSocketIo(c)
		// 	return
		// }
		// deErr := protos.Decode(dataBase64, data)
		// if deErr != nil {
		// 	res.Data = "[Encryption]" + deErr.Error()
		// 	res.Code = 10008
		// 	res.CallSocketIo(c)
		// 	return
		// }

		// // log.Info("开始处理加密内容")
		// getAesKey := conf.EncryptionClient.GetUserAesKeyByKey(conf.Redisdb, data.Key)
		// if getAesKey == nil {
		// 	res.Code = 10008
		// 	res.CallSocketIo(c)
		// 	return
		// }
		// userAesKey = getAesKey
		// aes.Key = getAesKey.AESKey
		// c.Set("userAesKey", userAesKey)

		// userInfo := c.GetSessionCache("userInfo").(*sso.UserInfo)

		// // log.Info(userAesKey, userInfo, userAesKey.Uid != userInfo.Uid || userAesKey.DeviceId != userInfo.UserAgent.DeviceId)
		// if userAesKey.Uid != userInfo.Uid || userAesKey.DeviceId != userInfo.UserAgent.DeviceId {
		// 	res.Code = 10008
		// 	res.CallSocketIo(c)
		// 	return
		// }
		// deStr, deStrErr := aes.DecryptWithString(data.Data, "")
		// if deStrErr != nil {
		// 	res.Data = "[Encryption]" + deStrErr.Error()
		// 	res.Code = 10008
		// 	res.CallSocketIo(c)
		// 	return
		// }
		// var dataMap map[string]interface{}
		// unErr := json.Unmarshal(deStr.Byte(), &dataMap)
		// if unErr != nil {
		// 	res.Data = "[Encryption]" + unErr.Error()
		// 	res.Code = 10008
		// 	res.CallSocketIo(c)
		// 	return
		// }
		// for key, item := range dataMap {
		// 	// fmt.Println(key, item)
		// 	c.Set(key, item)
		// }

		c.Next()
		return
	}
}

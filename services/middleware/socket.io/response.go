package socketioMiddleware

import (
	"github.com/cherrai/nyanyago-utils/nsocketio"
)

func Response() nsocketio.HandlerFunc {
	return func(c *nsocketio.EventInstance) error {
		defer func() {
			// var res response.ResponseProtobufType
			// log.Info("Response")
			getProtobufDataResponse, _ := c.Get("protobuf")
			log.Info("getProtobufDataResponse", getProtobufDataResponse)
			// userAesKeyInterface, exists := c.Get("userAesKey")
			// log.Info("userAesKeyInterface, exists", userAesKeyInterface, exists, !exists)
			// if !exists {
			// 	return
			// }
			// log.Info("userAesKeyInterface, exists", userAesKeyInterface, exists, !exists)
			// userAesKey := userAesKeyInterface.(*encryption.UserAESKey)
			// log.Info("userAesKey", userAesKey)
			requestId := c.GetParamsString("requestId")
			log.Info("requestId", requestId)
			// log.Info(res.Encryption(userAesKey, getProtobufDataResponse))
			// c.Emit(map[string]interface{}{
			// 	"data":      res.Encryption(userAesKey.AESKey, getProtobufDataResponse),
			// 	"requestId": requestId,
			// })

			// log.Info("getProtobufDataResponse", getProtobufDataResponse)
		}()

		c.Next()
		return nil
	}
}

package middleware

import (
	"encoding/json"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/response"

	sso "github.com/cherrai/saki-sso-go"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, isStaticServer := c.Get("isStaticServer"); isStaticServer {
			c.Next()
			return
		}
		if _, isWsServer := c.Get("WsServer"); isWsServer {
			c.Next()
			return
		}

		roles := new(RoleOptionsType)
		getRoles, isRoles := c.Get("roles")
		if isRoles {
			roles = getRoles.(*RoleOptionsType)
		}

		// Log.Info("------Authorize------", roles.Authorize)

		if roles.Authorize {
			res := response.ResponseProtobufType{}
			res.Code = 10004

			// 解析用户数据
			token := c.GetString("token")
			deviceId := c.GetString("deviceId")
			userAgent := new(sso.UserAgent)
			userAgentAny, isUserAgentAny := c.Get("userAgent")
			if !isUserAgentAny {
				res.Code = 10004
				res.Call(c)
				c.Abort()
				return
			}
			if token == "" || deviceId == "" || userAgent == nil {
				res.Code = 10004
				res.Call(c)
				c.Abort()
				return
			}
			userAgent = userAgentAny.(*sso.UserAgent)
			// Log.Info("token, deviceId, userAgent", deviceId, *userAgent)
			ret, err := conf.SSO.Verify(token, deviceId, *userAgent)
			if err != nil {
				// Log.Info("jwt: ", err)
				res.Call(c)
				c.Abort()
				return
			}
			// log.Info("token", token)
			// log.Info("deviceId", deviceId)
			// log.Info("userAgent", userAgent)

			// Log.Info("ret", ret, ret.Payload)
			if ret != nil && ret.Payload.Uid != 0 {

				// if isExchangeKey := strings.Contains(c.Request.URL.Path, "encryption/exchangeKey"); !isExchangeKey {
				// 	// 要求登录的同时还没有key就说不过去了
				// 	// userAesKeyInterface, err := c.Get("userAesKey")
				// 	// if userAesKeyInterface != nil || !err {
				// 	// 	userAesKey := userAesKeyInterface.(*encryption.UserAESKey)
				// 	// 	if userAesKey.Uid != ret.Payload.Uid || userAesKey.DeviceId != ret.Payload.UserAgent.DeviceId {
				// 	// 		res.Code = 10008
				// 	// 		res.Call(c)
				// 	// 		c.Abort()
				// 	// 		return
				// 	// 	}
				// 	// }
				// }
				c.Set("userInfo", &ret.Payload)

				c.Next()
				return
			}
			res.Code = 10004
			res.Call(c)
			// Log.Info(res)
			c.Abort()
			// res.Call(c)
			// c.Abort()
			return
		}

		c.Next()
	}
}

func ConvertResponseJson(jsonStr []byte) (sso.UserInfo, error) {
	var m sso.UserInfo
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		Log.Info("Unmarshal with error: %+v\n", err)
		return m, err
	}
	return m, nil
}

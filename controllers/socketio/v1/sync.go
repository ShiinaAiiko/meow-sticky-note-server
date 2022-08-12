package socketIoControllersV1

import (

	// "github.com/cherrai/saki-sso-go"

	"errors"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/response"
	"github.com/cherrai/nyanyago-utils/nsocketio"
	"github.com/cherrai/nyanyago-utils/nstrings"
	sso "github.com/cherrai/saki-sso-go"
)

type SyncControllers struct {
}

func (s *SyncControllers) SyncData(e *nsocketio.EventInstance) error {
	// Conn := e.Conn()
	// c := e.ConnContext()
	log.Info("------SyncData------")

	return nil
}

func (s *SyncControllers) NewConnect(e *nsocketio.EventInstance) error {
	log.Info("/Chat 开始连接")
	var res response.ResponseProtobufType
	c := e.ConnContext()

	Conn := e.Conn()

	userInfoAny := e.GetSessionCache("userInfo")
	if userInfoAny == nil {
		return errors.New("userinfo does not exist")
	}
	userInfo := userInfoAny.(*sso.UserInfo)

	// getUser, err := conf.SSO.Verify(queryData.Token, queryData.DeviceId, queryData.UserAgent)
	if userInfo.Uid == 0 {
		res.Code = 10004
		res.Data = "userinfo does not exist"
		Conn.Emit(conf.SocketRouterEventNames["Error"], res.GetResponse())
		defer Conn.Close()
		return nil
	} else {
		log.Warn(e.Namespace()+" => 连接成功！", Conn.ID(), nstrings.ToString(userInfo.Uid)+", Connection to Successful.")
		// c.SetRoomId(nstrings.ToString(getUser.Payload.Uid), getUser.Payload.UserAgent.DeviceId)
		// c.SetCustomId(cipher.MD5(nstrings.ToString(getUser.Payload.Uid) + getUser.Payload.UserAgent.DeviceId))
		c.SetTag("Uid", nstrings.ToString(userInfo.Uid))
		c.SetTag("DeviceId", userInfo.UserAgent.DeviceId)
	}
	return nil

	// var res response.ResponseProtobufType

	// // fmt.Println(msgRes)
	// defer func() {
	// 	// fmt.Println("Error middleware.2222222222222")
	// 	if err := recover(); err != nil {
	// 		res.Code = 10001
	// 		res.Data = err.(error).Error()
	// 		Conn.Emit(conf.SocketRouterEventNames["Error"], res.GetReponse())
	// 		defer Conn.Close()
	// 	}
	// }()
	// sc := methods.SocketConn{
	// 	Conn: Conn,
	// }
	// // s.SetContext("dsdsdsd")
	// // 申请一个房间

	// query := new(typings.SocketEncryptionQuery)

	// err := qs.Unmarshal(query, Conn.URL().RawQuery)

	// if err != nil {
	// 	res.Code = 10002
	// 	res.Data = err.Error()
	// 	Conn.Emit(conf.SocketRouterEventNames["Error"], res.GetReponse())
	// 	defer Conn.Close()
	// 	return err
	// }
	// sc.Query = query
	// // log.Info("query", query)

	// queryData := new(typings.SocketQuery)
	// deQueryDataErr := sc.Decryption(queryData)
	// // log.Info("deQueryDataErr", deQueryDataErr != nil, deQueryDataErr)
	// if deQueryDataErr != nil {
	// 	res.Code = 10009
	// 	res.Data = deQueryDataErr.Error()
	// 	Conn.Emit(conf.SocketRouterEventNames["Error"], res.GetReponse())
	// 	defer Conn.Close()
	// 	return deQueryDataErr
	// }

	// getUser, err := conf.SSO.Verify(queryData.Token, queryData.DeviceId, queryData.UserAgent)
	// if err != nil || getUser == nil || getUser.Payload.Uid == 0 {
	// 	res.Code = 10004
	// 	res.Data = "SSO Error: " + err.Error()
	// 	Conn.Emit(conf.SocketRouterEventNames["Error"], res.GetReponse())
	// 	defer Conn.Close()
	// 	return err
	// } else {

	// 	// 检测之前是否登录过了，登录过把之前的实例干掉
	// 	for _, v := range conf.SocketIO.GetConnContextByTag(conf.SocketRouterNamespace["Base"], "DeviceId", getUser.Payload.UserAgent.DeviceId) {
	// 		// 1、发送信息告知对方下线
	// 		if userAesKey := conf.EncryptionClient.GetUserAesKeyByDeviceId(conf.Redisdb, getUser.Payload.UserAgent.DeviceId); userAesKey != nil {
	// 			var res response.ResponseProtobufType
	// 			res.Code = 200
	// 			res.Data = protos.Encode(&protos.OnForceOffline_Response{})
	// 			eventName := conf.SocketRouterEventNames["OnForceOffline"]
	// 			responseData := res.Encryption(userAesKey.AESKey, res.GetReponse())
	// 			isEmit := v.Emit(eventName, responseData)
	// 			if isEmit {
	// 				// 2、断开连接
	// 				log.Info("有另外一个设备在线", userAesKey)
	// 				// go v.Close()
	// 			} else {
	// 			}
	// 		}

	// 	}
	// 	log.Info("UID " + strconv.FormatInt(getUser.Payload.Uid, 10) + ", Connection to Successful.")
	// 	log.Info("/ UID", getUser.Payload.Uid)

	// 	c.SetSessionCache("loginTime", time.Now().Unix())
	// 	c.SetSessionCache("userInfo", &getUser.Payload)
	// 	c.SetSessionCache("deviceId", queryData.DeviceId)
	// 	c.SetSessionCache("userAgent", &queryData.UserAgent)
	// 	c.SetTag("Uid", nstrings.ToString(getUser.Payload.Uid))
	// 	c.SetTag("DeviceId", getUser.Payload.UserAgent.DeviceId)

	// 	// log.Info("SocketIO Client连接成功：", Conn.ID())

	// 	// sc.SetUserInfo(&ret.Payload)

	// 	// log.Info("------ 1、检测其他设备是否登录------")
	// 	// 1、检测其他设备是否登录
	// 	sc := e.ServerContext()
	// 	getConnContext := sc.GetConnContextByTag(conf.SocketRouterNamespace["Base"], "Uid", nstrings.ToString(getUser.Payload.Uid))
	// 	// log.Info("当前ID", c.ID())
	// 	// log.Info("有哪些设备在线", getConnContext)

	// 	onlineDeviceListMap := map[string]*protos.OnlineDeviceList{}
	// 	onlineDeviceList := []*protos.OnlineDeviceList{}
	// 	// 2、遍历设备实例、告诉对方上线了
	// 	for _, cctx := range getConnContext {
	// 		// log.Info(c)
	// 		// uid := cctx.GetTag("Uid")
	// 		// log.Info(uid)
	// 		deviceId := cctx.GetTag("DeviceId")
	// 		// log.Info(deviceId)

	// 		// userInfo
	// 		cctxSsoUser := new(protos.SSOUserInfo)
	// 		cctxUserInfoInteface := cctx.GetSessionCache("userInfo")
	// 		if cctxUserInfoInteface == nil {
	// 			continue
	// 		}
	// 		cctxUserInfo := cctxUserInfoInteface.(*sso.UserInfo)
	// 		copier.Copy(cctxSsoUser, cctxUserInfo)

	// 		// userAgent
	// 		cctxProtoUserAgent := new(protos.UserAgent)
	// 		cctxUserAgentInteface := cctx.GetSessionCache("userAgent")
	// 		if cctxUserAgentInteface == nil {
	// 			continue
	// 		}
	// 		cctxUserAgent := cctxUserAgentInteface.(*sso.UserAgent)
	// 		copier.Copy(cctxProtoUserAgent, cctxUserAgent)

	// 		// loginTime
	// 		cctxLoginTimeInterface := cctx.GetSessionCache("loginTime")
	// 		if cctxLoginTimeInterface == nil {
	// 			continue
	// 		}
	// 		onlineDeviceListMap[deviceId] = &protos.OnlineDeviceList{
	// 			UserInfo:  cctxSsoUser,
	// 			LoginTime: cctxLoginTimeInterface.(int64),
	// 			UserAgent: cctxProtoUserAgent,
	// 			Location:  "",
	// 			DeviceId:  deviceId,
	// 		}
	// 		onlineDeviceList = append(onlineDeviceList, onlineDeviceListMap[deviceId])
	// 	}
	// 	// log.Info("onlineDeviceList", onlineDeviceListMap)

	// 	currentDevice := onlineDeviceListMap[queryData.DeviceId]
	// 	for _, cctx := range getConnContext {
	// 		deviceId := cctx.GetTag("DeviceId")
	// 		// log.Info(deviceId)

	// 		// if deviceId == queryData.DeviceId {
	// 		// 	log.Info("乃是自己也")
	// 		// }else{

	// 		// }
	// 		// userAesKey1 := conf.EncryptionClient.GetUserAesKeyByDeviceId(conf.Redisdb, deviceId)

	// 		if userAesKey := conf.EncryptionClient.GetUserAesKeyByDeviceId(conf.Redisdb, deviceId); userAesKey != nil {
	// 			// log.Info("userAesKey SendJoinAnonymousRoomMessage", userAesKey)

	// 			var res response.ResponseProtobufType
	// 			res.Code = 200

	// 			res.Data = protos.Encode(&protos.OtherDeviceOnline_Response{
	// 				CurrentDevice:    currentDevice,
	// 				OnlineDeviceList: onlineDeviceList,
	// 			})

	// 			eventName := conf.SocketRouterEventNames["OtherDeviceOnline"]
	// 			responseData := res.Encryption(userAesKey.AESKey, res.GetReponse())
	// 			isEmit := cctx.Emit(eventName, responseData)
	// 			if isEmit {
	// 				// 发送成功或存储到数据库
	// 			} else {
	// 				// 存储到数据库作为离线数据
	// 			}
	// 		}

	// 	}
	// }

	return nil
}

func (s *SyncControllers) Disconnect(e *nsocketio.EventInstance) error {
	// c := e.ConnContext()

	log.Info("已经断开了", e.Reason)

	// // 1、检测其他设备是否登录
	// sc := e.ServerContext()

	// getConnContext := sc.GetConnContextByTag(conf.SocketRouterNamespace["Base"], "Uid", c.GetTag("Uid"))
	// log.Info("当前ID", c.ID())
	// log.Info("有哪些设备在线", getConnContext)

	// onlineDeviceListMap := map[string]*protos.OnlineDeviceList{}
	// onlineDeviceList := []*protos.OnlineDeviceList{}
	// // 2、遍历设备实例、告诉对方上线了
	// for _, cctx := range getConnContext {
	// 	// log.Info(c)
	// 	// uid := cctx.GetTag("Uid")
	// 	// log.Info(uid)
	// 	deviceId := cctx.GetTag("DeviceId")
	// 	// log.Info(deviceId)

	// 	// userInfo
	// 	cctxSsoUser := new(protos.SSOUserInfo)
	// 	cctxUserInfoInteface := cctx.GetSessionCache("userInfo")
	// 	if cctxUserInfoInteface == nil {
	// 		continue
	// 	}
	// 	cctxUserInfo := cctxUserInfoInteface.(*sso.UserInfo)
	// 	copier.Copy(cctxSsoUser, cctxUserInfo)

	// 	// userAgent
	// 	cctxProtoUserAgent := new(protos.UserAgent)
	// 	cctxUserAgentInteface := cctx.GetSessionCache("userAgent")
	// 	if cctxUserAgentInteface == nil {
	// 		continue
	// 	}
	// 	cctxUserAgent := cctxUserAgentInteface.(*sso.UserAgent)
	// 	copier.Copy(cctxProtoUserAgent, cctxUserAgent)

	// 	// loginTime
	// 	cctxLoginTimeInterface := cctx.GetSessionCache("loginTime")
	// 	if cctxLoginTimeInterface == nil {
	// 		continue
	// 	}
	// 	onlineDeviceListMap[deviceId] = &protos.OnlineDeviceList{
	// 		UserInfo:  cctxSsoUser,
	// 		LoginTime: cctxLoginTimeInterface.(int64),
	// 		UserAgent: cctxProtoUserAgent,
	// 		Location:  "",
	// 		DeviceId:  deviceId,
	// 	}
	// 	onlineDeviceList = append(onlineDeviceList, onlineDeviceListMap[deviceId])
	// }

	// deviceId := c.GetSessionCache("deviceId")
	// if deviceId == nil {
	// 	return nil
	// }
	// currentDevice := onlineDeviceListMap[deviceId.(string)]

	// for _, cctx := range getConnContext {
	// 	deviceId := cctx.GetTag("DeviceId")
	// 	// log.Info(deviceId)

	// 	// if deviceId == queryData.DeviceId {
	// 	// 	log.Info("乃是自己也")
	// 	// }else{

	// 	// }
	// 	// userAesKey1 := conf.EncryptionClient.GetUserAesKeyByDeviceId(conf.Redisdb, deviceId)

	// 	if userAesKey := conf.EncryptionClient.GetUserAesKeyByDeviceId(conf.Redisdb, deviceId); userAesKey != nil {
	// 		// log.Info("userAesKey SendJoinAnonymousRoomMessage", userAesKey)

	// 		var res response.ResponseProtobufType
	// 		res.Code = 200

	// 		res.Data = protos.Encode(&protos.OtherDeviceOffline_Response{
	// 			CurrentDevice:    currentDevice,
	// 			OnlineDeviceList: onlineDeviceList,
	// 		})

	// 		eventName := conf.SocketRouterEventNames["OtherDeviceOffline"]
	// 		responseData := res.Encryption(userAesKey.AESKey, res.GetReponse())
	// 		isEmit := cctx.Emit(eventName, responseData)
	// 		if isEmit {
	// 			// 发送成功或存储到数据库
	// 		} else {
	// 			// 存储到数据库作为离线数据
	// 		}
	// 	}

	// }

	return nil
}

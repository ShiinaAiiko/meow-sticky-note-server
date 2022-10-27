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

}

func (s *SyncControllers) Disconnect(e *nsocketio.EventInstance) error {
	// c := e.ConnContext()

	log.Info("已经断开了", e.Reason)

	return nil
}

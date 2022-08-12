package socketioRouter

import (
	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	socketIoControllersV1 "github.com/ShiinaAiiko/meow-sticky-note/server/controllers/socketio/v1"
)

func (v V1) InitSync() {
	r := v.Router

	sc := socketIoControllersV1.SyncControllers{}

	// s.OnConnect(r.Chat, func(s socketio.Conn) error {
	// 	fmt.Println(r.Chat+"连接成功：", s.ID())
	// 	return nil
	// })
	v.Server.Router(r.Sync, conf.SocketRouterEventNames["SyncData"], sc.SyncData)

	// r := v.Router
	v.Server.OnConnect(r.Sync, sc.NewConnect)
	v.Server.OnDisconnect(r.Base, sc.Disconnect)

}

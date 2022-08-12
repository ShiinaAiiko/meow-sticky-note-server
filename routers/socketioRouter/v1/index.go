package socketioRouter

import (
	socketIoControllersV1 "github.com/ShiinaAiiko/meow-sticky-note/server/controllers/socketio/v1"

	"github.com/cherrai/nyanyago-utils/nsocketio"
)

type V1 struct {
	Server *nsocketio.NSocketIO
	Router RouterV1
}

type RouterV1 struct {
	Sync string
	Base string
}

func (v V1) Init() {
	r := v.Router

	bc := socketIoControllersV1.BaseControllers{}
	// s.OnConnect(r.Chat, func(s socketio.Conn) error {
	// 	fmt.Println(r.Chat+"连接成功：", s.ID())
	// 	return nil
	// })

	// r := v.Router
	v.Server.OnConnect(r.Base, bc.NewConnect)
	v.Server.OnDisconnect(r.Base, bc.Disconnect)

	v.InitSync()
}

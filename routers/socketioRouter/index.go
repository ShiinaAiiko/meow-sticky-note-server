package socketioRouter

import (
	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	"github.com/ShiinaAiiko/meow-sticky-note/server/routers/socketioRouter/v1"
)

func InitRouter() {
	// fmt.Println(conf.SocketIoServer)
	rv1 := socketioRouter.V1{
		Server: conf.SocketIO,
		Router: socketioRouter.RouterV1{
			Sync: "/sync",
			Base: "/",
		},
	}
	rv1.Init()
}

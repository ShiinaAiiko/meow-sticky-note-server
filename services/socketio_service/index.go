package socketio_service

import (
	"github.com/ShiinaAiiko/meow-sticky-note/server/routers/socketioRouter"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/gin_service"

	socketioMiddleware "github.com/ShiinaAiiko/meow-sticky-note/server/services/middleware/socket.io"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	"github.com/cherrai/nyanyago-utils/nlog"
	"github.com/cherrai/nyanyago-utils/nsocketio"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var (
	log = nlog.New()
)

// var Router *gin.Engine
// 个人也是一个room，roomId：U+UID
// 群组也是一个room，roomId：G+UID
// ChatMessage 发送消息
// 直接发给对应的roomId即可

// InputStatus 发送正在输入状态
// 直接发给对应的roomId即可

// OnlineStatus 在线状态
// 发送给好友关系存在的roomId

func Init() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	gin.SetMode(conf.Config.Server.Mode)

	gin_service.Router = gin.New()

	// fmt.Println("Server", Server)
	// conf.SocketIoServer = socketIoServer

	conf.SocketIO = nsocketio.New(&nsocketio.Options{
		RDB: conf.Redisdb,
		RedisAdapterOptions: &socketio.RedisAdapterOptions{
			Addr:    conf.Config.Redis.Addr,
			Prefix:  "socket.io",
			Network: "tcp",
		},
	})

	// 处理中间件
	conf.SocketIO.Use(socketioMiddleware.RoleMiddleware())
	conf.SocketIO.Use(socketioMiddleware.Response())
	conf.SocketIO.Use(socketioMiddleware.Error())
	conf.SocketIO.Use(socketioMiddleware.Decryption())

	socketioRouter.InitRouter()
	log.Info("[Socket.IO] server created successfully.")

}

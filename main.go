package main

import (
	"context"
	"os"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	mongodb "github.com/ShiinaAiiko/meow-sticky-note/server/db/mongo"
	redisdb "github.com/ShiinaAiiko/meow-sticky-note/server/db/redis"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/gin_service"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/socketio_service"

	"github.com/cherrai/nyanyago-utils/nlog"
	"github.com/cherrai/nyanyago-utils/nredis"
	"github.com/cherrai/nyanyago-utils/saass"
	sso "github.com/cherrai/saki-sso-go"

	// sfu "github.com/pion/ion-sfu/pkg/sfu"

	"github.com/go-redis/redis/v8"
)

var (
	log = nlog.New()
)

// 文件到期后根据时间进行删除 未做
func main() {
	nlog.SetPrefixTemplate("[{{Timer}}] [{{Type}}] [{{Date}}] [{{File}}]@{{Name}}")
	nlog.SetName("SAaSS")

	conf.G.Go(func() {
		configPath := ""
		for k, v := range os.Args {
			switch v {
			case "--config":
				if os.Args[k+1] != "" {
					configPath = os.Args[k+1]
				}

			}
		}
		if configPath == "" {
			log.Error("Config file does not exist.")
			return
		}
		conf.GetConfig(configPath)

		// Connect to redis.
		redisdb.ConnectRedis(&redis.Options{
			Addr:     conf.Config.Redis.Addr,
			Password: conf.Config.Redis.Password, // no password set
			DB:       conf.Config.Redis.DB,       // use default DB
		})

		conf.Redisdb = nredis.New(context.Background(), &redis.Options{
			Addr:     conf.Config.Redis.Addr,
			Password: conf.Config.Redis.Password, // no password set
			DB:       conf.Config.Redis.DB,       // use default DB
		}, conf.BaseKey)
		conf.Redisdb.CreateKeys(conf.RedisCacheKeys)

		conf.SSO = sso.New(&sso.SakiSsoOptions{
			AppId:  conf.Config.SSO.AppId,
			AppKey: conf.Config.SSO.AppKey,
			Host:   conf.Config.SSO.Host,
			RedisOptions: &redis.Options{
				Addr:     conf.Config.Redis.Addr,
				Password: conf.Config.Redis.Password,
				DB:       conf.Config.Redis.DB,
			},
		})

		conf.SAaSS = saass.New(&saass.Options{
			AppId:      "1e816914-64d2-477a-8e35-427d947ecf50",
			AppKey:     "0cbd470b-5091-418a-8cfd-4349404879b9",
			BaseUrl:    conf.Config.StaticPathDomain,
			ApiVersion: "v1",
		})
		// Connect to mongodb.
		mongodb.ConnectMongoDB(conf.Config.Mongodb.Currentdb.Uri, conf.Config.Mongodb.Currentdb.Name)

		// syncDbx := new(dbxV1.SyncDbx)

		// debounce := nutils.NewDebounce(1000 * time.Millisecond)

		// g, ctx := errgroup.WithContext(context.Background())

		// g.Go(func() error {
		// 	result := make(chan string)
		// 	go func() {
		// 		fmt.Println("will sleep 3s")
		// 		time.Sleep(3 * time.Second)
		// 		fmt.Println("sleep 3s done")
		// 		result <- "done888"
		// 	}()

		// 	select {
		// 	case <-ctx.Done():
		// 		fmt.Println("3333333")
		// 		return nil
		// 		// return fmt.Errorf("errgroup will exit")
		// 	case r := <-result:
		// 		fmt.Printf("1res: %s\n", r)
		// 		return nil
		// 	}
		// })

		// g.Go(func() error {
		// 	fmt.Println("will sleep 1s")
		// 	time.Sleep(1 * time.Second)
		// 	// 返回错误一秒后主程序就会推出
		// 	// return fmt.Errorf("err3")
		// 	// 不返回错误，另外一个 协程会执行完毕
		// 	return nil
		// })

		// g.Go(func() error {
		// 	defer func() error {
		// 		if r := recover(); r != nil {
		// 			log.Error(r)
		// 			return r.(error)
		// 		}
		// 		return nil
		// 	}()
		// 	panic(errors.New("panicpanicpanic"))
		// 	return nil
		// })

		// err := g.Wait()

		// if err != nil {
		// 	log.Error("err group wait: %v", err)
		// } else {
		// 	log.Error("err group wait all success")
		// }

		// result := make(chan string)

		// go func() {
		// 	for data := range result {
		// 		// 打印通道数据
		// 		fmt.Println("dsadsad", data)
		// 		// 当遇到数据0时, 退出接收循环
		// 		// if data == 0 {
		// 		// 	break
		// 		// }
		// 	}
		// }()
		// PanicGoroutine(f func()){

		// }

		// go func() error {
		// 	return nil
		// }()

		// g := goroutinepanic.G

		// g.Go(func() {
		// 	time.Sleep(1 * time.Second)
		// 	panic(errors.New("panicpanicpanic1"))
		// })
		// g.Go(func() {
		// 	panic(errors.New("11111"))
		// })
		// log.Info(12345)

		// g.Wait()
		// ntimer.SetTimeout1(func() {
		// 	log.Info("sadsadsa")
		// }, 1000)
		// ntimer.SetTimeout(func() {
		// 	note, _ := syncDbx.GetNoteBySAaSS("0f82932b-d569-532b-b874-f137d5bc4686", 1)
		// 	log.Info(note)
		// 	note.Name = "A30002"

		// 	syncDbx.UpdateNote(note, note.Id, 1)

		// 	// ntimer.SetTimeout(func() {
		// 	// 	debounce.Add(func() {
		// 	// 		log.Info("sasasasa1")
		// 	// 	})
		// 	// }, 300)
		// 	// ntimer.SetTimeout(func() {
		// 	// 	debounce.Add(func() {
		// 	// 		log.Info("sasasasa2")
		// 	// 	})
		// 	// }, 500)
		// 	// ntimer.SetTimeout(func() {
		// 	// 	debounce.Add(func() {
		// 	// 		log.Info("sasasasa3")
		// 	// 	})
		// 	// }, 700)
		// }, 500)
		socketio_service.Init()
		gin_service.Init()
	})

	conf.G.Error(func(err error) {
		log.Error(err)
	})
	conf.G.Wait()
}

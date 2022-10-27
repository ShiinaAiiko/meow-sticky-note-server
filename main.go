package main

import (
	"context"
	"os"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
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
			AppId:      conf.Config.Saass.AppId,
			AppKey:     conf.Config.Saass.AppKey,
			BaseUrl:    conf.Config.Saass.BaseUrl,
			ApiVersion: conf.Config.Saass.ApiVersion,
		})
		socketio_service.Init()
		gin_service.Init()
	})

	conf.G.Error(func(err error) {
		log.Error(err)
	})
	conf.G.Wait()
}

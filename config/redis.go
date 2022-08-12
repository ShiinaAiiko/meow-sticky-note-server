package conf

import (
	"time"

	"github.com/cherrai/nyanyago-utils/nredis"
)

var Redisdb *nredis.NRedis

var BaseKey = "meow-whisper"

var RedisCacheKeys = map[string]*nredis.RedisCacheKeysType{
	"Token": {
		Key:        "GetFriend",
		Expiration: 5 * 60 * time.Second,
	},
	"GetFriendIds": {
		Key:        "GetFriendIds",
		Expiration: 5 * 60 * time.Second,
	},
	"GetNote": {
		Key:        "GetNote",
		Expiration: 10 * 60 * time.Second,
	},
}

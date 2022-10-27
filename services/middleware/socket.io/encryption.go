package socketioMiddleware

import (
	"github.com/cherrai/nyanyago-utils/nsocketio"
)

// 解密Request
func Decryption() nsocketio.HandlerFunc {
	return func(c *nsocketio.EventInstance) (err error) {
		c.Next()
		return
	}
}

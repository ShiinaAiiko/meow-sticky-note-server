package routerV1

import (
	controllersV1 "github.com/ShiinaAiiko/meow-sticky-note/server/controllers/v1"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/middleware"
)

func (r Routerv1) InitFile() {
	fc := new(controllersV1.FileController)

	role := middleware.RoleMiddlewareOptions{
		BaseUrl: r.BaseUrl,
	}
	r.Group.GET(
		role.SetRole("/file/getUploadToken", &middleware.RoleOptionsType{
			CheckApp:           true,
			Authorize:          true,
			RequestEncryption:  false,
			ResponseEncryption: false,
			ResponseDataType:   "protobuf",
		}),
		fc.GetUploadToken)
}

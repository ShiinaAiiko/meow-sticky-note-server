package routerV1

import (
	controllersV1 "github.com/ShiinaAiiko/meow-sticky-note/server/controllers/v1"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/middleware"
)

func (r Routerv1) InitSync() {
	sc := new(controllersV1.SyncController)

	role := middleware.RoleMiddlewareOptions{
		BaseUrl: r.BaseUrl,
	}

	r.Group.POST(
		role.SetRole("/sync/toserver", &middleware.RoleOptionsType{
			CheckApp:           true,
			Authorize:          true,
			RequestEncryption:  false,
			ResponseEncryption: false,
			ResponseDataType:   "protobuf",
		}),
		sc.SyncToServer)

	r.Group.GET(
		role.SetRole("/sync/geturls", &middleware.RoleOptionsType{
			CheckApp:           true,
			Authorize:          true,
			RequestEncryption:  false,
			ResponseEncryption: false,
			ResponseDataType:   "protobuf",
		}),
		sc.GetUrls)

	r.Group.GET(
		role.SetRole("/sync/getfolderfiles", &middleware.RoleOptionsType{
			CheckApp:           true,
			Authorize:          true,
			RequestEncryption:  false,
			ResponseEncryption: false,
			ResponseDataType:   "protobuf",
		}),
		sc.GetFolderFiles)

	r.Group.GET(
		role.SetRole("/sync/getnote", &middleware.RoleOptionsType{
			CheckApp:           true,
			Authorize:          true,
			RequestEncryption:  false,
			ResponseEncryption: false,
			ResponseDataType:   "protobuf",
		}),
		sc.GetNote)
}

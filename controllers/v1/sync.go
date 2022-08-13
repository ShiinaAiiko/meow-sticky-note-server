package controllersV1

import (
	"strings"
	"time"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	dbxV1 "github.com/ShiinaAiiko/meow-sticky-note/server/dbx/v1"
	"github.com/ShiinaAiiko/meow-sticky-note/server/protos"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/methods"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/response"
	"github.com/cherrai/nyanyago-utils/nlog"
	"github.com/cherrai/nyanyago-utils/nstrings"
	"github.com/cherrai/nyanyago-utils/validation"
	sso "github.com/cherrai/saki-sso-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

var (
	log     = nlog.New()
	syncDbx = new(dbxV1.SyncDbx)
)

type SyncController struct {
}

// 目前暂时仅支持PNG和JPG
func (fc *SyncController) SyncToServer(c *gin.Context) {
	// 1、请求体
	log.Info("------SyncToServer------")
	// var res response.ResponseType
	var res response.ResponseProtobufType
	res.Code = 200

	// 2、获取参数
	data := new(protos.SyncToServer_Request)
	var err error
	if err = protos.DecodeBase64(c.GetString("data"), data); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	// 3、验证参数
	if err = validation.ValidateStruct(
		data,
		validation.Parameter(&data.Type, validation.Type("string"), validation.Enum([]string{"Note", "Category", "Page"}), validation.Required()),
		validation.Parameter(&data.Methods, validation.Type("string"), validation.Enum([]string{"Add", "Update", "Delete", "Sort"}), validation.Required()),
	); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	if !(data.Type == "Note" && data.Methods == "Sort") {
		if err = validation.ValidateStruct(
			data.Options,
			validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
		); err != nil {
			res.Error = err.Error()
			res.Code = 10002
			res.Call(c)
			return
		}
	}
	authorId := c.MustGet("userInfo").(*sso.UserInfo).Uid
	deviceId := c.MustGet("deviceId").(string)

	log.Info("add", data.Type, data.Methods, data.Options)

	// 4 获取Note
	var note *protos.NoteItem

	if !(data.Type == "Note" && (data.Methods == "Add" || data.Methods == "Sort")) {
		note, err = syncDbx.GetNoteBySAaSS(data.Options.NoteId, authorId)
		if err != nil {
			res.Error = err.Error()
			res.Code = 10021
			res.Call(c)
			return
		}
	}
	// log.Info(note)
	// note.Name = "A30002"

	switch data.Type {
	case "Note":
		if data.Methods == "Add" || data.Methods == "Update" {

			if err = validation.ValidateStruct(
				data,
				validation.Parameter(&data.Data, validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			if err = validation.ValidateStruct(
				data.Data,
				validation.Parameter(&data.Data.Note, validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			if err = validation.ValidateStruct(
				data.Data.Note,
				validation.Parameter(&data.Data.Note.Name, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Data.Note.LastUpdateTime, validation.Type("int64"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
		}
		switch data.Methods {
		case "Add":
			err = syncDbx.AddNote(data.Data.Note, data.Data.Note.Id, authorId)
			if err != nil {
				res.Error = err.Error()
				res.Code = 10011
				res.Call(c)
				return
			}

		case "Update":
			note.Name = data.Data.Note.Name
			note.LastUpdateTime = data.Data.Note.LastUpdateTime
		case "Delete":
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			err = syncDbx.DeleteNote(data.Options.NoteId, authorId)
			if err != nil {
				res.Error = err.Error()
				res.Code = 10022
				res.Call(c)
				return
			}
		case "Sort":
			if err = validation.ValidateStruct(
				data.Data,
				validation.Parameter(&data.Data.Sort, validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			for _, v := range data.Data.Sort.List {
				note, err = syncDbx.GetNoteBySAaSS(v.Id, authorId)
				if err != nil {
					res.Errors(err)
					res.Code = 10011
					res.Call(c)
					return
				}
				note.Sort = v.Sort
				note.LastUpdateTime = v.LastUpdateTime
				// log.Info(note.Name, v.Sort, v.LastUpdateTime)
				err = syncDbx.UpdateNote(note, note.Id, authorId)
				if err != nil {
					res.Errors(err)
					res.Code = 10011
					res.Call(c)
					return
				}
			}
		}

	case "Category":
		if data.Methods == "Add" || data.Methods == "Update" {
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.CategoryId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			if err = validation.ValidateStruct(
				data.Data.Category,
				validation.Parameter(&data.Data.Category.Name, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Data.Category.LastUpdateTime, validation.Type("int64"), validation.Required()),
				// validation.Parameter(&data.Data.Page.Sort, validation.Type("int64"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
		}
		switch data.Methods {
		case "Add":
			note.Categories = append(note.Categories, data.Data.Category)
		case "Update":
			for _, v := range note.Categories {
				if v.Id != data.Options.CategoryId {
					continue
				}
				log.Info("data.Data.Category.Name", data.Data.Category.Name)
				if data.Data.Category.Name != "" {
					v.Name = data.Data.Category.Name
				}
				v.LastUpdateTime = data.Data.Category.LastUpdateTime
				break
			}
		case "Delete":
			// log.Info("note.Categories", note.Categories)
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.CategoryId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			for i, v := range note.Categories {
				if v.Id != data.Options.CategoryId {
					continue
				}
				note.Categories = append(note.Categories[:i], note.Categories[i+1:]...)
				break
			}
		case "Sort":
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			if err = validation.ValidateStruct(
				data.Data.Note,
				validation.Parameter(&data.Data.Note.Categories, validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			note.Categories = data.Data.Note.Categories
		}
	case "Page":
		if data.Methods == "Add" || data.Methods == "Update" {
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.CategoryId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.PageId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			if err = validation.ValidateStruct(
				data.Data.Page,
				// validation.Parameter(&data.Data.Page.Id, validation.Type("string"), validation.Required()),
				// validation.Parameter(&data.Data.Page.Title, validation.Type("string"), validation.Required()),
				// validation.Parameter(&data.Data.Page.Content, validation.Type("string"), validation.Required()),
				// validation.Parameter(&data.Data.Page.CreateTime, validation.Type("int64"), validation.Required()),
				validation.Parameter(&data.Data.Page.LastUpdateTime, validation.Type("int64"), validation.Required()),
				// validation.Parameter(&data.Data.Page.Sort, validation.Type("int64"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
		}
		switch data.Methods {
		case "Add":
			// log.Info("note.Categories", note.Categories)
			for _, v := range note.Categories {
				if v.Id != data.Options.CategoryId {
					continue
				}

				v.Data = append(v.Data, data.Data.Page)
				break
			}
		case "Update":
			// log.Info("note.Categories", note.Categories)
			for _, v := range note.Categories {
				if v.Id != data.Options.CategoryId {
					continue
				}
				for _, sv := range v.Data {
					if sv.Id != data.Options.PageId {
						continue
					}
					if data.Data.Page.Title != "__nonec__xxx1,df,a" {
						sv.Title = data.Data.Page.Title
					}
					if data.Data.Page.Content != "__nonec__xxx1,df,a" {
						sv.Content = data.Data.Page.Content
					}
					// sv.Sort = data.Data.Page.Sort
					sv.LastUpdateTime = data.Data.Page.LastUpdateTime
					break
				}
				break
			}
		case "Delete":
			// log.Info("note.Categories", note.Categories)
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.CategoryId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.PageId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			for _, v := range note.Categories {
				if v.Id != data.Options.CategoryId {
					continue
				}
				for si, sv := range v.Data {
					if sv.Id != data.Options.PageId {
						continue
					}
					v.Data = append(v.Data[:si], v.Data[si+1:]...)
					break
				}
				break
			}
		case "Sort":
			if err = validation.ValidateStruct(
				data.Options,
				validation.Parameter(&data.Options.NoteId, validation.Type("string"), validation.Required()),
				validation.Parameter(&data.Options.CategoryId, validation.Type("string"), validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			if err = validation.ValidateStruct(
				data.Data.Category,
				validation.Parameter(&data.Data.Category.Data, validation.Required()),
			); err != nil {
				res.Error = err.Error()
				res.Code = 10002
				res.Call(c)
				return
			}
			for _, v := range note.Categories {
				if v.Id != data.Options.CategoryId {
					continue
				}
				v.Data = data.Data.Category.Data
				break
			}
		}
	}
	// log.Info("note", note)
	lastUpdateTime := int64(0)

	if data.Type != "Note" || (data.Type == "Note" && data.Methods == "Update") {

		lastUpdateTime = time.Now().Unix()
		note.LastUpdateTime = lastUpdateTime
		err = syncDbx.UpdateNote(note, note.Id, authorId)
		if err != nil {
			res.Error = err.Error()
			res.Code = 10011
			res.Call(c)
			return
		}
	}

	protoData := &protos.SyncToServer_Response{
		LastUpdateTime: lastUpdateTime,
	}

	res.Data = protos.Encode(protoData)

	res.Call(c)

	conf.G.Go(func() {
		// cc := nsocketio.ConnContext{
		// 	ServerContext: conf.SocketIO,
		// }
		sc := conf.SocketIO
		getConnContext := sc.GetConnContextByTag(conf.SocketRouterNamespace["Sync"], "Uid", nstrings.ToString(authorId))
		// log.Info("当前ID", c.ID())
		// log.Info("有哪些设备在线", getConnContext)
		for _, v := range getConnContext {
			if v.GetSessionCache("deviceId") == deviceId {
				continue
			}

			var res response.ResponseProtobufType
			res.Code = 200
			sdr := &protos.SyncData_Response{}

			copier.Copy(sdr, data)

			res.Data = protos.Encode(sdr)

			eventName := conf.SocketRouterEventNames["SyncData"]
			// log.Info("responseData", responseData)
			v.Emit(eventName, res.ResponseProtoEncode())
		}

	})

}

func (fc *SyncController) GetUrls(c *gin.Context) {
	// 1、请求体
	var res response.ResponseProtobufType
	res.Code = 200

	// 2、获取参数
	data := new(protos.GetUrls_Request)
	var err error
	if err = protos.DecodeBase64(c.GetString("data"), data); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	// 3、验证参数
	if err = validation.ValidateStruct(
		data,
		validation.Parameter(&data.NoteId, validation.Type("string"), validation.Required()),
	); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	authorId := c.MustGet("userInfo").(*sso.UserInfo).Uid

	// 4、获取URLS

	path, fileName := methods.GetPathAndFileName(data.NoteId, authorId)

	urls, err := conf.SAaSS.GetUrls(path, fileName)
	log.Info("urls", urls)
	if err != nil {
		res.Error = err.Error()
		res.Code = 10020
		res.Call(c)
		return
	}

	// 按照顺序排序
	res.Data = protos.Encode(&protos.GetUrls_Response{
		Urls: &protos.Urls{
			DomainUrl:     urls.DomainUrl,
			EncryptionUrl: urls.EncryptionUrl,
			Url:           urls.Url,
		},
	})
	res.Call(c)
}

func (fc *SyncController) GetFolderFiles(c *gin.Context) {
	// 1、请求体
	var res response.ResponseProtobufType
	res.Code = 200

	// 2、获取参数
	data := new(protos.GetFolderFiles_Request)
	var err error
	if err = protos.DecodeBase64(c.GetString("data"), data); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	authorId := c.MustGet("userInfo").(*sso.UserInfo).Uid

	// 4、获取URLS

	path, _ := methods.GetPathAndFileName("", authorId)

	files, err := conf.SAaSS.GetFolderFiles(path)
	// log.Info("files", files)
	if err != nil {
		res.Error = err.Error()
		res.Code = 10020
		res.Call(c)
		return
	}

	protoData := &protos.GetFolderFiles_Response{
		Total: files.Total,
	}

	list := []*protos.GetFolderFiles_Response_UrlsItem{}

	for _, v := range files.List {
		ut := new(protos.GetFolderFiles_Response_UrlsItem)
		ut.Urls = &protos.Urls{
			DomainUrl:     v.Urls.DomainUrl,
			EncryptionUrl: v.Urls.EncryptionUrl,
			Url:           v.Urls.Url,
		}
		ut.Id = strings.Replace(v.FileName, ".note", "", 1)
		note, err := syncDbx.GetNoteBySAaSS(ut.Id, authorId)
		if err != nil {
			continue
		}
		ut.LastUpdateTime = note.LastUpdateTime
		list = append(list, ut)
	}

	protoData.List = list

	res.Data = protos.Encode(protoData)

	res.Call(c)
}

func (fc *SyncController) GetNote(c *gin.Context) {
	// 1、请求体
	var res response.ResponseProtobufType
	res.Code = 200

	// 2、获取参数
	data := new(protos.GetNote_Request)
	var err error
	if err = protos.DecodeBase64(c.GetString("data"), data); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	// 3、验证参数
	if err = validation.ValidateStruct(
		data,
		validation.Parameter(&data.Id, validation.Type("string"), validation.Required()),
	); err != nil {
		res.Error = err.Error()
		res.Code = 10002
		res.Call(c)
		return
	}

	authorId := c.MustGet("userInfo").(*sso.UserInfo).Uid

	// 4、获取Note
	note, err := syncDbx.GetNoteBySAaSS(data.Id, authorId)
	// log.Info(note.Id, err, data.Id, authorId)
	if err != nil {
		res.Error = err.Error()
		res.Code = 10021
		res.Call(c)
		return
	}
	protoData := &protos.GetNote_Response{
		Note: note,
	}

	res.Data = protos.Encode(protoData)

	res.Call(c)
}

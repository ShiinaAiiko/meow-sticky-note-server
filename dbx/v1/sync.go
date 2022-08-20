package dbxV1

import (
	"bytes"
	"encoding/json"
	"time"

	conf "github.com/ShiinaAiiko/meow-sticky-note/server/config"
	"github.com/ShiinaAiiko/meow-sticky-note/server/protos"
	"github.com/ShiinaAiiko/meow-sticky-note/server/services/methods"
	"github.com/cherrai/nyanyago-utils/nfile"
	"github.com/cherrai/nyanyago-utils/nlog"
	"github.com/cherrai/nyanyago-utils/nstrings"
	"github.com/cherrai/nyanyago-utils/nutils"
	"github.com/cherrai/nyanyago-utils/saass"
)

type SyncDbx struct {
}

var (
	log = nlog.New()
)

func (s *SyncDbx) GetNoteBySAaSS(noteId string, authorId int64) (*protos.NoteItem, error) {
	note := new(protos.NoteItem)
	key := conf.Redisdb.GetKey("GetNote")
	path, fileName := methods.GetPathAndFileName(noteId, authorId)

	err := conf.Redisdb.GetStruct(key.GetKey(noteId), note)
	// conf.Redisdb.Delete(key.GetKey(noteId))

	// log.Info("err", err, err != nil, note)
	// log.Info(noteId)
	// log.Info("GetNoteBySAaSS", path, fileName)
	if err != nil {
		reader, err := conf.SAaSS.GetFile(path, fileName)
		// log.Info("GetFile", reader, err, path, fileName)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(reader.Bytes(), note); err != nil {
			return nil, err
		}
		if err = conf.Redisdb.SetStruct(key.GetKey(noteId), note, key.GetExpiration()); err != nil {
			log.Info(err)
			return nil, err
		}
	}

	return note, nil
}
func (s *SyncDbx) AddNote(note *protos.NoteItem, noteId string, authorId int64) error {
	key := conf.Redisdb.GetKey("GetNote")
	if err := conf.Redisdb.SetStruct(key.GetKey(noteId), note, key.GetExpiration()); err != nil {
		return err
	}

	s.SyncSAaSS(noteId, authorId)
	return nil
}
func (s *SyncDbx) DeleteNote(noteId string, authorId int64) error {
	key := conf.Redisdb.GetKey("GetNote")
	if err := conf.Redisdb.Delete(key.GetKey(noteId)); err != nil {
		return err
	}
	path, fileName := methods.GetPathAndFileName(noteId, authorId)

	log.Info("path, fileName", path, fileName)
	return conf.SAaSS.DeleteFile(path, fileName, time.Now().Add(30*24*3600*time.Second))
}

func (s *SyncDbx) UpdateNote(note *protos.NoteItem, noteId string, authorId int64) error {
	key := conf.Redisdb.GetKey("GetNote")
	// path, fileName := methods.GetPathAndFileName(noteId, authorId)
	// log.Info("getRedis1", note)
	if err := conf.Redisdb.SetStruct(key.GetKey(noteId), note, key.GetExpiration()); err != nil {
		return err
	}

	// s.SyncSAaSS(noteId, authorId)
	nutils.NewDebounce(30*time.Second, "UpdateNote_"+noteId+"_"+nstrings.ToString(authorId)).Add(func() {

		s.SyncSAaSS(noteId, authorId)
	})
	return nil
}

func (s *SyncDbx) SyncSAaSS(noteId string, authorId int64) (*saass.Urls, error) {
	note := new(protos.NoteItem)
	key := conf.Redisdb.GetKey("GetNote")
	// log.Info("getRedis2", note)
	if err := conf.Redisdb.GetStruct(key.GetKey(noteId), note); err != nil {
		log.Info("内容已经被删除了")
		return nil, err
	}
	// conf.Redisdb.Delete(key.GetKey(noteId))
	// log.Info("data", data.Options)
	// log.Info("data", data.Data.Note)

	// 4、获取数据
	// log.Info("authorId", authorId)

	// "Image", "Video", "Audio", "Text", "File"
	path, fileName := methods.GetPathAndFileName(noteId, authorId)

	// log.Info("path", path)
	// log.Info("fileName", fileName)
	// log.Info("./temp" + path + fileName)

	// noteFile, err := ioutil.ReadFile("./temp" + path + fileName)
	// if err != nil {
	// 	// fmt.Println(err)
	// 	log.Error(err)
	// }
	// if noteFile==nil{

	// }

	noteByte, err := json.Marshal(note)
	if err != nil {
		// log.Info("err", 3, err)
		return nil, err
	}

	reader := bytes.NewReader(noteByte)
	size := reader.Size()

	hash := nfile.MuseGetHashByBytes(noteByte)
	// log.Info("hash", hash)
	// log.Info("size", size)

	// 如果该路径已经存在文件怎么办？
	// 1、覆盖更新操作
	// 2、不覆盖，返回原来文件下载地址

	ut, err := conf.SAaSS.CreateChuunkUploadToken(&saass.CreateUploadTokenOptions{
		FileInfo: &saass.FileInfo{
			Name:         fileName,
			Size:         size,
			Type:         ".note",
			Suffix:       ".note",
			LastModified: time.Now().UnixMilli(),
			Hash:         hash,
		},
		Path:           path,
		FileName:       fileName,
		ChunkSize:      256 * 1024,
		VisitCount:     -1,
		ExpirationTime: -1,
		// Type:           "File",
		FileConflict: "Replace",
		OnProgress: func(progress saass.Progress) {
			// log.Info("progress", progress)
		},
		OnSuccess: func(urls saass.Urls) {
			// log.Info("urls", urls)
		},
		OnError: func(err error) {
			// log.Info("err", err)
		},
	})
	if err != nil {
		return nil, err
	}
	// log.Info("ut, err", ut, err, ut.Token != "")

	// log.Info(ut.Urls.DomainUrl + ut.Urls.EncryptionUrl)
	// log.Info(ut.Urls.DomainUrl + ut.Urls.Url)
	// if ut.Token == "" {
	// 	return &ut.Urls, nil
	// }
	urls, err := ut.ChunkUpload(noteByte)
	// log.Info(urls, err)
	if err != nil {
		return nil, err
	}
	return &urls, nil
	// return nil, nil
}

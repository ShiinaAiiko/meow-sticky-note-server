package methods

import (
	"github.com/cherrai/nyanyago-utils/cipher"
	"github.com/cherrai/nyanyago-utils/nstrings"
)

func GetPathAndFileName(noteId string, authorId int64) (path, fileName string) {
	path = "/" + cipher.MD5(nstrings.ToString(authorId)) + "/notes/"
	fileName = noteId + ".note"
	return
}

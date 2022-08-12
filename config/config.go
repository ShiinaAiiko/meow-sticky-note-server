package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ShiinaAiiko/meow-sticky-note/server/services/typings"
	"github.com/cherrai/nyanyago-utils/goroutinepanic"
	"github.com/cherrai/nyanyago-utils/nsocketio"
	"github.com/cherrai/nyanyago-utils/saass"
	sso "github.com/cherrai/saki-sso-go"
)

var (
	Config        *typings.Config
	SSO           *sso.SakiSSO
	SAaSS         *saass.SAaSS
	SocketIO      *nsocketio.NSocketIO
	G             = goroutinepanic.G
	FileTokenSign = "saass_2022_6_4"
	// 文件到期后根据时间进行删除 未做
	// []string{"Image", "Video", "Audio", "Text", "File"}
	FileExpirationRemovalDeadline = 60 * 3600 * 24 * time.Second
	// 临时文件删除期限
	TempFileRemovalDeadline = 60 * 3600 * 24 * time.Second
)

func GetConfig(configPath string) {
	jsonFile, _ := os.Open(configPath)

	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)

	conf := new(typings.Config)
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	Config = conf
}

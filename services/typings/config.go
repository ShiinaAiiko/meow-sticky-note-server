package typings

type Config struct {
	Server Server
	SSO    Sso
	Saass  SAaSS
	Redis  Redis
	// StaticPathDomain string
	Mongodb Mongodb
}

type Server struct {
	Port int
	Cors struct {
		AllowOrigins []string
	}
	// mode: release debug
	Mode string
}
type Sso struct {
	AppId  string
	AppKey string
	Host   string
}
type SAaSS struct {
	AppId      string
	AppKey     string
	BaseUrl    string
	ApiVersion string
}
type Redis struct {
	Addr     string
	Password string
	DB       int
}
type Mongodb struct {
	Currentdb struct {
		Name string
		Uri  string
	}
}

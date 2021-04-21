package common

// URL
const (
	HomePage      = "127.0.0.1:11111"
	Callback      = "http://" + HomePage + "/oauth/redirect"
	OAuthPage     = "http://127.0.0.1:22222"
	OAuthWelcome  = "/oauth/welcome"
	OAuthCreate   = "/oauth/create"
	OAuthToken    = "/oauth/token"
	OAuthResource = "/oauth/resource"
	OAuthRefresh  = "/oauth/refresh"
)

// FILE
const (
	LogFile = "output/log/gin.log"
)

var (
	AppID     string
	AppSecret string
)

// CONST
const (
	StringAll         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	StringUpper       = 64
	TokenType         = "bearer"
	AuthorizationCode = "authorization_code"
	RefreshToken      = "refresh_token"
	Resource1         = "资源1："
	Resource2         = "资源2："
	Resource3         = "资源3："
)

package model

import (
	"crypto/rsa"
	"database/sql"
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

var DAO *sql.DB
var XDAO *xorm.Engine
var GinEngine *gin.Engine
var PrivateKey *rsa.PrivateKey
var JwtKey []byte
var ReCaptchaTokenKey string
var ReCaptchaURL string
var PublicFSPath string
var AuthedFSPath string

type ResMessage struct {
	OkFlag     bool
	Message    string
	Data       any
	ErrDetails error
}

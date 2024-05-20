package User

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserVO struct {
	ID             int64  ``
	Username       string `xorm:"'user_name'"`
	Password       string
	Email          string
	ImgUrl         string
	CreateTime     time.Time `xorm:"created 'create_time'"`
	Ban            bool
	IsEmailChecked bool
}

type AccountInfo struct {
	Username       string
	ImgUrl         string
	Email          string
	IsEmailChecked bool
}

func (p *UserVO) TableName() string {
	return "sm_user"
}

type User struct {
	Username string
	Password string
	Email    string
	Token    string
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

type Session struct {
	Username   string
	UserImgUrl string
	SessionID  string
	TS         time.Time
}

type CheckMode string

const (
	Email    CheckMode = "Email"
	UserName CheckMode = "UserName"
)

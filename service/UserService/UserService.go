package UserService

import (
	"awesomeProject/model"
	"awesomeProject/model/User"
	"awesomeProject/util"
	"github.com/google/uuid"
	"time"
)

var OnlineUser = make(map[string]User.Session, 10)
var OnlineUserSessionMap = make(map[string]string, 10)

// LoginUser /**登录
func LoginUser(user User.User) (bool, string, string, error) {
	//检查Token
	captchaToken, err := util.CheckReCAPTCHAToken(user.Token)
	if err != nil {
		return false, "ReCAPTCHAToken校验错误", "", err
	}
	if !captchaToken {
		return false, "ReCAPTCHAToken校验错误", "", err
	}
	isRight, msg, err := checkPassword(user)
	if !isRight || err != nil {
		return false, msg, "", err
	}

	token, err := doUserLogin(user)
	if err != nil {
		return false, "生成Token失败", "", err
	}
	return true, "登录成功", token, nil

}

func checkPassword(user User.User) (bool, string, error) {
	username := user.Username
	password := user.Password
	password, err := util.RsaDecrypt(password)
	if err != nil {
		return false, "登录失败，密码解析失败", err
	}
	row := model.DAO.QueryRow("select COUNT(1) from sm_user where user_name = ? and password = ?", username, password)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, "登录失败，查询失败", err
	}
	if count == 1 {
		return true, "校验成功", nil
	} else {
		return false, "登录失败，密码错误或用户不存在", nil
	}
}

func doUserLogin(user User.User) (string, error) {
	sessionID, err := generateSessionID(user.Username)

	return sessionID, err
}

func generateSessionID(userName string) (string, error) {
	//检查是否用户已经登录
	sessionID, ok := OnlineUserSessionMap[userName]
	//如果登录要先使会话失效
	if ok {
		delete(OnlineUser, sessionID)
	}
	//查询用户头像地址
	userImgUrl := ""
	err := model.DAO.QueryRow("select img_url from sm_user where user_name = ?", userName).Scan(&userImgUrl)
	if err != nil {
		return "", err
	}
	sessionID = uuid.NewString()
	now := time.Now()
	OnlineUser[sessionID] = User.Session{
		Username:   userName,
		SessionID:  sessionID,
		UserImgUrl: userImgUrl,
		TS:         now,
	}
	OnlineUserSessionMap[userName] = sessionID
	return sessionID, nil
}

func UpdateSessionTS(sessionID string) {
	sessionData, onLine := OnlineUser[sessionID]
	if onLine {
		sessionData.TS = time.Now()
	}
}

func LogoutUser(sessionID string) {
	//如果登录要先使会话失效
	delete(OnlineUser, sessionID)
}

func RegisterUser(user User.User) (bool, string, string, error) {
	//检查Token
	captchaToken, err := util.CheckReCAPTCHAToken(user.Token)
	if err != nil || !captchaToken {
		return false, "ReCAPTCHAToken校验错误", "", err
	}

	//查询用户名
	row := model.DAO.QueryRow("select count(1) from sm_user where user_name = ?", user.Username)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, "注册失败，查询用户名失败", "", err
	}
	if count > 0 {
		return false, "用户名重复", "", err
	}

	//邮箱唯一性
	row = model.DAO.QueryRow("select count(1) from sm_user where email = ?", user.Email)
	err = row.Scan(&count)
	if err != nil {
		return false, "注册失败，查询用户名失败", "", err
	}
	if count > 0 {
		return false, "邮箱重复", "", err
	}

	//还原密码
	password := user.Password
	password, err = util.RsaDecrypt(password)
	if err != nil {
		return false, "注册失败，密码解析失败", "", err
	}

	userVO := User.UserVO{Username: user.Username, Password: password, Email: user.Email}
	_, err = model.XDAO.Insert(userVO)
	if err != nil {
		return false, "注册失败, 插入数据失败", "", err
	}

	token, err := doUserLogin(user)
	if err != nil {
		return false, "生成Token失败", "", err
	}

	return true, "注册成功", token, nil
}

func UpdateUserImg(sessionID string, pubFileULR string) error {
	session := OnlineUser[sessionID]
	username := session.Username
	_, err := model.XDAO.Exec("update sm_user set img_url = ? where user_name = ?", pubFileULR, username)
	if err != nil {
		return err
	}
	return nil
}

func GetAccountInfo(sessionID string) (User.AccountInfo, error) {
	session := OnlineUser[sessionID]
	username := session.Username
	var user []User.UserVO
	err := model.XDAO.Where("user_name = ?", username).Find(&user)
	var account User.AccountInfo
	if err != nil {
		return account, err
	}
	account = User.AccountInfo{
		Username:       user[0].Username,
		Email:          user[0].Email,
		ImgUrl:         user[0].ImgUrl,
		IsEmailChecked: user[0].IsEmailChecked,
	}

	return account, nil
}

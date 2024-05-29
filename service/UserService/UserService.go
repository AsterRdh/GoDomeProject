package UserService

import (
	"awesomeProject/model"
	"awesomeProject/model/Email"
	"awesomeProject/model/User"
	"awesomeProject/service/MailService"
	"awesomeProject/util"
	"encoding/json"
	"github.com/google/uuid"
	"net/url"
	"strings"
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

	userVO := User.UserVO{
		Username: user.Username,
		Password: password,
		Email:    user.Email,
		Salt:     uuid.NewString(),
	}
	_, err = model.XDAO.Insert(userVO)
	if err != nil {
		return false, "注册失败, 插入数据失败", "", err
	}

	token, err := doUserLogin(user)
	if err != nil {
		return false, "生成Token失败", "", err
	}

	//发送邮件
	//生成确认地址
	err = SendCheckEmail(user.Email, userVO.Salt)
	if err != nil {

	}

	return true, "注册成功", token, nil
}

func SendCheckEmail(userEmail string, salt string) error {

	var keyObject = Email.EMailKey{
		Email: userEmail,
		Ts:    time.Now(),
	}
	marshal, err := json.Marshal(keyObject)
	if err != nil {
		return err
	}
	key := string(marshal)

	var AESKey = util.GetAesKey(salt)
	key = util.AesEncryptByECB(key, AESKey)

	var url = model.BaseURL + "/checkemail?email=" + url.QueryEscape(userEmail) + "&key=" + url.QueryEscape(key)
	sprintf := strings.Replace(Email.Template.ResCheckEmail, "{t_url}", url, 2)
	err = MailService.SendMail(model.EMail.Accounts.Account, userEmail, sprintf, "确认邮件地址")
	return err
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

func CheckEmail(email string, key string) (bool, string, error, string) {
	//var salts []string
	//err := model.XDAO.Where("email = ?", email).Table("sm_user").Cols("salt").Find(&salts)
	var users []User.UserVO
	err := model.XDAO.Where("email = ?", email).Find(&users)
	if err != nil || len(users) != 1 {
		return false, "未知的邮件地址", err, "E000001"
	}
	var user = users[0]
	var salt = user.Salt
	aesKey := util.GetAesKey(salt)
	ecb, errMsg, err := util.AesDecryptByECB(key, aesKey)
	if err != nil {
		return false, errMsg, err, "E000001"
	}

	var keyObject Email.EMailKey
	err = json.Unmarshal([]byte(ecb), &keyObject)
	if err != nil {
		return false, "未知的邮件地址", err, "E000001"
	}

	if email != keyObject.Email {
		return false, "密文不合法", nil, "E000001"
	}

	sub := time.Now().Sub(keyObject.Ts)
	if sub.Minutes() > 15 {
		return false, "密钥过期", nil, "E000002"
	}
	//设置已认证
	user.IsEmailChecked = true
	_, err = model.XDAO.ID(user.ID).AllCols().Update(&user)
	if err != nil {
		return false, "未知错误", nil, "E000001"
	}

	//登录
	sessionID, err := doUserLogin(User.User{Username: user.Username})
	if err != nil {
		return false, "", err, ""
	}

	return true, sessionID, nil, "ok"
}

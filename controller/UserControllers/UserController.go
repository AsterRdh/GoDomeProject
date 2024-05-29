package UserControllers

import (
	"awesomeProject/model"
	"awesomeProject/model/User"
	"awesomeProject/service/FileService"
	"awesomeProject/service/UserService"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginAction /** 登录
func LoginAction(c *gin.Context) {
	var user User.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, "无效的 JSON 格式")
		return
	}
	isOK, msg, token, err := UserService.LoginUser(user)
	resMsg := model.ResMessage{
		OkFlag:     isOK,
		Message:    msg,
		Data:       nil,
		ErrDetails: err,
	}
	if isOK {
		c.SetCookie("session_id", token, 3600, "/", "", false, false)
	}

	c.JSON(200, resMsg)
}

// GetUserBaseDataAction /** 获取用户基本信息
func GetUserBaseDataAction(c *gin.Context) {
	cookie, err := c.Cookie("session_id")
	if err != nil {
		c.SetCookie("session_id", "", 3600, "/", "", false, false)
		c.JSON(http.StatusUnprocessableEntity, "cookie error")

	} else {
		session := UserService.OnlineUser[cookie]

		resMsg := &model.ResMessage{
			OkFlag:     true,
			Message:    "ok",
			Data:       session,
			ErrDetails: err,
		}
		c.JSON(200, resMsg)
	}
}

// LogoutAction /** 登出
func LogoutAction(c *gin.Context) {
	cookie, err := c.Cookie("session_id")
	if err != nil {
		UserService.LogoutUser(cookie)
		c.SetCookie("session_id", "", 3600, "/", "", false, false)
	}
	resMsg := &model.ResMessage{
		OkFlag:     true,
		Message:    "ok",
		Data:       "",
		ErrDetails: err,
	}
	c.JSON(200, resMsg)
}

// RegisterAction /** 注册
func RegisterAction(c *gin.Context) {
	var user User.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, "无效的 JSON 格式")
		return
	}

	isOK, msg, token, err := UserService.RegisterUser(user)
	resMsg := &model.ResMessage{
		OkFlag:     isOK,
		Message:    msg,
		Data:       nil,
		ErrDetails: err,
	}
	if isOK {
		c.SetCookie("session_id", token, 3600, "/", "", false, false)
	}

	c.JSON(200, resMsg)
}

// UpdateUserImgAction /** 更新用户头像
func UpdateUserImgAction(c *gin.Context) {
	cookie, err := c.Cookie("session_id")
	if err != nil {
		c.SetCookie("session_id", "", 3600, "/", "", false, false)
		c.JSON(http.StatusUnprocessableEntity, "cookie error")
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传图片出错")
	}
	pubFileULR, err := FileService.UploadPubFile(c, file)
	if err != nil {
		c.String(500, err.Error())
	}
	err = UserService.UpdateUserImg(cookie, pubFileULR)
	if err != nil {
		c.String(500, err.Error())
	}
	resMsg := &model.ResMessage{
		OkFlag:     true,
		Message:    "ok",
		Data:       pubFileULR,
		ErrDetails: err,
	}
	c.JSON(200, resMsg)
}

// GetAccountInfoAction /** 加载账户信息
func GetAccountInfoAction(c *gin.Context) {
	cookie, err := c.Cookie("session_id")
	if err != nil {
		c.SetCookie("session_id", "", 3600, "/", "", false, false)
		c.JSON(http.StatusUnprocessableEntity, "cookie error")
	}
	info, err := UserService.GetAccountInfo(cookie)
	var resMsg *model.ResMessage
	if err != nil {
		resMsg = &model.ResMessage{
			OkFlag:     false,
			Message:    "未找到用户信息",
			Data:       nil,
			ErrDetails: err,
		}
		c.JSON(200, resMsg)
	}

	resMsg = &model.ResMessage{
		OkFlag:     true,
		Message:    "查询成功",
		Data:       info,
		ErrDetails: err,
	}
	c.JSON(200, resMsg)
}

// CheckEmail /**校验邮件地址
func CheckEmail(c *gin.Context) {

	email := c.Query("email")
	key := c.Query("key")

	checkEmail, sessionID, err, errorCode := UserService.CheckEmail(email, key)
	if err != nil || !checkEmail {
		c.Redirect(http.StatusTemporaryRedirect, model.BaseWebURL+"failCheckEmail?code="+errorCode)
	}
	c.SetCookie("session_id", sessionID, 3600, "/", "", false, false)
	c.Redirect(http.StatusTemporaryRedirect, model.BaseWebURL)

}

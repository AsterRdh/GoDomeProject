package main

import (
	"awesomeProject/configs"
	"awesomeProject/controller/FileControllers"
	"awesomeProject/controller/Filters"
	"awesomeProject/controller/UserControllers"
	"awesomeProject/model"
	"awesomeProject/model/Email"
	"awesomeProject/service/MailService"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "net/http"
	"os"
	"strings"
)

// 初始化配置文件
func initConfigs() (err error) {
	file, err := loadPrivateKeyFile("./config/private.pem")
	if err != nil {
		return err
	}
	model.PrivateKey = file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	model.JwtKey = []byte(viper.GetString("jwtKey"))
	model.ReCaptchaTokenKey = viper.GetString("ReCaptchaTokenKey")
	model.ReCaptchaURL = viper.GetString("ReCaptchaURL")
	model.AuthedFSPath = viper.GetString("AuthedFSPath")
	model.PublicFSPath = viper.GetString("PublicFSPath")
	model.AESKey = viper.GetString("AESKey")
	model.BaseURL = viper.GetString("BaseURL")
	model.BaseWebURL = viper.GetString("BaseWebURL")

	model.EMail = model.EMailConfig{
		Server: viper.GetString("mail.server"),
		Portal: viper.GetInt("mail.portal"),
		Accounts: model.EMailAccounts{
			Account: model.EMailAccountConfig{
				Account:  viper.GetString("mail.account.account"),
				Password: viper.GetString("mail.account.password"),
			},
		},
	}
	return err
}

func loadEmailTemplate() error {
	var EMailTemplatePath = "./template"
	CheckEmailBuffer, err := os.ReadFile(EMailTemplatePath + "/CheckEmail.html")
	if err != nil {
		return err
	}
	Email.Template.ResCheckEmail = string(CheckEmailBuffer)
	return nil
}

// 加载RSA密钥
func loadPrivateKeyFile(keyfile string) (*rsa.PrivateKey, error) {
	keyBuffer, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(keyBuffer))
	if block == nil {
		return nil, errors.New("private key error!")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}

// 初始化数据库链接
func initDB() (err error) {
	return configs.InitBD()
}

// 初始化服务器
func initServer() {
	r := gin.Default()
	model.GinEngine = r
}

// 初始化接口
func initController() {
	//用户登录接口
	model.GinEngine.POST("/login", UserControllers.LoginAction)
	//注册接口
	model.GinEngine.POST("/register", UserControllers.RegisterAction)
	//公共资源接口
	model.GinEngine.Static("/fs", model.PublicFSPath)
	//已授权接口
	adminGroup := model.GinEngine.Group("/authed")
	adminGroup.Use(Filters.AuthFilter()) //设置过滤器
	//获取用户基本信息接口
	adminGroup.GET("/getUserBaseData", UserControllers.GetUserBaseDataAction)
	//用户登出接口
	adminGroup.GET("/logout", UserControllers.LogoutAction)

	//授权资源接口
	adminGroup.Static("/fs", model.AuthedFSPath)
	//上传文件接口
	adminGroup.POST("/fs/upload/pub", FileControllers.UploadPub)
	adminGroup.POST("/fs/upload/auth", FileControllers.UploadAuthed)
	adminGroup.POST("/updateUserImg", UserControllers.UpdateUserImgAction)
	//加载用户信息接口
	adminGroup.POST("/getUserAccountInfo", UserControllers.GetAccountInfoAction)

	//测试接口
	adminGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Message": "pong",
			"OkFlag":  true,
		})
	})

	model.GinEngine.GET("/ping", func(c *gin.Context) {

		sprintf := strings.Replace(Email.Template.ResCheckEmail, "{t_url}", "http://baidu.com", 2)
		err := MailService.SendMail(model.EMail.Accounts.Account, "asterrdh@live.cn", sprintf, "测试")
		if err != nil {
			c.JSON(500, err)
		}
		c.JSON(200, gin.H{
			"Message": "pong",
			"OkFlag":  true,
		})
	})

	model.GinEngine.GET("/checkemail", UserControllers.CheckEmail)

}

func main() {

	var err error
	err = initConfigs()
	if err != nil {
		fmt.Printf("open server faild,err:%v\n", err)
		return
	}

	err = loadEmailTemplate()
	if err != nil {
		fmt.Printf("Load Email template faild,err:%v\n", err)
		return
	}

	err = initDB()
	if err != nil {
		fmt.Printf("open server faild,err:%v\n", err)
		return
	}

	initServer()
	initController()

	err = model.GinEngine.Run()
	if err != nil {
		fmt.Printf("open server faild,err:%v\n", err)
		return
	}
	fmt.Printf("server started")
}

package Filters

import (
	"awesomeProject/service/UserService"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		}
		sessionData, ok := UserService.OnlineUser[sessionID]
		if !ok {
			c.SetCookie("session_id", "", -1, "/", "", false, false)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		}
		ts := sessionData.TS
		sessionIDDuration := nowTime.Sub(ts)
		if sessionIDDuration > time.Duration(10)*time.Minute {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			UserService.LogoutUser(sessionID)
		}
		//更新时间戳
		UserService.UpdateSessionTS(sessionID)
		c.Next()
	}
}

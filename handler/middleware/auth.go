package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"normal_web/util"
	"normal_web/util/logger"
)

const (
	UIN_IN_TOKEN = "uid"
)

var (
	appConf = util.CreateConfig("app")
)

func GetUidFromJWT(token string) int {
	_, payload, err := util.VerifyJWT(token, appConf.GetString("jwt"))
	if err != nil {
		return 0
	}
	for k, v := range payload.UserDefined {
		if k == UIN_IN_TOKEN {
			return int(v.(float64))
		}
	}

	return 0
}

func GetLoginUid(cxt *gin.Context) int {
	token := cxt.Request.Header.Get("auth_token")

	//方案二: 依靠浏览器自动回传的cookie，提取出refresh_token，再由服务端拿refresh_token查redis得到auth_token。增加了访问redis的频次，auth_token不需要传给前端；不同浏览器窗口共享同一个登录的uid。方案二实际上就是基于cookie的认证。
	// var token string
	//http协议里没有cookie这个概念，cookie本质上是header里的一对KV
	// for _, cookie := range strings.Split(ctx.Request.Header.Get("cookie"), ";") {
	// 	arr := strings.Split(cookie, "=")
	// 	key := strings.TrimSpace(arr[0])
	// 	value := strings.TrimSpace(arr[1])
	// 	if key == "refresh_token" {
	// 		token = database.GetToken(value)
	// 	}
	// }
	//或者直接使用封装好的Request.Cookies()
	// for _, cookie := range ctx.Request.Cookies() {
	// 	if cookie.Name == "refresh_token" {
	// 		fmt.Println(cookie.Value)
	// 		token = database.GetToken(cookie.Value)
	// 	}
	// }

	logger.Debug(fmt.Sprintf("get token from header %s", token))
	return GetUidFromJWT(token)
}

func Auth() gin.HandlerFunc {
	return func(cxt *gin.Context) {
		uid := GetLoginUid(cxt)
		if uid > 0 {
			cxt.Set("uid", uid)
			cxt.Next()
		} else {
			cxt.JSON(http.StatusForbidden, gin.H{"err_msg": "auth failed"})
			cxt.Abort()
		}
	}
}

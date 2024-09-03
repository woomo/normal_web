package main

import (
	"github.com/gin-gonic/gin"
	"normal_web/util"
	"normal_web/util/logger"
)

func Init() {
	logger.Init("log")
}
func main() {
	Init()
	logger.Error("test")
	route := gin.Default()
	appConf := util.CreateConfig("app")
	route.Run(appConf.GetString("host") + ":" + appConf.GetString("port"))

}

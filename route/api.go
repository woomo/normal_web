package route

import "github.com/gin-gonic/gin"

func RegisterApi(r *gin.Engine) {
	api := r.Group("/api")
	api.Use()
	{
		loginG := api.Group("/login")
		loginG.Use()
		{
			loginG.GET("", func(ctx *gin.Context) {
				ctx.HTML(200, "login.html", nil)
			})
			//loginG.POST("/submit", handler.Login)
		}
		//api.POST("/token", handler.GetAuthToken)

		blogG := api.Group("/blog")
		blogG.Use()
		{
			//blogG.GET("/blog/belong", handler.BlogBelong)
			//restful风格，参数放在url路径里
			//blogG.GET("/blog/list/:uid", handler.BlogList)
			//blogG.GET("/blog/:bid", handler.BlogDetail)
			//blogG.POST("/blog/update", middleware.Auth(), handler.BlogUpdate)
		}
	}

}

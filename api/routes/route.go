package routes

import (
	"truonghoang/go-scam/api/handles"
	// "truonghoang/go-scam/api/middleware"

	"github.com/gin-gonic/gin"
)

func RouteAccount(route *gin.RouterGroup) {
	accountGroup := route.Group("/account")
	{
		accountGroup.POST("/register", func(ctx *gin.Context) {
			handles.Register(ctx)
		})
		accountGroup.POST(("/login"), func(ctx *gin.Context) {
			handles.Login(ctx)
		})

	}
}

func RouteUserScam(route *gin.RouterGroup) {
	routeUserScam := route.Group("/report")
	// .Use(middleware.MiddleWare())
	{
		routeUserScam.GET("/", func(ctx *gin.Context) {
            handles.ListUserScam(ctx)
		})

		routeUserScam.DELETE("/", func(ctx *gin.Context) {
            
		})

		routeUserScam.GET("/:id", func(ctx *gin.Context) {
            
		})
	}
}

func RouteApi(route *gin.Engine) {
	apiGroup := route.Group("/api")
	{
		RouteAccount(apiGroup)
		RouteUserScam(apiGroup)
	}
}

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

func RouteUser(route *gin.RouterGroup) {
	accountGroup := route.Group("/user")
	{	// select user
		accountGroup.GET("/select", func(ctx *gin.Context) {
			handles.SelectUser(ctx)
		})
		//search user by phone
		accountGroup.GET(("/search"), func(ctx *gin.Context) {
			handles.GetDetailUser(ctx, true)
		})
		//detail user
		accountGroup.GET("/:id", func(ctx *gin.Context) {
			handles.GetDetailUser(ctx, false)
		})
		//list user pagination
		accountGroup.GET("", func(ctx *gin.Context) {
			handles.ListUser(ctx)
		})

	}
}
func RouteUserScam(route *gin.RouterGroup) {
	routeUserScam := route.Group("/report")
	// .Use(middleware.MiddleWare())
	{
		routeUserScam.GET("", func(ctx *gin.Context) {
			handles.ListReport(ctx)
		})

		routeUserScam.POST("", func(ctx *gin.Context) {
			handles.AddReport(ctx)
		})

		routeUserScam.GET("/:id", func(ctx *gin.Context) {
			handles.DetailReport(ctx)
		})
	}
}

func RouteApi(route *gin.Engine) {
	apiGroup := route.Group("/api")
	{
		RouteAccount(apiGroup)
		RouteUserScam(apiGroup)
		RouteUser(apiGroup)
	}
}

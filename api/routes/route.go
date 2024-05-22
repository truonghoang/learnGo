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
	routeReport := route.Group("/report")
	// .Use(middleware.MiddleWare())
	{
		routeReport.GET("", func(ctx *gin.Context) {
			handles.ListReport(ctx)
		})

		routeReport.POST("", func(ctx *gin.Context) {
			handles.AddReport(ctx)
		})
		routeReport.GET("/search",func(ctx *gin.Context) {
			handles.SearchPhoneReport(ctx)
		})
		routeReport.GET("/filter",func(ctx *gin.Context){
			handles.FilterReportByReason(ctx)
		})
		routeReport.GET("/:id", func(ctx *gin.Context) {
			handles.DetailReport(ctx)
		})
		routeReport.DELETE("/:id",func(ctx *gin.Context){
			handles.DeleteReport(ctx)
		})
		routeReport.GET("/detail/account",func(ctx *gin.Context) {
			handles.GetListAccountByDetail(ctx)
		})
		routeReport.GET("/detail/link",func(ctx *gin.Context) {
			handles.GetListLinkByDetail(ctx)
		})
		routeReport.GET("/detail/list",func(ctx *gin.Context) {
			handles.GetListReportByPeerId(ctx)
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

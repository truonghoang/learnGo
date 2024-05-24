package routes

import (
	
	"truonghoang/go-scam/api/handles"
	"truonghoang/go-scam/api/middleware"

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
	accountGroup := route.Group("/user").Use(middleware.MiddleWare())
	{	// select user
		accountGroup.GET("/select", func(ctx *gin.Context) {
			handles.SelectUser(ctx)
		})
		
		//search user by phone
		accountGroup.GET(("/search"), func(ctx *gin.Context) {
			handles.GetDetailUser(ctx, true)
		})
		accountGroup.GET("/banned",func(ctx *gin.Context) {
			handles.ListUserBan(ctx)
		})
		accountGroup.POST("/banned",func(ctx *gin.Context) {
			handles.BanAndUnBanUser(ctx)
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
	routeReport := route.Group("/report").Use(middleware.MiddleWare())
	{
		routeReport.GET("", func(ctx *gin.Context) {
			handles.ListReport(ctx)
		})
		routeReport.GET("/search",func(ctx *gin.Context) {
			handles.SearchReport(ctx)
		})
		routeReport.GET("/filter",func(ctx *gin.Context){
			handles.FilterReportByReason(ctx)
		})
		routeReport.POST("/process",func(ctx *gin.Context) {
			handles.HandleProcessReadReport(ctx)
		})

		routeReport.GET("/history/:id",func(ctx *gin.Context) {
			handles.HandleHistoryReport(ctx)
		})
		routeReport.GET("/:id", func(ctx *gin.Context) {
			handles.DetailReport(ctx)
		})
		
		routeReport.GET("/detail/account",func(ctx *gin.Context) {
			handles.GetListAccountByDetail(ctx)
		})
		routeReport.GET("/detail/alias",func(ctx *gin.Context) {
			handles.GetListLinkByDetail(ctx)
		})
		
		routeReport.GET("/reporter/:id",func(ctx *gin.Context) {
			handles.HandleListReportByReporter(ctx)
		})
		routeReport.GET("/reporter/filter/:id",func(ctx *gin.Context) {
			handles.FilterReportByReporter(ctx)
		})
		routeReport.GET("/reported-user/list",func(ctx *gin.Context) {
			handles.GetListReportByPeerId(ctx)
		})
		routeReport.GET("/reported-user/:id",func(ctx *gin.Context) {
			
			handles.HandleDetailReportByPeerId(ctx)
		})
		routeReport.DELETE("/:id/:process",func(ctx *gin.Context){
			handles.HandleAccessOrDenyReport(ctx)
		})
		routeReport.GET("/reported-user/list/filter",func(ctx *gin.Context) {
			handles.FilterReportBannedByReason(ctx)
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

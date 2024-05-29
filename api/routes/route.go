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
		// xác nhận đọc tin
		routeReport.POST("/process",func(ctx *gin.Context) {
			handles.HandleProcessReadReport(ctx)
		})
		//xem lịch sử người ban
		routeReport.GET("/history/:id",func(ctx *gin.Context) {
			handles.HandleHistoryReport(ctx)
		})
		// xem chi tiết báo cáo từ trang báo cáo
		routeReport.GET("/:id", func(ctx *gin.Context) {
			handles.DetailReport(ctx)
		})
		//get số lượng tài khoản chung số điện thoại
		routeReport.GET("/detail/account",func(ctx *gin.Context) {
			handles.GetListAccountByDetail(ctx)
		})
		// get số lượng bí danh
		routeReport.GET("/detail/alias",func(ctx *gin.Context) {
			handles.GetListLinkByDetail(ctx)
		})
		routeReport.GET("/reporter/filter/:id",func(ctx *gin.Context) {
			handles.FilterReportByReporter(ctx)
		})
		// lấy danh sách báo cáo của reporter
		routeReport.GET("/reporter/:id",func(ctx *gin.Context) {
			handles.HandleListReportByReporter(ctx)
		})
		// lọc theo lí do của người tố cáo
		
		//lấy danh sách bị báo cáo theo id người bị tố cáo
		routeReport.GET("/reported-user/list",func(ctx *gin.Context) {
			handles.GetListReportByPeerId(ctx)
		})
		//xem chi tiết 1 báo cáo của người bị báo cáo

		routeReport.GET("/reported-user/:id",func(ctx *gin.Context) {
			
			handles.HandleDetailReportByPeerId(ctx)
		})

		//chấp thuận báo cáo vi phạm
		routeReport.DELETE("/:id/:process",func(ctx *gin.Context){
			handles.HandleAccessOrDenyReport(ctx)
		})
		//lọc danh sách vi phạm theo lí do người bị vi phạm
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

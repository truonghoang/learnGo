package handles

import (
	"strconv"
	"sync"
	"truonghoang/go-scam/api/query"
	"truonghoang/go-scam/connection"
	"truonghoang/go-scam/response"

	"github.com/gin-gonic/gin"
)

type FormReport struct {
	RpterId int `json:"reporter_id" `
	RptedId int `json:"reported_id"`
	Message string `json:"message"`
	Reason string `json:"reason"`
}
func ListReport(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	parseLimit, err := strconv.Atoi(limit)
	if parseLimit <= 0 {
		parseLimit = 1
	}
	if err != nil {
		response.Res400(ctx, "Invalid limit")
		return
	}
	parsePage, err := strconv.Atoi(page)

	if parsePage <= 0 {
		parsePage = 1
	}
	if err != nil {
		response.Res400(ctx, "Invalid page")
		return
	}

	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}

	// query
	var wg sync.WaitGroup
	ch_report := make(chan query.Response)
	wg.Add(1)

	go query.QueryListReport(db, parsePage, parseLimit, ch_report, &wg)
	dataResponse := <-ch_report

	go func() {
		wg.Wait()
		close(ch_report)
		db.Close()
	}()
	if dataResponse.Err {
		response.Res400(ctx, "query db fail")
		return
	}

	response.Res200(ctx, "list data", dataResponse)
}

func DetailReport(ctx *gin.Context) {

}

func AddReport (ctx *gin.Context){

	var reportInfo FormReport
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "err db:"+err.Error())
		return
	}
	if err := ctx.BindJSON(&reportInfo); err != nil {
		response.Res400(ctx, "err bind:"+err.Error())
		return
	}

	var wg sync.WaitGroup
	ch_add := make(chan bool)
	wg.Add(1)

	go query.InsertReport(db,reportInfo.RpterId,reportInfo.RptedId,reportInfo.Message,reportInfo.Reason,ch_add,&wg)
	go func (){
		wg.Wait()
		db.Close()
	}()

}

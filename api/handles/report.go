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
	RpterId int    `json:"reporter_id" `
	RptedId int    `json:"reported_id"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
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
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "err db:"+err.Error())
		return
	}

	var wg sync.WaitGroup

	ch_detail := make(chan query.ResponseDetail)

	wg.Add(1)

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.Res400(ctx, "parse id failure")
		return
	}
	go query.QueryDetailReport(db, id, ch_detail, &wg)

	detail := <-ch_detail
	go func() {
		wg.Wait()
		close(ch_detail)
		db.Close()
	}()

	if detail.Err {
		response.Res400(ctx, "get user  failure")
		return
	}

	response.Res200(ctx, "get user success", detail.Data)

}

func AddReport(ctx *gin.Context) {

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

	go query.InsertReport(db, reportInfo.RpterId, reportInfo.RptedId, reportInfo.Message, reportInfo.Reason, ch_add, &wg)

	success := <-ch_add

	go func() {
		wg.Wait()
		close(ch_add)
		db.Close()
	}()

	if !success {
		response.Res400(ctx, "send report failure")
		return
	}
	response.Res201(ctx, "send report successfully")

}

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
type DeleteResponse struct {
	Data bool `json:"data"`
}

type BanPayload struct {
	PeerId   int    `json:"peer_id"`
	Ban      int    `json:"ban"`
	Reason   int    `json:"reason"`
	Phone    string `json:"phone"`
	AdminBan string
}

func ListReport(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	orderBy := ctx.Query("sort")
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

	go query.QueryListReport(db,orderBy, parsePage, parseLimit, ch_report, &wg)

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

func FilterReportByReason(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	reason := ctx.Query("reason")
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
	parserReason, err := strconv.Atoi(reason)
	if err != nil {
		response.Res400(ctx, "Invalid reason")
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

	go query.FilterByTypeReason(db, parserReason, parsePage, parseLimit, ch_report, &wg)
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

	ch_detail := make(chan query.Basic)

	wg.Add(1)

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.Res400(ctx, "parse id failure")
		return
	}
	infoUser, err := query.CountAccountAndLinkByPhone(db, id)
	if err != nil {
		
		response.Res400(ctx, "get data fail")
		return
	}
	go query.BasicDetail(db, infoUser.PeerId, infoUser.Phone, ch_detail, &wg)

	detail := <-ch_detail

	result := query.RecordDetail{}
	result.FirstName = infoUser.FirstName
	result.LastName = infoUser.LastName
	result.Phone = infoUser.Phone
	result.TotalAccount = detail.TotalAccount
	result.TotalLink = detail.TotalLink
	result.PeerId = infoUser.PeerId

	if detail.Err {
		response.Res400(ctx, "get report  failure")
		return
	}

	go func() {
		wg.Wait()
		close(ch_detail)
		db.Close()
	}()

	response.Res200(ctx, "get report success", result)

}

func HandleAccessOrDenyReport(ctx *gin.Context) {
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect db failure")
		return
	}
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Res400(ctx, "parser param error")
	}
	deletedParam, err := strconv.Atoi(ctx.Param("process"))
	if err != nil {
		response.Res400(ctx, "parser param error")
	}
	var wg sync.WaitGroup
	ch_delete := make(chan bool)
	wg.Add(1)
	
	go query.DeleteReport(db, idParam,deletedParam, ch_delete, &wg)
	success := <-ch_delete
	if !success {
		response.Res400(ctx, "process failure")
		return
	}
	response.Res200(ctx, "access successfully", DeleteResponse{Data: true})
}


func HandleProcessReadReport(ctx *gin.Context) {
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect db failure")
		return
	}
	bodyParams:=query.ProcessReport{}

	if err := ctx.BindJSON(&bodyParams); err != nil {
		response.Res400(ctx, "err bind:"+err.Error())
		return
	}
	var wg sync.WaitGroup
	ch_delete := make(chan bool)
	wg.Add(1)
	
	go query.ProcessReadReport(db, bodyParams.Id,bodyParams.Process, ch_delete, &wg)
	success := <-ch_delete
	if !success {
		response.Res400(ctx, "process failure")
		return
	}
	response.Res200(ctx, "access successfully", DeleteResponse{Data: true})
}

func SearchReport(ctx *gin.Context) {

	page := ctx.Query("page")
	limit := ctx.Query("limit")
	search := ctx.Query("keySearch")

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

	go query.SearchReportByKeySearch(db, search, parseLimit, parsePage, ch_report, &wg)

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

func GetListAccountByDetail(ctx *gin.Context) {
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}
	phone := ctx.Query("phone")

	result, err := query.QueryListAccountWithAndNumberReport(db, phone)
	if err != nil {
		response.Res400(ctx, "query failure")
		return
	}
	db.Close()
	response.Res200(ctx, "list account success", result)
}

func GetListLinkByDetail(ctx *gin.Context) {
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		response.Res400(ctx, "fail parser")
		return
	}

	result, err := query.QueryListLinkByPeerId(db, id)
	if err != nil {
		response.Res400(ctx, "query failure")
		return
	}
	db.Close()
	response.Res200(ctx, "list account success", result)
}

func GetListReportByPeerId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))

	if err != nil {
		response.Res400(ctx, "parser failure")
		return
	}
	orderBy := ctx.Query("sort")
	
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect db failure")
		return
	}
	result, err := query.QueryListReportByPeerId(db, id,orderBy)
	if err != nil {
		
		response.Res400(ctx, "query failure")
		return
	}
	db.Close()
	response.Res200(ctx, "get list success", result)
}



func HandleDetailReportByPeerId(ctx *gin.Context){
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Res400(ctx, "parser param error")
	}
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}

	success,err:= query.DetailReportByPeerId(db,id)
	if err!=nil {
		response.Res400(ctx,err.Error())
		return
	}
	response.Res200(ctx,"get history success",success)
}



func HandleListReportByReporter(ctx *gin.Context){
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Res400(ctx, "parser param error")
	}
	orderBy := ctx.Query("sort")
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}
	success,err:= query.QueryListReportByReporter(db,id,orderBy)
	if err!=nil {
		response.Res400(ctx,err.Error())
		return
	}
	response.Res200(ctx,"get list by reporter success",success)
}

func FilterReportByReporter(ctx *gin.Context) {

	reason := ctx.Query("reason")
	
	parserReason, err := strconv.Atoi(reason)
	if err != nil {
		response.Res400(ctx, "Invalid reason")
		return
	}

    reporter := ctx.Param("id")
	
	parserId, err := strconv.Atoi(reporter)
	if err != nil {
		response.Res400(ctx, err.Error())
		return
	}
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}

	// query
	var wg sync.WaitGroup
	ch_report := make(chan query.ResponseFilterOwnerByReason)
	wg.Add(1)

	go query.FilterReportOfReporterByReason(db, parserReason,parserId, ch_report, &wg)
	dataResponse := <-ch_report

	go func() {
		wg.Wait()
		close(ch_report)
		db.Close()
	}()
	if dataResponse.Error {
		response.Res400(ctx, "query db fail")
		return
	}

	response.Res200(ctx, "list data", dataResponse.Data)
}

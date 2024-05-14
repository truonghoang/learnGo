package handles

import (
	"fmt"
	"strconv"
	"truonghoang/go-scam/connection"
	"truonghoang/go-scam/response"

	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	Page      int `json:"page"`
	TotalPage int `json:"totalPage"`
	Data      []Product `json:"data"`
}

type Product struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Price string `json:"price" db:"price"`
	Quality string `json:"quality" db:"quality"`

}

const getProduct = "select * from product "
const baseLimit = "limit"
const offset ="offset"
const orderBy ="Order By"



func ListUserScam(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	parseLimit, err := strconv.Atoi(limit)
	if parseLimit<=0{
		parseLimit =1
	} 
	if err != nil  {
		response.Res400(ctx,"Invalid limit")
		return
	}
	parsePage, err := strconv.Atoi(page)
	fmt.Println(parsePage)
	
	if parsePage <=0 {
		parsePage =1
	}
	if err != nil  {
		response.Res400(ctx,"Invalid page")
		return
	}
	
	offset := (parsePage - 1) * parseLimit
    
	db,err:= connection.ConnectDb()
	if err!=nil{
		response.Res400(ctx,"connect Db fail")
		return
	}
	fmt.Println(parseLimit,offset)
	dataQuery:= []Product{}
	orderBy:=" id "
	sortOrder:=" ASC "
    if err:= db.Select(&dataQuery, "SELECT * FROM product ORDER BY" +orderBy+ sortOrder+ "LIMIT ? OFFSET ? ",parseLimit,offset);err!=nil{
		response.Res400(ctx,err.Error())
		return
	}

	result := DataResponse{
		Page: 1,
		TotalPage: 1,
		Data: dataQuery,
	}

	response.Res200(ctx,"list data",result)
}

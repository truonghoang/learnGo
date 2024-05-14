package handles

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	Page      int `json:"page"`
	TotalPage int `json:"totalPage"`
	data      []interface{}
}

func ListUserScam(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	parsePage, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}
	parseLimit, err := strconv.Atoi(limit)
	if err != nil {
		panic(err)
	}

	if parsePage <= 0 || parseLimit <= 0 {
		panic("Invalid page or limit value")
	}
	offet := (parsePage - 1) * parseLimit
    fmt.Println(offet)
}

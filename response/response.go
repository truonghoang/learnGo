package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Res400(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code":    0,
		"status":  http.StatusBadRequest,
		"message": message,
	})

}

func Res201(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusCreated, gin.H{
		"code":    1,
		"status":  http.StatusCreated,
		"message": message,
	})

}

func Res403(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"code":    0,
		"status":  http.StatusForbidden,
		"message": "Authorization invalid",
	})

}

func Res401(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"code":    0,
		"status":  http.StatusUnauthorized,
		"message": "token invalid",
	})

}

func Res200(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"status":  http.StatusOK,
		"message": message ,
		"data":    data,
	})

}

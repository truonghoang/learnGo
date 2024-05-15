package main

import (
	"net/http"
	"truonghoang/go-scam/api/middleware"
	"truonghoang/go-scam/api/routes"

	"github.com/gin-gonic/gin"
)

const port = "localhost:80"

func main() {

	r := gin.Default()
	r.Use(middleware.ConfigCors())
	r.LoadHTMLGlob("dist/*.html")
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })
    
	routes.RouteApi(r)
	r.Run(port)

}

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"truonghoang/go-scam/api/middleware"
	"truonghoang/go-scam/api/routes"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const port = "localhost:80"

func main() {
	f, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // Thiết lập gin.DefaultWriter để ghi log vào file và stdout
    gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	
	r := gin.Default()
	r.Use(middleware.ConfigCors())
	r.Use(static.Serve("/assets/*",static.LocalFile("/dist/assets",false)))
	r.LoadHTMLGlob("dist/*.html")
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })
    
	routes.RouteApi(r)
	r.Run(port)

}

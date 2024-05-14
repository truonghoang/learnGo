package main

import (
	"github.com/gin-gonic/gin"

	"truonghoang/go-scam/api/middleware"
	"truonghoang/go-scam/api/routes"
)

const port = "localhost:80"

func main() {

	r := gin.Default()
	r.Use(middleware.ConfigCors())
	routes.RouteApi(r)

	r.Run(port)

}

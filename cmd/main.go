package main

import (
	"fmt"
	"net/http"
	"os"
	"truonghoang/go-scam/config"
	"truonghoang/go-scam/connection"

	"github.com/gin-gonic/gin"
)

type DataLogin struct {
	Email    string `json:"email"`
	Password string `json:password`
}

type Register struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Account struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

const port = "localhost:8080"

const schema = `CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);`

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	db, err := connection.ConnectDb()
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected to Db", db)
	// db.MustExec(schema)
	fmt.Println(pwd)
	cfg, err := config.LoadConfig(pwd + "/config/config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/login", func(c *gin.Context) {
		var loginAcc DataLogin
		if err := c.BindJSON(&loginAcc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "login success",
			"token":   loginAcc,
		})
	})
	router.POST("/register", func(ctx *gin.Context) {
		var account Register
		if err := ctx.BindJSON(&account); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err,
			})
		}
		fmt.Println(account.FirstName)
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)", account.FirstName, account.LastName, account.Email)
		// tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (?,?,?)", "John", "Doe", "johndoeDNE@gmail.net")

		if err:=tx.Commit();err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err,
			})
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "create successfully",
		})


	})
	router.GET("/account/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		people := []Account{}
    db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
		ctx.JSON(http.StatusOK, gin.H{
			"messgage": "success",
			"name":id,
			"data": people,
		})
	})

	router.Run(port)

}

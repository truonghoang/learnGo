package handles

import (
	"time"
	"truonghoang/go-scam/connection"
	"truonghoang/go-scam/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginAccount struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterAccount struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Account struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type CustomClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type ResponseData struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func Register(ctx *gin.Context) {
	var account RegisterAccount
	// connect db
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, err.Error())
		return
	}
	// bind json data from payload
	if err := ctx.BindJSON(&account); err != nil {
		response.Res400(ctx, err.Error())
		return
	}
	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), 10)
	if err != nil {
		response.Res400(ctx, err.Error())
	}

	//insert db
	tx := db.MustBegin()

	tx.MustExec("INSERT INTO account(name,email,password) VALUES(?,?,?)", account.Name, account.Email, hashPassword)

	if err := tx.Commit(); err != nil {
		response.Res400(ctx, err.Error())
	}

	//res
	response.Res201(ctx, "register success")

}

func Login(ctx *gin.Context) {
	var account LoginAccount
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, err.Error())
		return
	}
	if err := ctx.BindJSON(&account); err != nil {
		response.Res400(ctx, err.Error())
		return
	}
	// query check email existed

	var accountCheck Account

	if err := db.Get(&accountCheck, "SELECT id,email,password FROM account WHERE email=?", account.Email); err != nil {
		response.Res400(ctx, "email not exist")
		return
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(accountCheck.Password), []byte(account.Password)); err != nil {
		response.Res400(ctx, "password invalid")
		return
	}
	// generate token

	var recretKey = []byte("scamreportserver")

	claims := CustomClaims{
		Id:    accountCheck.Id,
		Email: accountCheck.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := token.SignedString(recretKey)

	if err != nil {
		response.Res400(ctx, err.Error())
		return
	}

	dataRes := ResponseData{
		Email: accountCheck.Email,
		Token: secret,
	}
	response.Res200(ctx, "login successfully", dataRes)

}

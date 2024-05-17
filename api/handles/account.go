package handles

import (
	"strconv"
	"sync"
	"time"
	"truonghoang/go-scam/api/query"
	"truonghoang/go-scam/connection"
	"truonghoang/go-scam/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type Account struct {
	Id       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
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
	var account RegisterUser
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
	var wg sync.WaitGroup
	ch_user := make(chan bool)

	wg.Add(1)

	go query.InsertUser(db, account.Phone, account.FirstName, account.LastName, string(hashPassword), ch_user, &wg)

	success := <-ch_user

	if !success {
		response.Res400(ctx, "add user failure")
		return
	}

	//start query get id user added by phone

	ch_id := make(chan int)

	wg.Add(1)

	go query.SelectUserId(db, account.Phone, ch_id, &wg)

	id_user := <-ch_id

	if id_user == 0 {
		response.Res400(ctx, "user not exist")
		return
	}

	// start insert user_name

	wg.Add(1)

	go query.InsertUserName(db, id_user, account.Email, ch_user, &wg)

	end := <-ch_user

	if !end {
		response.Res400(ctx, "insert into user_name failure")
		return
	}
	go func() {
		wg.Wait()
		close(ch_user)
		close(ch_id)
		db.Close()
	}()

	response.Res201(ctx, "register success")

}

func Login(ctx *gin.Context) {
	var account LoginUser
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "err db:"+err.Error())
		return
	}
	if err := ctx.BindJSON(&account); err != nil {
		response.Res400(ctx, "err bind:"+err.Error())
		return
	}
	// query check email existed
	var wg sync.WaitGroup
	ch_login := make(chan query.ResultLogin)

	wg.Add(1)

	go query.QueryEmail(db, account.Email, ch_login, &wg)

	resultData := <-ch_login

	go func() {
		wg.Wait()
		close(ch_login)
		db.Close()
	}()

	if resultData.Error {
		response.Res400(ctx, "request data login failure")
		return
	}
	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(resultData.User.Password), []byte(account.Password)); err != nil {
		response.Res400(ctx, "password invalid")
		return
	}

	// generate token

	var recretKey = []byte("scamreportserver")

	claims := CustomClaims{
		Id:    strconv.Itoa(resultData.User.Id),
		Email: resultData.User.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := token.SignedString(recretKey)

	if err != nil {
		response.Res400(ctx, "sign token:"+err.Error())
		return
	}

	dataRes := ResponseData{
		Email: resultData.User.Email,
		Token: secret,
	}
	response.Res200(ctx, "login successfully", dataRes)

}

func GetDetailUser(ctx *gin.Context, search bool) {
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "err db:"+err.Error())
		return
	}

	var wg sync.WaitGroup

	ch_detail := make(chan query.ResultUserInfo)

	wg.Add(1)
	if search {
		phone := ctx.Query("phone")
		go query.QueryUserByPhone(db, phone, ch_detail, &wg)
	} else {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			response.Res400(ctx, "parse id failure")
			return
		}
		go query.QueryUserById(db, id, ch_detail, &wg)
	}

	detail := <-ch_detail
	go func() {
		wg.Wait()
		close(ch_detail)
		db.Close()
	}()

	if detail.Error {
		response.Res400(ctx, "get user  failure")
		return
	}

	response.Res200(ctx, "get user success", detail.User)

}

func ListUser (ctx *gin.Context){
	db,err := connection.ConnectDb()
	
	if err!=nil {
		response.Res400(ctx,"connect db failure")
		return
	}
	ch_list_user := make(chan query.ResponseListUser)
	var wg sync.WaitGroup
	wg.Add(1)
	parsePage,err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Res400(ctx,"parser page err")
		return
	}
	parseLimit,err:= strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Res400(ctx,"parser limit err")
		return
	}

	go query.ListUser(db,parseLimit,parsePage,ch_list_user,&wg)

	resultQuery := <- ch_list_user
	if resultQuery.Err {
		response.Res400(ctx,"query data failure")
		return
	}
	go func (){
		wg.Wait()
		close(ch_list_user)
		db.Close()
	}()

	response.Res200(ctx,"list user successfully",resultQuery)
}


func SelectUser (ctx *gin.Context){
	db,err := connection.ConnectDb()
	
	if err!=nil {
		response.Res400(ctx,"connect db failure")
		return
	}
	ch_sel_user := make(chan query.ResponseUserName)
	var wg sync.WaitGroup
	wg.Add(1)
	
	go query.ListUserSelect(db,ch_sel_user,&wg)

	resultQuery := <- ch_sel_user
	if resultQuery.Err {
		response.Res400(ctx,"query data failure")
		return
	}
	go func (){
		wg.Wait()
		close(ch_sel_user)
		db.Close()
	}()

	response.Res200(ctx,"list user successfully",resultQuery.Data)
}


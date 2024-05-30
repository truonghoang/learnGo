package handles

import (
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
	"truonghoang/go-scam/api/query"
	"truonghoang/go-scam/config"
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



func Login(ctx *gin.Context) {
	var account LoginUser
	if err := ctx.BindJSON(&account); err != nil {
		response.Res400(ctx, "err bind:"+err.Error())
		return
	}
	pwd, err := os.Getwd()
	if err != nil {
		response.Res400(ctx, "Not load path from system")
		return
	}
	cfg, err := config.LoadConfigAccount(pwd + "/config/admin.json")
	if err != nil {
		response.Res400(ctx, "Not load account from system")
		return
	}
	// Define a regex pattern for validating email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	success := emailRegex.MatchString(account.Email)
	if !success {
		response.Res400(ctx, "email invalid")
		return
	}
	exist := false
	userSystem := config.Account{}
	
	for _, acc := range cfg.Data {
		
		if acc.Account == account.Email {
			userSystem.Password = acc.Password
			userSystem.Role = acc.Role
			userSystem.Account = acc.Account
			exist = true
			break
		}

	}
	if !exist {
		response.Res400(ctx, "email not exist")
		return
	}
	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(userSystem.Password), []byte(account.Password)); err != nil {
		response.Res400(ctx, "password invalid")
		return
	}

	// generate token

	var recretKey = []byte("scamreportserver")

	claims := CustomClaims{
		Id:    userSystem.Role,
		Email: userSystem.Account,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(55 * time.Minute).Unix(),
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
		Email: userSystem.Account,
		Token: secret,
	}
	response.Res200(ctx, "login successfully", dataRes)

}

func SearchBannedUser(ctx *gin.Context){
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	sort := ctx.Query("sort")
	ban := ctx.Query("ban")
	keySearch := ctx.Query("keySearch")
	parseBan, err := strconv.Atoi(ban)
	if err != nil {
		response.Res400(ctx, "Invalid ban")
		return
	}
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

	result, err := query.QuerySearchUser(db, parsePage, parseLimit, sort, parseBan,keySearch)

	if err != nil {
		response.Res400(ctx, err.Error())
		return
	}

	response.Res200(ctx, "get list user banned successfully", result)
}

func BanAndUnBanUser(ctx *gin.Context) {
	var reportInfo BanPayload
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "err db:"+err.Error())
		return
	}
	if err := ctx.BindJSON(&reportInfo); err != nil {
		response.Res400(ctx, "err bind:"+err.Error())
		return
	}

	admin := ctx.MustGet("email").(string)

	reportInfo.AdminBan = admin
	
	success, err := query.InsertBanUserAndHistory(db, reportInfo.PeerId, reportInfo.Ban, reportInfo.Reason, reportInfo.AdminBan)

	if err != nil && success == 0 {
		response.Res400(ctx, err.Error())
		return
	}
	result, err := query.BannedOrtherAccountWithPhone(db, reportInfo.PeerId, reportInfo.Ban, reportInfo.Phone, reportInfo.Reason, reportInfo.AdminBan)

	if err != nil && result == 0 {
		response.Res400(ctx, "Query db fail")
		return
	}
	response.Res201(ctx, "Ban successfully")

}

func ListUserBan(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	sort := ctx.Query("sort")
	ban := ctx.Query("ban")
	parseBan, err := strconv.Atoi(ban)
	if err != nil {
		response.Res400(ctx, "Invalid ban")
		return
	}
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

	result, err := query.ListUserBanned(db, parsePage, parseLimit, sort, parseBan)

	if err != nil {
		response.Res400(ctx, "query failure")
		return
	}

	response.Res200(ctx, "get list user banned successfully", result)

}

func HandleHistoryReport(ctx *gin.Context){

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Res400(ctx, "parser param error")
	}
	db, err := connection.ConnectDb()
	if err != nil {
		response.Res400(ctx, "connect Db fail")
		return
	}

	success,err:= query.GetHistoryBan(db,id)
	if err!=nil {
		response.Res400(ctx,err.Error())
		return
	}
	response.Res200(ctx,"get history success",success.Data)

}


func FilterReportBannedByReason(ctx *gin.Context) {

	reason := ctx.Query("reason")
	
	parserReason, err := strconv.Atoi(reason)
	if err != nil {
		response.Res400(ctx, "Invalid reason")
		return
	}

 peer_id := ctx.Query("id")
	
	parserId, err := strconv.Atoi(peer_id)
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
	ch_report := make(chan query.ResponseFilterOwnerByReason)
	wg.Add(1)

	go query.FilterOwnerByTypeReason(db, parserReason,parserId, ch_report, &wg)
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

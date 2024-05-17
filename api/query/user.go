package query

import (
	"fmt"
	"math"
	"sync"

	"github.com/jmoiron/sqlx"
)

type UserLogin struct {
	Id       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
type UserInfo struct {
	Id        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
}

type UserName struct {
	Id        int    `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
}

type ResponseUserName struct {
	Data []UserName `json:"data"`
	Err  bool
}

type ResultUserInfo struct {
	User  UserInfo
	Error bool
}

type ResponseListUser struct {
	Data  []UserInfo `json:"data"`
	Err   bool       `json:"error"`
	Limit int        `json:"limit"`
	Page  int        `json:"page"`
	Total int        `json:"totalPage"`
}

type ResultLogin struct {
	User  UserLogin
	Error bool
}

func ListUserSelect(db *sqlx.DB, ch chan ResponseUserName, wg *sync.WaitGroup) {
	defer wg.Done()
	const limit = 100
	var resultResponse ResponseUserName
	result := []UserName{}

	err := db.Select(&result, `select id,first_name,last_name from user limit ?`, limit)
	if err != nil {
		resultResponse.Err = true
		ch <- resultResponse
		return
	}
	resultResponse.Err = false
	resultResponse.Data = result

	ch <- resultResponse
}

func ListUser(db *sqlx.DB, limit int, page int, ch chan ResponseListUser, wg *sync.WaitGroup) {
	defer wg.Done()
	offset := (page - 1) * limit
	var resultResponse ResponseListUser
	result := make([]UserInfo, limit)

	err := db.Select(&result, `select user.id,user.first_name,user.last_name,user.phone,user_name.email from user JOIN user_name ON user.id=user_name.uid  limit ? offset ?`, limit, offset)
	if err != nil {
		resultResponse.Err = true
		ch <- resultResponse
		return
	}

	var count int
	countQuery := `SELECT count(id) as totalPage from user`
	error2 := db.QueryRow(countQuery).Scan(&count)

	if error2 != nil {
		resultResponse.Err = true
		ch <- resultResponse
		return
	}
	totalPage:=math.Ceil(float64(count/limit))
	fmt.Println(count,totalPage)
	resultResponse.Err = false
	resultResponse.Data = result
	resultResponse.Limit = limit
	resultResponse.Page = page
	resultResponse.Total = int(totalPage)

	ch <- resultResponse
}

func InsertUser(db *sqlx.DB, phone string, first_name string, last_name string, password string, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO user(phone, first_name, last_name, password) VALUES(?, ?, ?, ?)", phone, first_name, last_name, password)
	if err := tx.Commit(); err != nil {
		fmt.Println("error user:" + err.Error())
		ch <- false
		return
	}
	ch <- true
}

func InsertUserName(db *sqlx.DB, uid int, email string, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO user_name(uid,email) VALUES(?,?)", uid, email)
	if err := tx.Commit(); err != nil {
		fmt.Println("ins user_name:" + err.Error())
		ch <- false
		return
	}
	ch <- true
}

func SelectUserId(db *sqlx.DB, phone string, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	var userId int

	err := db.Get(&userId, "SELECT id FROM user WHERE phone = ?", phone)
	if err != nil {
		ch <- 0
	}
	ch <- userId
}

func QueryEmail(db *sqlx.DB, email string, ch chan ResultLogin, wg *sync.WaitGroup) {
	defer wg.Done()
	resultQuery := ResultLogin{}
	err := db.Get(&resultQuery.User, "SELECT user.id,user.password,user_name.email from user inner join user_name on user.id=user_name.uid where user_name.email= ?", email)
	if err != nil {

		resultQuery.Error = true
		ch <- resultQuery
	}
	resultQuery.Error = false
	fmt.Println(resultQuery)
	ch <- resultQuery

}

func QueryUserById(db *sqlx.DB, id int, ch chan ResultUserInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	resultQuery := ResultUserInfo{}
	err := db.Get(&resultQuery.User, "SELECT user.id,user.phone,user.first_name,user.last_name,user_name.email FROM user inner join user_name ON user.id=user_name.uid WHERE user.id= ?", id)
	if err != nil {

		resultQuery.Error = true
		ch <- resultQuery
	}
	resultQuery.Error = false
	fmt.Println(resultQuery)
	ch <- resultQuery

}

func QueryUserByPhone(db *sqlx.DB, phone string, ch chan ResultUserInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	resultQuery := ResultUserInfo{}
	err := db.Get(&resultQuery.User, "SELECT user.id,user.phone,user.first_name,user.last_name,user_name.email FROM user inner join user_name ON user.id=user_name.uid WHERE user.phone= ?", phone)
	if err != nil {

		resultQuery.Error = true
		ch <- resultQuery
	}
	resultQuery.Error = false
	fmt.Println(resultQuery)
	ch <- resultQuery

}

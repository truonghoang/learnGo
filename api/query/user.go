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
	Link     string `json:"link" db:"link"`
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

type ResponseUserWithPhone struct {
	User []UserInfo
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

	err := db.Select(&result, `select users.id,users.first_name,users.last_name,users.phone,username.link from users JOIN username ON users.id=username.peer_id AND username.peer_type=2 ORDER by id DESC  limit ? offset ?`, limit, offset)
	if err != nil {
		resultResponse.Err = true
		ch <- resultResponse
		return
	}

	var count float64
	countQuery := `SELECT count(id) as totalPage from users`
	error2 := db.QueryRow(countQuery).Scan(&count)

	if error2 != nil {
		resultResponse.Err = true
		ch <- resultResponse
		return
	}
	parseLimit := float64(limit)
	totalPage:=math.Ceil(float64(count/parseLimit))
	
	fmt.Println(count,limit,totalPage)
	
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
	err := db.Get(&resultQuery.User, "SELECT users.id,user.phone,user.first_name,user.last_name,user_name.email FROM user inner join user_name ON user.id=user_name.uid WHERE user.phone= ?", phone)
	if err != nil {

		resultQuery.Error = true
		ch <- resultQuery
	}
	resultQuery.Error = false
	fmt.Println(resultQuery)
	ch <- resultQuery

}

func FilterUsersWithPhone(db *sqlx.DB,phone string,ch chan ResponseUserWithPhone,wg *sync.WaitGroup ){
	defer wg.Done()
	
	var resultResponse ResponseUserWithPhone
	result := []UserInfo{}

	err := db.Select(&result, `select users.id,users.first_name,users.last_name,users.phone,username.link from users JOIN username ON users.id=username.peer_id ORDER by id DESC `)
	if err != nil {
		resultResponse.Error = true
		ch <- resultResponse
		return
	}

	resultResponse.Error = false
	resultResponse.User = result
	
	ch <- resultResponse
}
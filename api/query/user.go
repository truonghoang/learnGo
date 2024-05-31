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

type ResponseHistoryban struct {
	Data []HistoryBan
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

func QuerySearchUser(db *sqlx.DB,page int,limit int,orderBy string,ban int, keySearch string)(*ResponseUserBan, error){
	offset := (page - 1) * limit
	var responseUser ResponseUserBan
	usersBan := []UserBanned{}
	search:= `%`+keySearch+`%`
	query := `SELECT u.first_name,ub.uid ,u.last_name ,u.phone,ub.created_at ,ub.ban from ban_user ub join users u on ub.uid = u.id where ub.ban= ? and u.first_name like ? or u.last_name like ? or u.phone like ? ORDER BY ub.created_at ` + orderBy + ` limit ? offset ?`

	countQuery := `SELECT count(ub.id) as totalPage from ban_user ub join users u on ub.uid = u.id  where ub.ban= ? and u.first_name like ? or u.last_name like ? or u.phone like ? `
	err := db.Select(&usersBan, query, ban,search,search,search, limit, offset)

	if err != nil {
		
		return nil, err
	}
	var count float64
	error2 := db.QueryRow(countQuery, ban,search,search,search).Scan(&count)
	if error2 != nil {
		return nil, error2
	}

	totalPage := math.Ceil(float64(count / float64(limit)))
	responseUser.Data = usersBan
	responseUser.Limit = limit
	responseUser.Page = page
	responseUser.TotalPage = int(totalPage)
	return &responseUser, nil

}


func GetHistoryBan(db *sqlx.DB, id int) (*ResponseHistoryban, error) {
	var responseBan ResponseHistoryban
	result := responseBan.Data

	query := `select h.id,h.ban,h.reason,h.admin_ban,h.created_at, u.first_name,u.last_name,u.phone from ban_history h join users u on h.uid=u.id where uid=? order by created_at desc`

	err := db.Select(&result, query, id)
	if err != nil {

		return nil, err
	}
	responseBan.Data = result

	return &responseBan, nil
}



func InsertBanUserAndHistory(db *sqlx.DB, peer_id int, ban int, reason int, admin_ban string) (int, error) {

	var user CheckBan
	err := db.Get(&user, "SELECT id, uid FROM ban_user WHERE  uid = ?", peer_id)
	if err != nil {
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO ban_user(uid,ban) VALUES(?,?)", peer_id, ban)
		tx.MustExec("INSERT INTO ban_history(uid,ban,reason,admin_ban) VALUES(?,?,?,?)", peer_id, ban, reason, admin_ban)
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return 1, nil
	} else {
		user.Ban = ban
		_, err := db.NamedExec("UPDATE ban_user SET ban = :ban  WHERE uid = :uid", user)

		if err != nil {
			return 0, err
		}

		tx := db.MustBegin()
		tx.MustExec("INSERT INTO ban_history(uid,ban,reason,admin_ban) VALUES(?,?,?,?)", peer_id, ban, reason, admin_ban)
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return 1, nil
	}

}

func BannedOrtherAccountWithPhone(db *sqlx.DB, peer_id int, ban int, phone string, reason int, admin_ban string) (int, error) {
	if phone ==""{
		return 1,nil
	}
	uid := []UidBanned{}
	err := db.Select(&uid, "SELECT id FROM users WHERE user_type!= 9 AND phone= ?", phone)
	if err != nil {
		return 1, nil
	} else {
			fmt.Print(uid)
		for i, user_id := range uid {
			if user_id.Id == peer_id {
				continue;
			}
			var user CheckBan
			err := db.Get(&user, "SELECT id, uid FROM ban_user WHERE  uid = ?", user_id.Id)
			fmt.Println(user)

			if err != nil {
				tx := db.MustBegin()
				tx.MustExec("INSERT INTO ban_user(uid,ban) VALUES(?,?)", user_id.Id, ban)
				tx.MustExec("INSERT INTO ban_history(uid,ban,reason,admin_ban) VALUES(?,?,?,?)", user_id.Id, ban, reason, admin_ban)
				if err := tx.Commit(); err != nil {
					return 0, err

				}
				
			} else {
				user.Ban = ban
				_, err := db.NamedExec("UPDATE ban_user SET ban = :ban  WHERE uid = :uid", user)
				if err != nil {
					return 0, err

				}
				tx := db.MustBegin()
				tx.MustExec("INSERT INTO ban_history(uid,ban,reason,admin_ban) VALUES(?,?,?,?)", user_id.Id, ban, reason, admin_ban)
				if err := tx.Commit(); err != nil {
					return 0, err

				}
				
			}
			if i==len(uid)-1 {
				return 1,nil
			}

		}
			
		return 1, nil
	}

}
func ListUserBanned(db *sqlx.DB, page int, limit int, orderBy string, ban int) (*ResponseUserBan, error) {
	offset := (page - 1) * limit
	var responseUser ResponseUserBan
	usersBan := []UserBanned{}
	query := `SELECT u.first_name,ub.uid ,u.last_name ,u.phone,ub.created_at ,ub.ban from ban_user ub join users u on ub.uid = u.id where ub.ban= ? ORDER BY ub.created_at ` + orderBy + ` limit ? offset ?`

	countQuery := `SELECT count(ub.id) as totalPage from ban_user ub  where ub.ban= ? `
	err := db.Select(&usersBan, query, ban, limit, offset)

	if err != nil {
		
		return nil, err
	}
	var count float64
	error2 := db.QueryRow(countQuery, ban).Scan(&count)
	if error2 != nil {
		return nil, error2
	}
	totalPage := math.Ceil(count / float64(limit))
	responseUser.Data = usersBan
	responseUser.Limit = limit
	responseUser.Page = page
	responseUser.TotalPage = int(totalPage)
	return &responseUser, nil

}

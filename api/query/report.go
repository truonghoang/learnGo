package query

import (
	"fmt"
	"math"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Record struct {
	Id            int    `json:"id" db:"id"`
	FirstName     string `json:"firstName" db:"first_name"`
	LastName      string `json:"last_name" db:"last_name"`
	Phone         string `json:"phone" db:"phone"`
	PeerId        int    `json:"peer_id" db:"peer_id"`
	TotalReported string `json:"total_reported" db:"total"`
	Time          string `json:"created_at" db:"created_at"`
}

type RecordDetail struct {
	FirstName    string `json:"first_name"  `
	LastName     string `json:"last_name" `
	Phone        string `json:"phone"`
	TotalLink    int    `json:"totalLink" `
	TotalAccount int    `json:"totalAccount" `
	PeerId       int    `json:"peer_id"`
}
type CountUserPhone struct {
	PeerId    int    `json:"id" db:"peer_id"`
	Phone     string `json:"phone" db:"phone"`
	FirstName string `json:"first_name" db:"first_name" `
	LastName  string `json:"last_name" db:"last_name"`
}
type Basic struct {
	TotalLink    int  `json:"totalLink" db:"totalLink"`
	TotalAccount int  `json:"totalAccount" db:"totalAccount"`
	Err          bool `json:"error"`
}
type Response struct {
	Data  []Record `json:"data"`
	Err   bool     `json:"error"`
	Limit int      `json:"limit"`
	Page  int      `json:"page"`
	Total int      `json:"totalPage"`
}

type UserBanned struct {
	Id        int    `json:"id " db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Uid       int    `json:"uid" db:"uid"`
	Phone     string `json:"phone" db:"phone"`
	Ban       int    `json:"ban" db:"ban"`
	Time      string `json:"created_at" db:"created_at"`
}
type ResponseUserBan struct {
	Data      []UserBanned `json:"data"`
	Page      int          `json:"page"`
	Limit     int          `json:"limit"`
	TotalPage int          `json:"totalPage"`
}

type ResultAccountNumberReport struct {
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Phone       string `json:"phone" db:"phone"`
	TotalReport int    `json:"totalReport" db:"totalReport"`
}

type LinkStruct struct {
	Id   int    `json:"id" db:"id"`
	Link string `json:"link" db:"link"`
}

type ReportPeerIdStruct struct {
	Id        int    `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Phone     string `json:"phone" db:"phone"`
	Reason    int    `json:"reason" db:"reason"`
	Content   string `json:"content" db:"content"`
	Time      string `json:"created_at" db:"created_at"`
}

type CheckBan struct {
	Id  int    `db:"id"`
	Uid string `db:"uid"`
	Ban int
}
type BanUserWithPhone struct {
	Id int `db:"id"`
}

type UidBanned struct {
	Id int `json:"id" db:"id"`
}
type HistoryBan struct {
	Id       int    `json:"id" db:"id"`
	Ban      int    `json:"ban" db:"ban"`
	Reason   int    `json:"reason" db:"reason"`
	AdminBan string `json:"admin_ban" db:"admin_ban"`
	Time     string `json:"created_at" db:"created_at"`
}
type ResponseHistoryban struct {
	Data []HistoryBan
}

type FilterOwnerByReason struct {
	Id        int    `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Phone     string `json:"phone" db:"phone"`
	Reason    int    `json:"reason" db:"reason"`
	Content   string `json:"content" db:"content"`
	Time      string `json:"created_at" db:"created_at"`
}
type ResponseFilterOwnerByReason struct {
	Data  []FilterOwnerByReason
	Error bool
}

type HiddenReport struct {
	Id      int `json:"id"`
	Deleted int `json:"deleted"`
}
type ProcessReport struct {
	Id      int `json:"id"`
	Process int `json:"process"`
}

func QueryListReport(db *sqlx.DB,orderBy string, page int, limit int, ch chan Response, wg *sync.WaitGroup) {
	defer wg.Done()
	offset := (page - 1) * limit
	var responseData Response
	dataResult := make([]Record, limit)
	query := `
	WITH uid_counts AS (
        SELECT peer_id, COUNT(*) AS total
         FROM reports
        GROUP BY peer_id
    )
    SELECT r.id, u2.first_name, u2.last_name, uc.peer_id,u2.phone, r.created_at, uc.total
    FROM (
        SELECT r.*
        FROM reports r
        WHERE r.peer_id != 0
        AND (r.id IN (
            SELECT MAX(id)
            FROM reports
            WHERE peer_id != 0
            GROUP BY peer_id
        ))
    ) r
    JOIN uid_counts uc ON r.peer_id = uc.peer_id
    JOIN users u1 ON r.user_id = u1.id
    JOIN users u2 ON r.peer_id = u2.id
	Order By r.created_at `+ orderBy +
    ` LIMIT ? OFFSET ?`

	err := db.Select(&dataResult, query, limit, offset)

	if err != nil {
 fmt.Print(err)
		responseData.Err = true
		ch <- responseData
		return
	}
	var count float64

	countQuery := `SELECT count(rp.id) as totalPage from  (
        SELECT r.*
        FROM reports r
        WHERE r.peer_id != 0
        AND (r.id IN (
            SELECT MAX(id)
            FROM reports
            WHERE peer_id != 0 
            GROUP BY peer_id
        ))
    ) rp JOIN users u1 ON rp.user_id=u1.id join users u2 on rp.peer_id =u2.id `

	error2 := db.QueryRow(countQuery).Scan(&count)

	if error2 != nil {

		responseData.Err = true
		ch <- responseData
		return
	}
	totalPage := math.Ceil(float64(count / float64(limit)))
	responseData.Total = int(totalPage)
	responseData.Err = false
	responseData.Data = dataResult
	responseData.Limit = limit
	responseData.Page = page

	ch <- responseData

}

func DeleteReport(db *sqlx.DB, id int, deleted int, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	reportResult := HiddenReport{}
	reportResult.Id = id
	reportResult.Deleted = deleted
	_, err := db.NamedExec("UPDATE reports SET deleted = :deleted  WHERE id = :id", reportResult)
	if err != nil {
		ch <- false
		return
	}

	ch <- true

}

func ProcessReadReport(db *sqlx.DB, id int, process int, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	reportResult := ProcessReport{}
	reportResult.Id = id
	reportResult.Process = process
	_, err := db.NamedExec("UPDATE reports SET process = :process  WHERE id = :id", reportResult)
	if err != nil {
		ch <- false
		return
	}

	ch <- true

}

func FilterByTypeReason(db *sqlx.DB, reason int, page int, limit int, ch chan Response, wg *sync.WaitGroup) {
	defer wg.Done()
	offset := (page - 1) * limit
	var responseData Response
	dataResult := make([]Record, limit)
	query := `
	WITH uid_counts AS (
        SELECT peer_id, COUNT(*) AS total
        FROM reports
		Where reports.reason = ?
        GROUP BY peer_id
    )
    SELECT r.id, u2.first_name, u2.last_name, u2.phone, r.created_at, uc.total
    FROM (
        SELECT r.*
        FROM reports r
        WHERE r.peer_id != 0
        AND (r.id IN (
            SELECT MAX(id)
            FROM reports
            WHERE peer_id != 0 AND reports.reason =?
            GROUP BY peer_id
        ))
    ) r
    JOIN uid_counts uc ON r.peer_id = uc.peer_id
    JOIN users u1 ON r.user_id = u1.id
    JOIN users u2 ON r.peer_id = u2.id
	Where NOT r.peer_id=0 AND r.reason =? order by r.id desc
    LIMIT ? OFFSET ?
	`

	countQuery := `SELECT count(rp.id) as totalPage from  (
        SELECT r.*
        FROM reports r
        WHERE r.peer_id != 0
        AND (r.id IN (
            SELECT MAX(id)
            FROM reports
            WHERE peer_id != 0 AND reports.reason =?
            GROUP BY peer_id
        ))
    ) rp JOIN users u1 ON rp.user_id=u1.id join users u2 on rp.peer_id =u2.id where rp.reason =?`
	err := db.Select(&dataResult, query, reason, reason, reason, limit, offset)

	if err != nil {
		fmt.Print("query" + err.Error())
		responseData.Err = true
		ch <- responseData
		return
	}
	var count float64
	error2 := db.QueryRow(countQuery, reason, reason).Scan(&count)
	if error2 != nil {

		responseData.Err = true
		ch <- responseData
		return
	}
	totalPage := math.Ceil(float64(count / float64(limit)))
	responseData.Total = int(totalPage)
	responseData.Err = false
	responseData.Data = dataResult
	responseData.Limit = limit
	responseData.Page = page

	ch <- responseData

}

func SearchReportByKeySearch(db *sqlx.DB, keySearch string, limit int, page int, ch chan Response, wg *sync.WaitGroup) {
	defer wg.Done()

	offset := (page - 1) * limit

	var responseData Response

	dataResult := make([]Record, limit)
   key:= `%` +keySearch + `%`
	querySearch := `
	WITH uid_counts AS (
        SELECT peer_id, COUNT(*) AS total
        FROM reports
        GROUP BY peer_id
    )
    SELECT r.id, u2.first_name,u2.last_name, u2.phone,r.created_at, uc.total
    FROM reports r
    JOIN uid_counts uc ON r.peer_id = uc.peer_id
    JOIN users u1 ON r.user_id=u1.id join users u2 on r.peer_id =u2.id 
    Where u2.phone like ? OR u2.first_name like ? OR u2.last_name like ? AND NOT r.peer_id=0 order by r.created_at desc  LIMIT ? offset ? `

	err := db.Select(&dataResult, querySearch, key,key,key, limit, offset)

	if err != nil {
    
		responseData.Err = true
		ch <- responseData
		return
	}

	var count float64

	countQuery := `SELECT count(rp.id) as totalPage from reports rp JOIN users u1 ON rp.user_id=u1.id join users u2 on rp.peer_id =u2.id  Where  u2.phone like ? OR u2.first_name like ? OR u2.last_name like ?`

	error2 := db.QueryRow(countQuery, key, key,key).Scan(&count)

	if error2 != nil {
		fmt.Print(error2)
		responseData.Err = true
		ch <- responseData
		return
	}

	totalPage := math.Ceil(float64(count / float64(limit)))

	responseData.Total = int(totalPage)

	responseData.Err = false

	responseData.Data = dataResult

	responseData.Limit = limit

	responseData.Page = page

	ch <- responseData

}

// / detail user  includes : basicDetail(), totalAccountWithPhone, totalLinkWithID
func BasicDetail(db *sqlx.DB, id int, phone string, ch chan Basic, wg *sync.WaitGroup) {
	defer wg.Done()
	var resultDetail Basic
	info := Basic{}
	if phone == "" {
		resultDetail.TotalAccount = 0
		queryDB := `
		SELECT  COUNT(*) as totalLink
		FROM username un
		WHERE un.peer_id = ?
		GROUP BY peer_id
	`
		err := db.Get(&info, queryDB, id)
		if err != nil {
			fmt.Print("err not phone:", err)
			resultDetail.Err = true
			ch <- resultDetail
			return
		}

		resultDetail.TotalLink = info.TotalLink
		resultDetail.Err = false
		ch <- resultDetail
	} else {
		queryDB := `
	WITH link_account AS (
		SELECT peer_id, COUNT(*) as totalLink
		FROM username un
		WHERE un.peer_id = ?
		GROUP BY peer_id
	)
	SELECT  uc.totalLink, COUNT(*) as totalAccount
	FROM users u1
	JOIN link_account uc 
	WHERE   u1.user_type != 9 AND u1.phone = ? 
	GROUP BY  u1.phone;
	`
		err := db.Get(&info, queryDB, id, phone)
		if err != nil {
			fmt.Print("err nè:", err)
			resultDetail.Err = true
			ch <- resultDetail
			return
		}

		resultDetail.TotalAccount = info.TotalAccount
		resultDetail.TotalLink = info.TotalLink
		resultDetail.Err = false
		ch <- resultDetail
	}

}

func CountAccountAndLinkByPhone(db *sqlx.DB, id int) (*CountUserPhone, error) {
	info := CountUserPhone{}
	queryDB := `
  	 	SELECT u.first_name,u.last_name,u.phone, r.peer_id
   		FROM reports r
		join users u on u.id=r.peer_id 
   		WHERE r.id=?
	`
	err := db.Get(&info, queryDB, id)

	if err != nil {
		fmt.Print("err:", err)
		return nil, err
	}
	fmt.Print(info)
	return &info, nil

}

func QueryListAccountWithAndNumberReport(db *sqlx.DB, phone string) (*[]ResultAccountNumberReport, error) {

	result := []ResultAccountNumberReport{}
	query := `
	SELECT  u1.first_name, u1.last_name, u1.phone,COALESCE(COUNT(uc.peer_id), 0) as totalReport
	FROM users u1 
	LEFT JOIN reports uc on u1.id=uc.peer_id 
	Where u1.phone = ? AND user_type != 9
	GROUP BY u1.id 
   `
	err := db.Select(&result, query, phone)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func QueryListLinkByPeerId(db *sqlx.DB, peer_id int) (*[]LinkStruct, error) {
	result := []LinkStruct{}
	query := `select un.id,un.link from username un where un.peer_id =? `
	err := db.Select(&result, query, peer_id)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func QueryListReportByPeerId(db *sqlx.DB, peer_id int) (*[]ReportPeerIdStruct, error) {
	result := []ReportPeerIdStruct{}
	query := `select r.id, u.first_name, u.last_name,u.phone,r.reason ,r.content,r.created_at  from reports r join users u on u.id= r.user_id join users u2 on u2.id =r.peer_id where r.peer_id =? `
	err := db.Select(&result, query, peer_id)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func QueryListReportByReporter(db *sqlx.DB, reporter int, orderBy string) (*[]ReportPeerIdStruct, error) {
	result := []ReportPeerIdStruct{}
	query := `select r.id, u.first_name, u.last_name,u.phone,r.reason ,r.content,r.created_at  from reports r join users u on u.id= r.user_id join users u2 on u2.id =r.peer_id where r.user_id =? order by created_at ` + orderBy
	err := db.Select(&result, query, reporter)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func FilterReportOfReporterByReason(db *sqlx.DB, reason int, reporter int, ch chan ResponseFilterOwnerByReason, wg *sync.WaitGroup) {
	defer wg.Done()

	var responseData ResponseFilterOwnerByReason
	dataResult := []FilterOwnerByReason{}

	query := `select r.id, u.first_name, u.last_name,u.phone,r.reason ,r.content,r.created_at  from reports r join users u on u.id= r.user_id join users u2 on u2.id =r.peer_id where r.user_id =? and r.reason= ? order by created_at DESC `

	err := db.Select(&dataResult, query, reporter, reason)

	if err != nil {
		fmt.Print("query" + err.Error())
		responseData.Error = true
		ch <- responseData
		return
	}

	responseData.Error = false
	responseData.Data = dataResult

	ch <- responseData

}

// Ban user

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
		user.Ban = 1
		_, err := db.NamedExec("UPDATE ban_user SET ban = :ban  WHERE uid = :uid", user)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

}

func BannedOrtherAccountWithPhone(db *sqlx.DB, peer_id int, ban int, phone string, reason int, admin_ban string) (int, error) {

	uid := []UidBanned{}
	err := db.Select(&uid, "SELECT id FROM users WHERE user_type= -1 AND phone= ?", phone)
	if err != nil {
		return 1, nil
	} else {

		for _, user_id := range uid {

			var user CheckBan
			err := db.Get(&user, "SELECT id, uid FROM ban_user WHERE  uid = ?", user_id.Id)
			if err != nil {
				tx := db.MustBegin()
				tx.MustExec("INSERT INTO ban_user(uid,ban) VALUES(?,?)", user_id.Id, ban)
				tx.MustExec("INSERT INTO ban_history(uid,ban,reason,admin_ban) VALUES(?,?,?,?)", user_id.Id, ban, reason, admin_ban)
				if err := tx.Commit(); err != nil {
					return 0, err

				}
				return 1, nil
			} else {
				user.Ban = 1
				_, err := db.NamedExec("UPDATE ban_user SET ban = :ban  WHERE uid = :uid", user)
				if err != nil {
					return 0, err

				}
				return 1, nil
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
		fmt.Print("query" + err.Error())
		return nil, err
	}
	var count float64
	error2 := db.QueryRow(countQuery, ban).Scan(&count)
	if error2 != nil {
		return nil, error2
	}
	responseUser.Data = usersBan
	responseUser.Limit = limit
	responseUser.Page = page
	responseUser.TotalPage = int(count)
	return &responseUser, nil

}

func GetHistoryBan(db *sqlx.DB, id int) (*ResponseHistoryban, error) {
	var responseBan ResponseHistoryban
	result := responseBan.Data

	query := `select id,ban,reason,admin_ban,created_at from ban_history where uid=?`

	err := db.Select(&result, query, id)
	if err != nil {

		return nil, err
	}
	responseBan.Data = result

	return &responseBan, nil
}

func FilterOwnerByTypeReason(db *sqlx.DB, reason int, peer_id int, ch chan ResponseFilterOwnerByReason, wg *sync.WaitGroup) {
	defer wg.Done()

	var responseData ResponseFilterOwnerByReason
	dataResult := []FilterOwnerByReason{}
	query := `select u.first_name, u.last_name,u.phone,r.reason ,r.content,r.created_at  from reports r join users u on u.id= r.user_id join users u2 on u2.id =r.peer_id where r.reason=? and r.peer_id =? `

	err := db.Select(&dataResult, query, reason, peer_id)

	if err != nil {
		fmt.Print("query" + err.Error())
		responseData.Error = true
		ch <- responseData
		return
	}

	responseData.Error = false
	responseData.Data = dataResult

	ch <- responseData

}

func DetailReportByPeerId(db *sqlx.DB, id int) (*ReportPeerIdStruct, error) {
	result := ReportPeerIdStruct{}

	query := ` select r.id, u.first_name, u.last_name,u.phone,r.reason ,r.content,r.created_at from reports r join users u on u.id=r.user_id where r.id=? `

	err := db.Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

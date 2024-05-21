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
	TotalReported string `json:"total_reported" db:"total"`
	Time          string `json:"created_at" db:"created_at"`
}

type RecordDetail struct {
	FirstName    string `json:"first_name"  `
	LastName     string `json:"last_name" `
	Phone        string `json:"phone"`
	TotalLink    int    `json:"totalLink" `
	TotalAccount int    `json:"totalAccount" `
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

func QueryListReport(db *sqlx.DB, page int, limit int, ch chan Response, wg *sync.WaitGroup) {
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
    SELECT r.id, u2.first_name, u2.last_name, u2.phone, r.created_at, uc.total
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
    LIMIT ? OFFSET ?`

	err := db.Select(&dataResult, query, limit, offset)

	if err != nil {

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

func InsertReport(db *sqlx.DB, reporterId int, reportedId int, message string, reason string, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Print(reportedId, reportedId, message, reason)
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO reports(reporter_id, report_id, message, reason) VALUES(?, ?, ?, ?)", reporterId, reportedId, message, reason)
	if err := tx.Commit(); err != nil {

		ch <- false
		return
	}
	ch <- true

}

func DeleteReport(db *sqlx.DB, id int, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	row, err := db.Query("delete from report where id=?", id)
	if err != nil {
		ch <- false
		return
	}
	fmt.Print(row)
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
	Where NOT r.peer_id=0 AND r.reason =?
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

func SearchReportByPhone(db *sqlx.DB, phone string, limit int, page int, ch chan Response, wg *sync.WaitGroup) {
	defer wg.Done()

	offset := (page - 1) * limit

	var responseData Response

	dataResult := make([]Record, limit)

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
    Where  u2.phone=? OR u1.phone= ? AND NOT r.peer_id=0 LIMIT ? offset ? `

	err := db.Select(&dataResult, querySearch, phone, phone, limit, offset)

	if err != nil {

		responseData.Err = true
		ch <- responseData
		return
	}

	var count float64

	countQuery := `SELECT count(rp.id) as totalPage from reports rp JOIN users u1 ON rp.user_id=u1.id join users u2 on rp.peer_id =u2.id where u1.phone =? or u2.phone=?`

	error2 := db.QueryRow(countQuery, phone, phone).Scan(&count)

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
	WHERE   u1.phone !='' AND u1.phone = ?
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


func QueryListAccountWithAndNumberReport(db *sqlx.DB,phone string,){
	query := `
	WITH uid_counts AS (
        SELECT peer_id, COUNT(*) AS total
         FROM reports
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
            WHERE peer_id != 0
            GROUP BY peer_id
        ))
    ) r
    JOIN uid_counts uc ON r.peer_id = uc.peer_id
    JOIN users u1 ON r.user_id = u1.id
    JOIN users u2 ON r.peer_id = u2.id
   `
   fmt.Print(query)
}
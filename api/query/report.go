package query

import (
	"fmt"
	"math"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Record struct {
	Id            int    `json:"id" db:"id"`
	FirstName  string `json:"firstName" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Phone string `json:"phone" db:"phone"`
	TotalReported string `json:"total_reported" db:"total"`
	Time          string `json:"created_at" db:"created_at"`
}

type RecordDetail struct {
	Id                int    `json:"id" db:"id"`
	FirstNameReporter string `json:"firstname_rpter" db:"firstname_rpter"`
	LastNameReporter  string `json:"lastname_rpter" db:"lastname_rpter"`
	EmailReporter     string `json:"email_rpter" db:"email_rpter"`
	PhoneReporter     string `json:"phone_rpter" db:"phone_rpter"`

	FirstNameReported string `json:"firstname_rpted" db:"firstname_rpted"`
	LastNameReported  string `json:"lastname_rpted" db:"lastname_rpted"`
	EmailReported     string `json:"email_rpted" db:"email_rpted"`
	ReportedPhone     string `json:"phone_rpted" db:"phone_rpted"`

	Message string `json:"message" db:"message"`
	Reason  string `json:"reason" db:"reason"`
	Time    string `json:"created_at" db:"created_at"`
}
type ResponseDetail struct {
	Data RecordDetail `json:"data"`
	Err  bool         `json:"error"`
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
    SELECT r.id, u2.first_name,u2.last_name, u2.phone,r.created_at, uc.total
    FROM reports r
    JOIN uid_counts uc ON r.peer_id = uc.peer_id
    JOIN users u1 ON r.user_id=u1.id join users u2 on r.peer_id =u2.id
    Where NOT r.peer_id=0 LIMIT ? offset ?`

	countQuery := `SELECT count(rp.id) as totalPage from reports rp JOIN users u1 ON rp.user_id=u1.id join users u2 on rp.peer_id =u2.id`
	err := db.Select(&dataResult, query, limit, offset)

	if err != nil {
 
		responseData.Err = true
		ch <- responseData
		return
	}
	var count float64
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

func QueryDetailReport(db *sqlx.DB, id int, ch chan ResponseDetail, wg *sync.WaitGroup) {
	defer wg.Done()

	var responseData ResponseDetail
	dataResult := RecordDetail{}

	query := `SELECT report.id,report.created_at,report.content,report.reason,u1.last_name AS lastname_rpter,u1.first_name AS firstname_rpter,u1.phone AS phone_rpter,u2.phone AS phone_rpted ,u3.email AS email_rpter,u4.email AS email_rpted, u2.first_name AS firstname_rpted, u2.last_name AS lastname_rpted
	  FROM report rp
	  JOIN user u1 ON rp.user_id=u1.id
	  JOIN user u2 ON  rp.peer_id =u2.id
	  JOIN username u3 ON rp.user_id =u3.peer_id
	  JOIN username u4 ON rp.peer_id =u4.peer_id
	  WHERE report.id=?`

	err := db.Get(&dataResult, query, id)
	if err != nil {
		fmt.Print("eror ne:" + err.Error())
		responseData.Err = true
		ch <- responseData
		return
	}

	responseData.Err = false
	responseData.Data = dataResult

	ch <- responseData
}

func InsertReport(db *sqlx.DB, reporterId int, reportedId int, message string, reason string, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Print(reportedId, reportedId, message, reason)
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO report(reporter_id, report_id, message, reason) VALUES(?, ?, ?, ?)", reporterId, reportedId, message, reason)
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

func FilterByTypeReason (db *sqlx.DB, reason int,page int, limit int, ch chan Response, wg *sync.WaitGroup) {
	defer wg.Done()
	offset := (page - 1) * limit
	var responseData Response
	dataResult := make([]Record, limit)
	query := `
	WITH uid_counts AS (
        SELECT peer_id, COUNT(*) AS total
        FROM reports r
        Where r.reason = ?
        GROUP BY peer_id
    )
    SELECT r.id, u2.first_name,u2.last_name, u2.phone,r.created_at, uc.total
    FROM reports r
    JOIN uid_counts uc ON r.peer_id = uc.peer_id
    JOIN users u1 ON r.user_id=u1.id join users u2 on r.peer_id =u2.id
    Where NOT r.peer_id=0 AND r.reason =? LIMIT ? offset ?
	`

	countQuery := `SELECT count(rp.id) as totalPage from reports rp JOIN users u1 ON rp.user_id=u1.id join users u2 on rp.peer_id =u2.id where rp.reason =?`
	err := db.Select(&dataResult, query,reason,reason, limit, offset)

	if err != nil {
		fmt.Print("query"+err.Error())
		responseData.Err = true
		ch <- responseData
		return
	}
	var count float64
	error2 := db.QueryRow(countQuery,reason).Scan(&count)
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


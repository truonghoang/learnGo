package query

import (
	"fmt"
	"math"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Record struct {
	Id            int    `json:"id" db:"id"`
	NameReporter  string `json:"name_reporter" db:"name_reporter"`
	ReportedName  string `json:"reported_name" db:"reported_name"`
	PhoneReporter string `json:"phone_reporter" db:"phone_reporter"`
	ReportedPhone string `json:"phone_reported" db:"phone_reported"`
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
	query := `SELECT  report.id,report.created_at,u1.first_name AS name_reporter,u2.first_name AS reported_name,u1.phone AS phone_reporter,u2.phone AS phone_reported 
	FROM report 
	JOIN user u1 ON report.reporter_id=u1.id 
	JOIN user u2 ON report.report_id =u2.id 
	ORDER BY id LIMIT ? OFFSET ?`
	countQuery := `SELECT count(id) as totalPage from report`
	err := db.Select(&dataResult, query, limit, offset)

	if err != nil {
		responseData.Err = true
		ch <- responseData
		return
	}

	var count int
	error2 := db.QueryRow(countQuery).Scan(&count)

	if error2 != nil {
		responseData.Err = true
		ch <- responseData
		return
	}
	totalPage:=math.Ceil(float64(count/limit))
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

	query := `SELECT report.id,report.created_at,report.message,report.reason,u1.last_name AS lastname_rpter,u1.first_name AS firstname_rpter,u1.phone AS phone_rpter,u2.phone AS phone_rpted ,u3.email AS email_rpter,u4.email AS email_rpted, u2.first_name AS firstname_rpted, u2.last_name AS lastname_rpted
	  FROM report
	  JOIN user u1 ON report.reporter_id=u1.id
	  JOIN user u2 ON  report.report_id =u2.id
	  JOIN user_name u3 ON report.reporter_id =u3.uid
	  JOIN user_name u4 ON report.report_id =u4.uid
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

func DeleteReport(db *sqlx.DB,id int, ch chan bool,wg *sync.WaitGroup){
	defer wg.Done()
	row,err:= db.Query("delete from report where id=?",id)
	if err!= nil {
		ch <- false
		return
	}
	fmt.Print(row)
	ch <- true


}

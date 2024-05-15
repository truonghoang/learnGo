package connection

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:truong652000@tcp(localhost:3306)/scamreport")
	if err != nil {
		return nil, err
	}

	return db, nil

}

package connection

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:truong27913211@tcp(localhost:3306)/hugmanh")
	if err != nil {
		return nil, err
	}

	return db,nil

}

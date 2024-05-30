package connection

import (
	"fmt"
	"os"
	"truonghoang/go-scam/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

 
func ConnectDb() (*sqlx.DB, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadConfig(pwd + "/config/connectdb.json")
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:3306)/%s`,cfg.User,cfg.Password,cfg.Host,cfg.DB_Name)
	
	db, err := sqlx.Connect("mysql",dsn)
	if err != nil {
		return nil, err
	}

	return db, nil

}

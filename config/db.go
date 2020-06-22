package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func init() {
	var (
		host     = os.Getenv("db_host")
		port, _  = strconv.Atoi(os.Getenv("db_port"))
		user     = os.Getenv("db_user")
		password = os.Getenv("db_pass")
		dbname   = os.Getenv("db_name")
	)

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open(os.Getenv("db_driver"), connString)

	if err != nil {
		fmt.Println(err)
		db = nil
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err)
		db = nil
	}
}

func GetDB() *sql.DB {
	return db
}

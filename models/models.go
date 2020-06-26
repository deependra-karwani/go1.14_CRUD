package models

import (
	"CRUD/config"
	"fmt"
)

var db = config.GetDB()

func CreateTables() {
	var err error

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, profPic VARCHAR(255), email VARCHAR(254) NOT NULL UNIQUE, mobile CHAR(10), username VARCHAR(50) NOT NULL UNIQUE, password CHAR(64) NOT NULL, fcm TEXT, token VARCHAR(255) UNIQUE)")
	if err != nil {
		fmt.Println("Users Table Error: ", err)
	} else {
		fmt.Println("Users Table Created if did not Exist")
	}
}

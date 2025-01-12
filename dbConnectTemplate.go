package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func connectToDB(name string) (*sql.DB,error) {

	config := mysql.Config{
		User: "[INSERTDBUSERHERE]",
		Passwd: "[INSERTPASSWORDHERE]",
		DBName: "dbProjInfo",
	}

	db, _ := sql.Open("mysql",config.FormatDSN())

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db,nil
}

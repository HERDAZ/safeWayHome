package main

import (
	"log"
	"fmt"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// connect to db
	db, err := connectToDB("dbProjInfo")
	if err != nil {log.Fatal(err)}

	fmt.Println("Connected")

	// create fake user
	var pos Position
	pos.userID = "ABCD"
	pos.time = "2022-02-16 15:03:12"
	pos.latitude = 12.8967081
	pos.longitude = 24.5683885

	//push data to db
	err = pushPositionToDB(db, pos.userID, pos.time, pos.latitude, pos.longitude)
	if err != nil {log.Printf("Could not push to DB : ", err)}

	//retrive data from db
	var rows *sql.Rows

	rows, err = getRowsFromTable(db,"coords")
	if err != nil {log.Fatal("Could not retrieve table from DB : ", err)}

	positions, _ := extractPositions(rows)
	fmt.Printf("%v",positions)
}

package main

import (
	"log"
	"fmt"
	"time"
)

func main() {
	// connect to db
	db, err := connectToDB("dbProjInfo")
	if err != nil {log.Fatal(err)}

	fmt.Println("Connected")

	// create fake user
	var pos PositionDB
	pos.UserID = "DEF"
	pos.Time = time.Now().Format(time.DateTime)
	pos.Latitude = 12.8967081
	pos.Longitude = 24.5683885

	//push data to db
	//err = pushPositionToDB(db, pos.UserID, pos.Time, pos.Latitude, pos.Longitude)
	//if err != nil {log.Printf("Could not push to DB : ", err)}

	//retrive data from db
	var positions []PositionDB

	positions, err = getUserPosition(db,"ABCD",true)
	if err != nil {log.Fatal("Could not retrieve position from DB : ", err)}

	fmt.Printf("%v\n",positions)

	positions, err = getRowsFromTable(db, "coords")
	if err != nil { log.Fatal("Could not retrieve table from DB : ", err) }

	fmt.Printf("%v",positions)
}


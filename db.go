package main

import (
	"errors"
	"log"
	"fmt"
	"database/sql"
)


func extractHomes(rows *sql.Rows) ([]HomeDB, int) {

	var homes []HomeDB
	var home HomeDB
	
	count := 0 // rows counter
	for rows.Next() {
		rows.Scan(&home.Time, &home.UserID, &home.Latitude, &home.Longitude)
		homes = append(homes, home)
		count++
	}

	return homes, count

}

func extractPositions(rows *sql.Rows) ([]PositionDB, int) {

	var positions []PositionDB
	var pos PositionDB
	
	count := 0 // rows counter
	for rows.Next() {
		rows.Scan(&pos.Time, &pos.UserID, &pos.Latitude, &pos.Longitude)
		positions = append(positions, pos)
		count++
	}

	return positions, count
}

func getUsersRows(db *sql.DB, userName string, tableName string, latest bool) (*sql.Rows, error) {
	var query string

	if latest { // only get the last position
		query = fmt.Sprintf("SELECT * FROM %s WHERE userID = '%s' ORDER BY -time LIMIT 1;", tableName, userName)
	} else { // get all positions
		query = fmt.Sprintf("SELECT * FROM %s WHERE userID = '%s' ORDER BY -time", tableName, userName)
	}

	var rows *sql.Rows
	var err error

	rows, err = db.Query(query)
	if err != nil { return nil, err }

	return rows, nil
}


func getUsersPosition(db *sql.DB, apikey string, userID string, latest bool) ([]PositionDB, error) {

	if len(apikey) != 32 { 

		err := fmt.Sprintf("Invalid API key '%s'", apikey)
		log.Printf("WARNING : %s", err)
		return []PositionDB{}, errors.New(err)
	}

	//TODO add permissions

	rows, err := getUsersRows(db, userID, "coords", latest)
	if err != nil { log.Printf("WARNING : Problem getting position for user '%s' : ", err) }

	positions, count := extractPositions(rows)
	if count == 0 { log.Printf("WARNING : No rows in table 'coords' for user '%s'\n", userID) }

	return positions, nil
}

func getUsersHome(db *sql.DB, apikey string, userID string) (HomeDB, error) {

	if len(apikey) != 32 { 

		err := fmt.Sprintf("Invalid API key '%s'", apikey)
		log.Printf("WARNING : %s", err)
		return HomeDB{}, errors.New(err)
	}

	rows, err := getUsersRows(db, userID, "home", false)
	if err != nil { log.Printf("WARNING : Problem getting position for user '%s' : ", err) }

	homes, count := extractHomes(rows)

	if count == 0 {

		log.Printf("WARNING : No rows in table 'home' for user '%s''\n", userID)
		return HomeDB{}, nil

	} else if count > 1 {

		log.Printf("WARNING : Too many homes found for user '%s' (%d)\n", userID, count)
	}

	return homes[0], nil
}

func pushHomeToDB(db *sql.DB, userID string, time string, latitude float64, longitude float64) error {

	var query string

	err := deleteHome(db, userID)
	if err != nil { log.Printf("WARNING : Can't delete home position for user '%s' %v\n", userID, err) }

	query = "INSERT INTO home (userID, time, latitude, longitude) VALUES "
	query += fmt.Sprintf("('%s','%s',%.7f,%.7f);", userID, time, latitude, longitude)

	_, err = db.Exec(query)
	if err != nil { return err }

	return nil
}

func pushPositionToDB(db * sql.DB, apikey string, time string, latitude float64, longitude float64) error {

	var query string

	userID, err := getUserFromAPIkey(db, apikey)
	if err != nil {
		log.Println("ERROR : Problem with APIkey to userID conversion :", err)
		return nil
	}

	query = "INSERT INTO coords (userID, time, latitude, longitude) VALUES " // static value
	query += fmt.Sprintf("('%s','%s',%.7f,%.7f);", userID, time, latitude, longitude) // dynamic values

	_, err = db.Exec(query)
	//TODO add more tourough testing
	if err != nil {
		return err
	}

	return nil
}

func deletePositions(db *sql.DB, userID string, timeUpTo string) error {
	delQuery := fmt.Sprintf("DELETE FROM coords WHERE userID = \"%s\" WHERE time <= \"%s\";", userID, timeUpTo)

	result, err := db.Exec(delQuery)
	if err != nil { return err }

	count, _ := result.RowsAffected()
	if count == 0 { log.Println("WARNING : No rows to delete") }

	return nil
}

func deleteHome(db *sql.DB, userID string) error {
	delQuery := fmt.Sprintf("DELETE FROM home WHERE userID = \"%s\";", userID)

	result, err := db.Exec(delQuery)
	if err != nil { return err }

	count, _ := result.RowsAffected()
	if count > 1 { log.Printf("WARNING : Multiple home position deleted for user '%s'\n", userID) }

	return nil
}

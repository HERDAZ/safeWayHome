package main

import (
	"log"
	"fmt"
	"database/sql"
)

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

func getRowsFromTable(db *sql.DB, tableName string) ([]PositionDB, error) {

	rows, err := db.Query("SELECT * FROM " + tableName);

	if err != nil {
		return nil, err
	}

	positions, count := extractPositions(rows)
	if count == 0 { log.Printf("WARNING : No rows in table '%s'\n", tableName) }

	return positions, nil
}

func getUserPosition(db *sql.DB, userName string, latest bool) ([]PositionDB, error) {

	var query string

	if latest { // only get the last position
		query = fmt.Sprintf("SELECT * FROM coords WHERE userID = '%s' ORDER BY time LIMIT 1;", userName)
	} else { // get all positions
		query = fmt.Sprintf("SELECT * FROM coords WHERE userID = '%s'", userName)
	}

	var rows *sql.Rows
	var err error

	rows, err = db.Query(query)
	if err != nil { return nil, err }

	positions, count := extractPositions(rows)
	if count == 0 { log.Printf("WARNING : No postion found for user '%s'\n", userName) }

	return positions, nil
}


func pushPositionToDB(db * sql.DB, userID string, time string, latitude float64, longitude float64) error {

	var query string

	query = "INSERT INTO coords (userID, time, latitude, longitude) VALUES " // static value
	query += fmt.Sprintf("('%s','%s',%.7f,%.7f);", userID, time, latitude, longitude) // dynamic values (line is split for visibility)

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func deletePosition(db *sql.DB, userID string) error {
	delQuery := fmt.Sprintf("DELETE FROM coords WHERE userID = \"%s\";", userID)

	result, err := db.Exec(delQuery)
	if err != nil { return err }
	count, _ := result.RowsAffected()
	if count != 1 { log.Println("WARNING : No rows to delete") }

	return nil
}

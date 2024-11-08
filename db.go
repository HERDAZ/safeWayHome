package main

import (
	"log"
	"fmt"
	"database/sql"
	
	"github.com/go-sql-driver/mysql"
)

func getRowsFromTable(db *sql.DB, table string) (*sql.Rows, error) {

	rows, err := db.Query("SELECT * FROM " + table);

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func deletePosition(db *sql.DB, userID string) error {
	delQuery := fmt.Sprintf("DELETE FROM coords WHERE userID = \"%s\";", userID)

	result, err := db.Exec(delQuery)
	if err != nil { return err }
	count, _ := result.RowsAffected()
	if count != 1 { log.Println("WARNING : No rows to delete") }

	return nil
}

func extractPositions(rows *sql.Rows) ([]Position, int) {
	var positions []Position
	var pos Position
	
	count := 0 // rows counter
	for rows.Next() {
		rows.Scan(&pos.Time, &pos.UserID, &pos.Latitude, &pos.Longitude)
		positions = append(positions, pos)
		count++
	}

	return positions, count
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



package main

import (
	"fmt"
	"database/sql"
	"log"
	"errors"
)

func getPermissions(db *sql.DB, userID string, friendID string) (Permissions, error) {

	query := fmt.Sprintf("SELECT seePosition, sendMessage FROM relations WHERE userID = '%s' AND friendID = '%s';", userID, friendID)

	var seePos []byte
	var sendMsg []byte

	row := db.QueryRow(query)
	err := row.Scan(&seePos, &sendMsg)

	if err == sql.ErrNoRows {
		errorMsg := fmt.Sprintf("ERROR : Couldn't get permissions for users '%s' and '%s' : %v", userID, friendID, err)
		return Permissions{}, errors.New(errorMsg)
	}

	var perms Permissions

	perms.seePosition = (seePos[0] == 0x1)
	perms.sendMessage = (sendMsg[0] == 0x1)

	return perms, nil
}

func validateNewRelation(db *sql.DB, userID string, friendID string) error {

	if userID == friendID {
		errorMsg := fmt.Sprintf("ERROR : Creating a looped-relation, where '%s' is friend with him/her/themselves", userID)
		return errors.New(errorMsg)
	}
	
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM relations WHERE userID = '%s' AND friendID = '%s';", userID, friendID)

	row := db.QueryRow(query)
	err := row.Scan(&count)

	if err == sql.ErrNoRows {
		errorMsg := fmt.Sprintf("ERROR : Couldn't validate new relation for users '%s' and '%s : %v", userID, friendID, err)
		return errors.New(errorMsg)
	}
	
	if count != 0 {
		errorMsg := fmt.Sprintf("WARNING : Relation between userID '%s' and userID '%s' already exist", userID, friendID)
		return errors.New(errorMsg)
	}
	
	return nil
}

func updatePermission(db *sql.DB, userID string, friendID string, perm string, status int) error {
	
	query := fmt.Sprintf("UPDATE relations SET %s = b'%b' WHERE userID = '%s' AND friendID = '%s';", perm, status, userID, friendID)

	result, err := db.Exec(query)
	if err != nil {
		errorMsg := fmt.Sprintf("ERROR := Cannot update '%s' to '%d' for userID '%s' and friendID '%s' : %v", perm, status, userID, friendID, err)
		return errors.New(errorMsg)
	}

	count, _ := result.RowsAffected()
	if count > 1 {
		log.Printf("WARNING : Multiple relations updated for userID '%s' and friendID '%s' (%d)\n", userID, count)
	}

	return nil
}

func createRelation(db *sql.DB, userID string, friendID string, addDate string) error {

	query := fmt.Sprint("INSERT INTO relations (userID, friendID, addDate) ")
	query += fmt.Sprintf("VALUES ('%s', '%s', '%s');", userID, friendID, addDate)

	_, err := db.Exec(query)
	if err != nil { return err }

	return nil
}

//func main() {
//	//db, _ := connectToDB("dbProjInfo")
//
//	//userID := "ABCD"
//	//friendID := "abcd"
//	//now := time.Now().Add(time.Hour).Format(time.DateTime)
//	//permType := "seePosition"
//	//status := 1
//
//	//err := createRelation(db, userID, friendID, now)
//	//err := updatePermission(db, userID, friendID, permType, status)
//	//err := validateNewRelation(db, userID, friendID)
//	//fmt.Println(err)
//	//perms, err := getPermissions(db, userID, friendID)
//	//fmt.Println(perms, err)
//}

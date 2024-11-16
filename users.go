package main

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"errors"
)

func pushUserToDB(db *sql.DB, userID string, lastLogin string, username string, email string, phoneNb string, passwdHash string) error {
	
	var query string

	query = "INSERT INTO users (userID, lastLogin, username, phoneNb, email, passwdHash) VALUES "
	query += fmt.Sprintf("('%s','%s','%s', '%s','%s','%s');", userID, lastLogin, username, phoneNb, email, passwdHash)

	_, err := db.Exec(query)
	if err != nil { return err }

	return nil
}

func generateUserID(db *sql.DB) string  {

	var userID string;

	for i := 0; i<4; i++ {

		max := big.NewInt(26+26+10)
		symbol_, _ := rand.Int(rand.Reader, max)
		symbol := symbol_.Int64()

		if symbol < 26 {
			//lowercase
			userID += fmt.Sprintf("%s", string(97+symbol))
		} else if symbol < 52 {
			//uppercase
			userID += fmt.Sprintf("%s", string(65+symbol-26))
		} else {
			//number
			userID += fmt.Sprintf("%s", string(48+symbol-52))
		}
	}
	
	//chech if userID is already taken
	query := fmt.Sprintf("SELECT * FROM users WHERE userID = '%s' LIMIT 1;", userID)

	row := db.QueryRow(query)
	err := row.Scan()

	if err != sql.ErrNoRows {
		return generateUserID(db)
	}
	
	return userID;
}


func validateNewUser(db *sql.DB, username string, email string, phoneNb string) error {

	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s' OR phoneNb = '%s' OR userName= '%s' LIMIT 1;", email, phoneNb, username)

	row := db.QueryRow(query)
	err := row.Scan()

	if err != sql.ErrNoRows {
		log.Printf("WARNING : User with email '%s' or phone number '%s' or username '%s' already exist\n", email, phoneNb, username)
		return errors.New("User already exist (phone number or email already associated with an user")
	}

	return nil
}

func pushNewUserToDB(db *sql.DB, lastLogin string, username string, email string, phoneNb string, passwd string) (string, error) {

	err := validateNewUser(db, username, email, phoneNb)
	if err != nil {
		return "", err
	}

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), 10)
	if err != nil {
		log.Printf("ALERT : Problem when hashing password '%s' %v\n", passwd, err)
		return "", err
	}

	userID := generateUserID(db)

	err = pushUserToDB(db, userID, lastLogin, username, email, phoneNb, string(passwdHash))
	if err != nil {
		log.Printf("WARNING : Could not insert user in DB\n")
		return "", err
	}

	return userID, nil
}

func deleteUser(db *sql.DB, userID string) error {
	delQuery := fmt.Sprintf("DELETE FROM users WHERE userID = '%s';", userID)

	result, err := db.Exec(delQuery)
	if err != nil { return err }

	count, _ := result.RowsAffected()
	if count == 0 {
		log.Printf("WARNING : No user to delete for userID '%s \n", userID)
	} else if count > 1 {
		log.Printf("WARNING : Multiple users deleted for userID '%s' (%d)\n", userID, count)
	}

	return nil
}

func generateAPIkey(db *sql.DB, count int) (string, error)  {

	if count > 10 {
		return "", errors.New("ERROR : 10 bad iteration of API key generation, recursion stoped")
	}

	var APIkey string;

	for i := 0; i<32; i++ {

		max := big.NewInt(26+26+10)
		symbol_, _ := rand.Int(rand.Reader, max)
		symbol := symbol_.Int64()

		if symbol < 26 {
			//lowercase
			APIkey += fmt.Sprintf("%s", string(97+symbol))
		} else if symbol < 52 {
			//uppercase
			APIkey += fmt.Sprintf("%s", string(65+symbol-26))
		} else {
			//number
			APIkey += fmt.Sprintf("%s", string(48+symbol-52))
		}
	}
	
	//chech if this API key is already taken
	query := fmt.Sprintf("SELECT * FROM users WHERE APIkey = '%s' LIMIT 1;", APIkey)

	row := db.QueryRow(query)
	err := row.Scan()

	if err != sql.ErrNoRows {
		return generateAPIkey(db, count+1)
	}
	
	return APIkey, nil;
}

func authenticateUser(db *sql.DB, username string, password string) (string, string, error) {
	
	var passwordHash string
	var userID string

	query := fmt.Sprintf("SELECT passwdHash, userID FROM users WHERE username = '%s';", username)

	row := db.QueryRow(query)
	err := row.Scan(&passwordHash, &userID)

	if err == sql.ErrNoRows {
		errorMsg := fmt.Sprintf("WARNING : No password hash for user '%s'\n", username)
		return "", "", errors.New(errorMsg)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		errorMsg := fmt.Sprintf("Bad password (%s) for user '%s'", password, username)
		return "", "", errors.New(errorMsg)
	}

	APIkey, err := generateAPIkey(db,0)

	return APIkey, userID, err
}

func updateAPIkey(db *sql.DB, userID string, APIkey string) error {
	query := fmt.Sprintf("UPDATE users SET APIkey = '%s' WHERE userID = '%s';", APIkey, userID)

	result, err := db.Exec(query)
	if err != nil { return err }

	count, _ := result.RowsAffected()

	if count == 0 {
		errorMsg := fmt.Sprintf("WARNING : No APIkeys updated, either no user with userID '%s', or same APIkey\n", userID)
		return errors.New(errorMsg)
	} else if count > 1 {
		errorMsg := fmt.Sprintln("WARNING : Multiple APIkeys updated for userID '%s' (%d)\n", userID, count)
		return errors.New(errorMsg)
	}

	return nil
}

func getUserFromAPIkey(db *sql.DB, APIkey string) (string, error) {

	query := fmt.Sprintf("SELECT userID FROM users WHERE APIkey = '%s';", APIkey)

	rows, err := db.Query(query)
	if err != nil {
		errorMsg := fmt.Sprintf("ERROR : Couldn't get userID from DB : ", err)
		return "", errors.New(errorMsg)
	}

	var IDs []string
	var userID string

	count := 0
	for rows.Next() {
		rows.Scan(&userID)
		IDs = append(IDs, userID)
		count += 1
	}
	
	if count > 1 {
		errorMsg := fmt.Sprintf("WARNING : Multiple user %v have the same API key '%s'", IDs, APIkey)
		return IDs[0], errors.New(errorMsg)
	}

	return userID, nil
}

//func main() {
//	db, _ := connectToDB("dbProjInfo")
//
//	username := "test"
//	password := "test"
//
//	APIkey, userID, err := authenticateUser(db, username, password)
//	fmt.Println(APIkey, userID, err)
//
//	err = updateAPIkey(db, userID, APIkey)
//	fmt.Println(err)
//
//	userID, err = getUserFromAPIkey(db, APIkey)
//	fmt.Println(userID, err)
//}

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"database/sql"
	"log"
	"fmt"
	"time"	
)

var db *sql.DB

func main() {

	var err error
	db, err = connectToDB("dbProjInfo")
	if err != nil { log.Fatal("Could not connect to database : ", err) }

	router := gin.Default()

	router.GET("/position", getposition)
	router.GET("/home", getHomeposition)
	router.GET("/login", login)

	router.POST("/home", postHomeposition)	
	router.POST("/position", postPosition)
	router.POST("/signup", signup)

	router.Run("87.106.79.94:8447")
}

func postPosition(c *gin.Context) {

	now := time.Now().Add(time.Hour).Format(time.DateTime)

	var newPosition PositionRequest

	err := c.BindJSON(&newPosition)
	fmt.Printf("%v\n", newPosition) //TODO WTF does latitude stays 0 ?!
	
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create new position"})
		fmt.Println(err)
		return
	}

	err = pushPositionToDB(db, newPosition.APIkey, now, newPosition.Latitude, newPosition.Longitude)
        if err != nil { log.Println("Could not push data to database : ", err) }


	c.IndentedJSON(http.StatusCreated, newPosition)
}

func getposition(c *gin.Context) {
     
	apikey := c.GetHeader("apikey")
	userID := c.GetHeader("userID")

        fmt.Println("API key :", apikey)
        fmt.Println("userID :", userID)

	positions, err := getUsersPosition(db, apikey, userID, false)
	if err != nil { log.Println("Could not retrieve data from table : ", err) }

	c.IndentedJSON(http.StatusCreated, positions[0]) //TODO fix error code 500
	

}

func getHomeposition(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	userID := c.GetHeader("userID")

        fmt.Println("API key :", apikey)
        fmt.Println("userID :", userID)

	homePosition, err := getUsersHome(db, apikey, userID)
	if err != nil { log.Println("WARNING : Could not retrieve home position for API key '%s'", apikey) }

	c.IndentedJSON(http.StatusCreated, homePosition)
}

func postHomeposition(c *gin.Context) {
	now := time.Now().Add(time.Hour).Format(time.DateTime)

	var newHomePosition HomeRequest
        if err := c.BindJSON(&newHomePosition); err != nil {
                c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to parse home position"})
                fmt.Println(err)
                return
        }

	userID, err := getUserFromAPIkey(db, newHomePosition.APIkey)
	if err != nil { fmt.Println("Couldn't get userID from API key :", err) }

        err = pushHomeToDB(db, userID, now, newHomePosition.Latitude, newHomePosition.Longitude)
        if err != nil { log.Println("Could not push data to database : ", err) }

        c.IndentedJSON(http.StatusCreated, newHomePosition)
	
}

func signup(c *gin.Context) {

	var user UserSignup
	c.BindJSON(&user)

	fmt.Println("user :", user)

	err := validateNewUser(db, user.UserName, user.Email, user.PhoneNb)

	if err != nil {
		response := fmt.Sprintf("Problem validating new user : %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
		
	userID, err := pushNewUserToDB(db, time.Now().Add(time.Hour).Format(time.DateTime), user.UserName, user.Email, user.PhoneNb, user.Password)

	if err == nil {
		response := fmt.Sprintf("New user '%s' with userID '%s' was created\n", user.UserName, userID)
		c.IndentedJSON(http.StatusCreated, response)
	} else {
		response := fmt.Sprintf("Problem creating user with username '%s'\n", user.UserName, userID)
		log.Printf(response)
		c.IndentedJSON(http.StatusInternalServerError, response)
	}
}

func login(c *gin.Context) {

	username := c.GetHeader("username")
	password := c.GetHeader("password")

	userKey, userID, err := authenticateUser(db, username, password)

	if err == nil {

		response := fmt.Sprintf("{'userKey' : '%s', 'userID' : '%s'}", userKey, userID)
		c.IndentedJSON(http.StatusCreated, response)

	} else {

		response := fmt.Sprintf("{'error' : '%s'}", err)
		c.IndentedJSON(http.StatusUnauthorized, response)
	}

}

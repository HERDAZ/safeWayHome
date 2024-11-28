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

	router.GET("/position", getPosition)
	router.GET("/home", getHomeposition)
	router.GET("/login", login)

	router.POST("/home", postHomeposition)	
	router.POST("/position", postPosition)
	router.POST("/signup", signup)

	router.Run("87.106.79.94:8447")
}

func makeErrMsg(err error) Error {
	response := Error{ErrorMsg : err.Error()}
	return response
}

func postPosition(c *gin.Context) {

	now := time.Now().Add(time.Hour).Format(time.DateTime)

	var newPosition PositionRequest

	err := c.BindJSON(&newPosition)
	
	if err != nil {

		log.Println("ERROR : Couldn't bind newPosition from JSON :", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = pushPositionToDB(db, newPosition.APIkey, now, newPosition.Latitude, newPosition.Longitude)
        if err != nil {

		response := makeErrMsg(err)
		log.Println("postPosition :", err.Error())
		c.IndentedJSON(http.StatusUnauthorized, response)

	} else {

		c.IndentedJSON(http.StatusCreated, newPosition)
	}

}

func getPosition(c *gin.Context) {
     
	//TODO find why the fuck there is an error 500 that is not this one
	apikey := c.GetHeader("apikey")
	userID := c.GetHeader("userID")

        fmt.Println("API key :", apikey)
        fmt.Println("userID :", userID)

	positions, err := getUsersPosition(db, apikey, userID, false)
	if err != nil {
		log.Println("ERROR : Could not retrieve data from table : ", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusInternalServerError, response)
	} 

	c.IndentedJSON(http.StatusCreated, positions[0])
	

}

func getHomeposition(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	userID := c.GetHeader("userID")

        fmt.Println("API key :", apikey)
        fmt.Println("userID :", userID)

	homePosition, err := getUsersHome(db, apikey, userID)
	if err != nil {

		log.Printf("WARNING : Could not retrieve home position for API key '%s'\n", apikey)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusCreated, response)
		return
	}

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
	if err != nil {
		fmt.Println("Couldn't get userID from API key :", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

        err = pushHomeToDB(db, userID, now, newHomePosition.Latitude, newHomePosition.Longitude)
        if err != nil {
		log.Println("Could not push data to database : ", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

        c.IndentedJSON(http.StatusCreated, newHomePosition)
	
}

func signup(c *gin.Context) {

	var user UserSignup
	c.BindJSON(&user)
	//TODO : add err catch for bind

	err := validateNewUser(db, user.UserName, user.Email, user.PhoneNb)

	if err != nil {
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
		
	userID, err := pushNewUserToDB(db, time.Now().Add(time.Hour).Format(time.DateTime), user.UserName, user.Email, user.PhoneNb, user.Password)

	if err != nil {
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusCreated, SignupResponse{ UserID: userID })
}

func login(c *gin.Context) {

	username := c.GetHeader("username")
	password := c.GetHeader("password")

	userKey, userID, err := authenticateUser(db, username, password)
	
	if err != nil {
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusUnauthorized, response)
	} else {

		response := LoginResponse{Apikey: userKey, Userid : userID}
		c.IndentedJSON(http.StatusCreated, response)
	}

}

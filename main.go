package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"database/sql"
	"log"
	"fmt"
	"time"	
	"errors"
	"slices"
)

var db *sql.DB
var isHome []string

func main() {

	var err error
	db, err = connectToDB("dbProjInfo")
	if err != nil { log.Fatal("Could not connect to database : ", err) }

	router := gin.Default()

	router.GET("/position", getPosition)
	router.GET("/home", getHomeposition)
	router.GET("/login", getLogin)

	router.POST("/home", postHomeposition)	
	router.POST("/position", postPosition)
	router.POST("/signup", postSignup)

	router.POST("/amHome", postAmHome)
	router.GET("/isHome", getIsHome)
	//router.POST("/cleanIsHome", postCleanIsHome) //TODO delete for production

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
	log.Printf("postPosition with %v", newPosition)
	
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
     
	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("getPosition with apikey : %s, friendID : %s\n", apikey, friendID)

	if apikey == "" || friendID == "" {
		errorMsg := makeErrMsg(errors.New("Empty apikey or friendID"))
		c.IndentedJSON(http.StatusBadRequest, errorMsg)
		return
	}

	userID, err := getUserFromAPIkey(db, apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}

	perms, err := getPermissions(db, userID, friendID)

	if perms.seePosition != true {

		log.Printf("WARNING : Insuficient perm seePosition for userID '%s' and friendID '%s'\n", userID, friendID)
		response := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusUnauthorized, response)
		return
	}

	positions, err := getUsersPosition(db, friendID, false)

	if err != nil {
		log.Println("ERROR : Could not retrieve data from table : ", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	} 

	c.IndentedJSON(http.StatusCreated, positions[0])

}

func getHomeposition(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("getHomePosition with apikey : %s, friendID : %s\n", apikey, friendID)

	userID, err := getUserFromAPIkey(db, apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(errorMsg)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}

	perms, err := getPermissions(db, userID, friendID)

	if perms.seePosition != true {

		log.Printf("WARNING : Insuficient perms for userID '%s' and friendID '%s' : ", userID, friendID)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusUnauthorized, response)
		return
	}

	homePosition, err := getUsersHome(db, friendID)
	if err != nil {

		log.Printf("WARNING : Could not retrieve home position for friendID '%s'\n", friendID)
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
		//TODO redo
		log.Println("Failed to parse home position")
                c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to parse home position"})
                return
        }
	log.Printf("postHomePosition with %v\n", newHomePosition)

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

func postSignup(c *gin.Context) {

	var user UserSignup
	c.BindJSON(&user)
	log.Printf("postSignup with %v\n", user)
	//TODO : add err catch for bind

	err := validateNewUser(db, user.UserName, user.Email, user.PhoneNb)

	if err != nil {
		response := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
		
	userID, err := pushNewUserToDB(db, time.Now().Add(time.Hour).Format(time.DateTime), user.UserName, user.Email, user.PhoneNb, user.Password)

	if err != nil {
		response := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	c.IndentedJSON(http.StatusCreated, SignupResponse{ UserID: userID })
}

func getLogin(c *gin.Context) {

	username := c.GetHeader("username")
	password := c.GetHeader("password")
	log.Printf("getLogin with username : %s, password : %s\n", username, password)

	userKey, userID, err := authenticateUser(db, username, password)
	
	if err != nil {
		log.Println(err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusUnauthorized, response)
	} else {
		log.Println(err)
		response := LoginResponse{Apikey: userKey, Userid : userID}
		c.IndentedJSON(http.StatusCreated, response)
	}

}

//func postCleanIsHome(c *gin.Context) {
//	isHome = []string{}
//	fmt.Println("isHome reset : ", isHome)
//}

func getIsHome(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("getIsHome with apikey : '%s', friendID : '%s'\n", apikey, friendID)

	userID, err := getUserFromAPIkey(db, apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}

	perms, err := getPermissions(db, userID, friendID)

	if perms.sendMessage != true {

		log.Printf("WARNING : Insuficient perm sendMessage for userID '%s' and friendID '%s'\n", userID, friendID)
		c.IndentedJSON(http.StatusUnauthorized, "GTFO, Ur not supposed to be here")
		return
	}

	if slices.Contains(isHome, friendID) { c.IndentedJSON(http.StatusOK, true) 
	} else { c.IndentedJSON(http.StatusOK, false) }
}

func postAmHome(c *gin.Context) {

	var apikey Apikey
	err := c.BindJSON(&apikey)
	log.Printf("postAmHome with apikey : '%s'\n", apikey)

	if err != nil {
		log.Println(err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userID, err := getUserFromAPIkey(db, apikey.Apikey)

	isHome = append(isHome, userID)
}

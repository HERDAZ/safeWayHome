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
var inEmergency []string
var isStoped []string

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

	// router.POST("/amHome", postAmHome) DEPRECATED, NOW IN POST /position
	//router.GET("/isHome", getIsHome) DEPRECATED, NOW HANDLED BY GET /getFriend
	router.POST("/cleanIsHome", postCleanIsHome) //TODO delete for production

	router.POST("/inEmergency", postInEmergency)
	router.GET("/inEmergency", getInEmergency)
	router.POST("/cleanInEmergency", postCleanIsHome) //TODO delete for production

	router.GET("/isStoped", getIsStoped)
	router.POST("/cleanIsStoped", postCleanIsStoped) //TODO delete for production

	router.POST("/addFriend", postAddFriend)
	router.GET("/getFriend", getFriendList)

	router.Run("87.106.79.94:8447")
}

func makeErrMsg(err error) Error {
	response := Error{ErrorMsg : err.Error()}
	return response
}

func deleteElement(slice []string, elem string) []string { //should create a util class at this point
	for i:=0; i<len(slice);i++ {
		if slice[i] == elem {
			slice[i] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}

func getFriendList(c * gin.Context) {

	apikey := c.GetHeader("apikey")
	log.Printf("INFO : getFriendList with apikey : %s\n", apikey)

	if apikey == "" {
		log.Println("Empty apikey")
		c.IndentedJSON(http.StatusBadRequest, "Empty apikey")
		return
	}

	userID, err := getUserFromAPIkey(db, apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}
	
	var relations []Relation
	relations, err = getUsersRelations(db, userID)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return
	}

	c.IndentedJSON(http.StatusOK, relations)
}

func postAddFriend(c * gin.Context) {

	now := time.Now().Add(time.Hour).Format(time.DateTime)
	
	var friendRequest AddFriend

	err := c.BindJSON(&friendRequest)

	log.Println("postAddFriend with ", friendRequest);

	if err != nil {
		log.Printf("ERROR : Couldn't bind friendRequest from JSON :", err);
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	err = createFriendship(db, friendRequest.APIkey, friendRequest.FriendID, now, true, true)
	if err != nil {
		errorMsg := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, errorMsg)
		return
	}
	c.IndentedJSON(http.StatusOK, "")

}

func postPosition(c *gin.Context) {

	now := time.Now().Add(time.Hour).Format(time.DateTime)

	var newPosition PositionRequest

	err := c.BindJSON(&newPosition)
	log.Printf("INFO : postPosition with %v", newPosition)
	
	if err != nil {

		log.Println("ERROR : Couldn't bind newPosition from JSON :", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}


	
	userID, err := getUserFromAPIkey(db, newPosition.APIkey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println("ERROR : Can't get userID of apikey '", newPosition.APIkey, "' :", err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return
	}

	if newPosition.IsHome && !slices.Contains(isHome, userID) {
		isHome = append(isHome, userID)
	} else if !newPosition.IsHome && slices.Contains(isHome, userID) {
		isHome = deleteElement(isHome, userID)
	}

	delta := 0.00002 // yes i know, no magic numbers, but for now idgas (it's the minimum step, if 2 latitude or longitude are within, they are considered the same)
	if !slices.Contains(isStoped, userID) {

		//check if user is at the same place as 2 minutes ago (and not already in isStoped)
		timeMinus := time.Now().Add(time.Hour-2*time.Minute-15*time.Second).Format(time.DateTime)
		timePlus  := time.Now().Add(time.Hour-2*time.Minute+15*time.Second).Format(time.DateTime)

		// TODO maybe refactor with COUNT(*)
		query := 	     "SELECT latitude, longitude FROM coords WHERE "
		query += fmt.Sprintf("userID = '%s' AND ", userID)
		query += fmt.Sprintf("time BETWEEN '%s' AND '%s' AND ", timeMinus, timePlus)
		query += fmt.Sprintf("latitude BETWEEN %.7f AND %.7f AND ", newPosition.Latitude-delta, newPosition.Latitude+delta)
		query += fmt.Sprintf("longitude BETWEEN %.7f AND %.7f ", newPosition.Longitude-delta, newPosition.Longitude+delta)
		query +=             "LIMIT 1;"

		fmt.Println("INFO : SQL query :", query)
		fmt.Println("INFO : Time :", now)

		row := db.QueryRow(query)
		err := row.Scan()

		if err != sql.ErrNoRows { //there is a problem : same position (within delta) than 2 minutes ago

			userID, err := getUserFromAPIkey(db, newPosition.APIkey)
			if err != nil {
				fmt.Println("ERROR : Couldn't get userID from API key :", err)
				errorMsg := makeErrMsg(err)
				c.IndentedJSON(http.StatusBadRequest, errorMsg)
				return
			}

			fmt.Printf("ALERT : UserID '%s' is stoped\n", userID)
			isStoped = append(isStoped, userID)
		}

	} else {

		//check if user started moving again (and in isStoped)
		query := fmt.Sprintf("SELECT latitude, longitude FROM coords WHERE userID = '%s' ORDER BY time DESC LIMIT 1;", userID)

		var lastPosition Position
		row := db.QueryRow(query)
		err := row.Scan(&lastPosition)

		if err != nil {
			errorMsg := makeErrMsg(err)
			fmt.Printf("ERROR : Couldn't scan into Position : %s\n", err)
			c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		}

		latDiff := lastPosition.Latitude - newPosition.Latitude
		lonDiff := lastPosition.Longitude - newPosition.Longitude

		if !( (-delta < latDiff && latDiff > delta) && (-delta < lonDiff && lonDiff > delta) ) { //if not still in the same place
			fmt.Printf("ALERT : UserID '%s' is moving again\n", userID)
			isStoped = deleteElement(isStoped, userID)
		}
	}

	//now we push the data to db
	err = pushPositionToDB(db, newPosition.APIkey, now, newPosition.Latitude, newPosition.Longitude)

        if err != nil {
		response := makeErrMsg(err)
		log.Println("ERROR : postPosition :", err)
		c.IndentedJSON(http.StatusInternalServerError, response)
		return 
	}

	c.IndentedJSON(http.StatusCreated, newPosition)
}

func getPosition(c *gin.Context) {
     
	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("INFO : getPosition with apikey : %s, friendID : %s\n", apikey, friendID)

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
	log.Printf("INFO : getHomePosition with apikey : %s, friendID : %s\n", apikey, friendID)

	userID, err := getUserFromAPIkey(db, apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(errorMsg)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}

	if userID != friendID {

		perms, err := getPermissions(db, userID, friendID)

		if perms.seePosition != true {
			log.Printf("WARNING : Insuficient perms for userID '%s' and friendID '%s' : ", userID, friendID)
			response := makeErrMsg(err)
			c.IndentedJSON(http.StatusUnauthorized, response)
			return
		}

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

	err := c.BindJSON(&newHomePosition)
        if err != nil {
		//TODO redo
		log.Println("ERROR : Failed to parse home position")
                c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to parse home position"})
                return
        }
	log.Printf("INFO : postHomePosition with %v\n", newHomePosition)

	userID, err := getUserFromAPIkey(db, newHomePosition.APIkey)
	if err != nil {
		fmt.Println("ERROR : Couldn't get userID from API key :", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

        err = pushHomeToDB(db, userID, now, newHomePosition.Latitude, newHomePosition.Longitude)
        if err != nil {
		log.Println("ERROR : Could not push data to database : ", err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

        c.IndentedJSON(http.StatusCreated, newHomePosition)
	
}

func postSignup(c *gin.Context) {

	var user UserSignup
	c.BindJSON(&user)
	log.Printf("INFO : postSignup with %v\n", user)
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
	log.Printf("INFO : getLogin with username : %s, password : %s\n", username, password)

	userKey, userID, err := authenticateUser(db, username, password)
	
	if err != nil {
		log.Println(err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusUnauthorized, response)
		return
	} else {
		log.Println(err)
		response := LoginResponse{Apikey: userKey, Userid : userID}
		c.IndentedJSON(http.StatusCreated, response)
		return
	}

}

func getIsHome(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("INFO : getIsHome with apikey : '%s', friendID : '%s'\n", apikey, friendID)

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


func postCleanIsHome(c *gin.Context) {
	isHome = []string{}
	fmt.Println("isHome reset : ", isHome)
}

func postInEmergency(c *gin.Context) {

	var apikey Apikey
	err := c.BindJSON(&apikey)
	log.Printf("INFO : postInEmergency with apikey : '%s'\n", apikey)

	if err != nil {
		log.Println(err)
		response := makeErrMsg(err)
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userID, err := getUserFromAPIkey(db, apikey.Apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}

	inEmergency = append(inEmergency, userID)
	c.IndentedJSON(http.StatusNoContent, nil)
}

func getInEmergency(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("INFO : getInEmergency with apikey : '%s', friendID : '%s'\n", apikey, friendID)

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

	if slices.Contains(inEmergency, friendID) { c.IndentedJSON(http.StatusOK, true) 
	} else { c.IndentedJSON(http.StatusOK, false) }
}

func postCleanInEmergency(c *gin.Context) {
	inEmergency = []string{}
	fmt.Println("inEmergency reset : ", inEmergency)
}

func getIsStoped(c *gin.Context) {

	apikey := c.GetHeader("apikey")
	friendID := c.GetHeader("friendID")
	log.Printf("INFO : getIsStoped with apikey : '%s', friendID : '%s'\n", apikey, friendID)

	userID, err := getUserFromAPIkey(db, apikey)

	if err != nil {
		errorMsg := makeErrMsg(err)
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, errorMsg)
		return 
	}

	perms, err := getPermissions(db, userID, friendID)

	if perms.seePosition != true {

		log.Printf("WARNING : Insuficient perm sendMessage for userID '%s' and friendID '%s'\n", userID, friendID)
		c.IndentedJSON(http.StatusUnauthorized, "GTFO, Ur not supposed to be here")
		return
	}

	if slices.Contains(isStoped, friendID) { c.IndentedJSON(http.StatusOK, true) 
	} else { c.IndentedJSON(http.StatusOK, false) }
}

func postCleanIsStoped(c *gin.Context) {
	isStoped = []string{}
	fmt.Println("isStoped reset : ", isStoped)
}

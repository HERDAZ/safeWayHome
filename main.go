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
	router.POST("/position", createposition)
	
	router.Run("87.106.79.94:8447")
}

func createposition(c *gin.Context) {

	now := time.Now()
        fmt.Println("YYYY.MM.DD : ", now.Format("2006.01.02 15:04:05"))




	var newPosition PositionRequest
	if err := c.BindJSON(&newPosition); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create new position"})
		fmt.Println(err)
		return
	}

	err := pushPositionToDB(db, newPosition.UserID, now.Format("2006.01.02 15:04:05"), newPosition.Latitude, newPosition.Longitude)
        if err != nil { log.Println("Could not push data to database : ", err) }


	c.IndentedJSON(http.StatusCreated, newPosition)
}

func getposition(c *gin.Context) {

	positions, err := getUsersPosition(db, "TEUB", false) // regarde la docu, la syntaxe à (oui encore) changée
	if err != nil {log.Fatal("Could not retrieve data from table : ", err) }

	c.IndentedJSON(http.StatusCreated, positions[0])
	

}

func getHomeposition(c *gin.Context) {

	homePosition, err := getUsersHome(db, "ABCD")
	if err != nil { log.Println("WARNING : Could not retrieve home position for user 'ABCD'") }

	c.IndentedJSON(http.StatusCreated, homePosition)
}

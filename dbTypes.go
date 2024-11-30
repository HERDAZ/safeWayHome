package main

type datetime struct {
	year 	   int8
	month 	   int8
	day 	   int8
	hour 	   int8
	minute 	   int8
}

type PositionRequest struct {
        APIkey     string  `json:"apikey"`
        Latitude   float64 `json:"latitude"`
        Longitude  float64 `json:"longitude"`
}

type PositionDB struct {
        UserID     string  `json:"userid"`
	Time       string  `json:"time"` // format "YY-MM-DD HH:MM:SS"
        Latitude   float64 `json:"latitude"`
        Longitude  float64 `json:"longitude"`
}

type HomeRequest struct {
 	APIkey 	    string  `json:"apikey"`
 	Latitude    float64 `json:"latitude"`
 	Longitude   float64 `json:"longitude"`
}

type HomeDB struct {
 	UserID 	    string  `json:"userid"`
	Time	    string  `json:"time"`
 	Latitude    float64 `json:"latitude"`
 	Longitude   float64 `json:"longitude"`
}

 type Relation struct {
 	UserID 	    string `json:"userid"`
	FriendID    string `json:"friendid"`
	Permissions byte   `json:"permissions"`
	AddDate     string `json:"addDate"`
}

type User struct {
	LastLogin    string `json:"lastLogin"`
	PhoneNb      string `json:"phoneNb"`
	Email 	     string `json:"email"`
	PasswdHash   string `json:"passwdHash"`
	UserID       string `json:"userid"`
}

type UserLogin struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
}

type LoginResponse struct {
        Apikey     string `json:"apikey"`
        Userid     string `json:"userid"`
}

type UserSignup struct {
	UserName     string `json:"username"`
	PhoneNb      string `json:"phoneNb"`
	Email 	     string `json:"email"`
	Password     string `json:"password"`
}

type SignupResponse struct {
	UserID	     string `json:"userid"`
	//TODO add APIkey maybe
}

 type UserAlerts struct {
	id 	     int     
 	UserID 	     string  `json:"userid"`
	Time 	     float32 `json:"time"`
	AlertType    string  `json:"alertType"`
 }

 type AlertNotif struct {
 	alertID	     int 	 
	sentTo	     string 
 }

 type Permissions struct {
 	seePosition  bool `json:"seeposition"`
	sendMessage  bool `json:"sendmessage"`
 }

 type Apikey struct {
 	Apikey      string `json:"apikey"`
}

 type Error struct {
 	ErrorMsg    string `json:"error"`
 }

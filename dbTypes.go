package main

type datetime struct {
	year 	   int8
	month 	   int8
	day 	   int8
	hour 	   int8
	minute 	   int8
}

type Position struct {
        UserID     string  `json:"userID"`
        Latitude   float64 `json:"latitude"`
        Longitude  float64 `json:"longitude"`


}


 type Relation struct {
 	UserID 	    string `json:"userID"`
	FriendID    string `json:"friendID"`
	Permissions byte   `json:"permissions"`
	AddDate     string `json:"addDate"`
}

type User struct {
	LastLogin    string `json:"lastLogin"`
	PhoneNb      string `json:"phoneNb"`
	Email 	     string `json:"email"`
	PasswdHash   string `json:"passwdHash"`
	UserID       string `json:"userID"`
}
 type UserAlerts struct {
	id 	     int     
 	UserID 	     string  `json:"userID"`
	Time 	     float32 `json:"time"`
	AlertType    string  `json:"alertType"`
 }

 type AlertNotif struct {
 	alertID	  int 	 
	sentTo	  string 
 }

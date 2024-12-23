package main

type CustomError struct {
 	errType ErrorType
	Stack	[]string
}


 type ErrorType string

 const (
 	BadAPIkey = 1
	BadUserID = 2
	BadFriendID = 3
	FailedJSONParsing = 4
)

var errorName = map[ErrorType]string{
 	BadAPIkey: "Bad API key"
	BadUserID: "Bad UserID"
	BadFriendID: "Bad FriendID"
	FailedJSONParsing: "Failed JSON parsing"
}

func 

func (err CustomError) AddToStack(newError string) {
	//todo
}

func main(){}

//I just want to preface that i dont know whatever the fuck i'm doing here, and any comment i make in this file is for helping my 5min-in-the-future-self find out whatever fuckery of an idea i had trying to implement this
package main

import (
	"fmt"
)

//The ErrorType definition, and it's String interface implementation
type ErrorType int
const (
 	BadAPIkey ErrorType = 1
	BadUserID ErrorType = 2
	BadFriendID ErrorType = 3
	FailedJSONParsing ErrorType = 4
	FailedSQLQuery ErrorType = 5
)
var errorName = map[ErrorType]string{
 	BadAPIkey: "Bad API key",
	BadUserID: "Bad UserID",
	BadFriendID: "Bad FriendID",
	FailedJSONParsing: "Failed JSON parsing",
	FailedSQLQuery: "Failed internal SQL query",
}
func (err ErrorType) String() string {
	return errorName[err]
}


//The ErrorStack, a wraper around ErrorType with the implementation of a stacktrace
type ErrorStack struct {
	Origin ErrorType
	Stack []string
}

//It's initialisation function
func NewErrorStack(errType ErrorType, stack string) ErrorStack {
	var err ErrorStack
	err.Origin = errType
	err.AddToStack(stack)
	fmt.Println(err.Stack)
	return err
}

//it's String interface implementation
func (err ErrorStack) String() string {
	log := fmt.Sprintf("Origin : %s\n", err.Origin)
	log += fmt.Sprintf("Stacktrace :\n")
	fmt.Println("AAAA", len(err.Stack))
	for i := 0; i<len(err.Stack); i++ {
		fmt.Println(i)
		log += fmt.Sprintf("	C%s%d\n", err.Stack[i], i)
	}
	return log
}

//to add to the stack
func (err ErrorStack) AddToStack(log string) {
	err.Stack = append(err.Stack, log)
}

func main() {
	err := NewErrorStack(BadAPIkey, "APIkey = AAAAAAAAAAA")
	fmt.Println(err)
}

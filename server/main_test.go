package main

import (
	"fmt"
	"testing"
)

// FullName string   `json:"name"`
// UserName string   `json:"username"`
// Password string   `json:"password"`
// Schedule Schedule `json:"schedule"`
// ClassTime string `json:"classtime"`
// DayOfWeek string `json:"weekday"`
// Date      string `json:"date"`
// Frequency string `json:"frequency"`
func TestRand(t *testing.T) {
	fmt.Println("in the test")
	u := User{FullName: "Alex", UserName: "alex@test.com", Password: "thisisasecret", Schedule: Schedule{ClassTime: "9:00am", DayOfWeek: "Monday", Date: "09/09/2020", Frequency: "9"}}
	u.CalculateSignUpTimes()
}

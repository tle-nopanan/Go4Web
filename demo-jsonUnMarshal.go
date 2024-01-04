package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee1 struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {

	e := employee1{}
	err := json.Unmarshal([]byte(`{"ID":101,"EmployeeName":"Tle Tle","Tel":"090000000","Email":"Tle@mail.com"}`), &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.EmployeeName)
}

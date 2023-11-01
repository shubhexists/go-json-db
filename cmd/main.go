package main

import (
	"encoding/json"
	"fmt"
	. "github.com/shubhexists/go-json-db/models" 
)
//EVERY THING HERE SHOULD BE IN MODELS FOR THE USER TO CREATE THESE DYNAMICALLY, THESE ARE JUST THE EXAMPLES.
//ALSO THINK OF A BASIC WAY TO ALLOW USERS TO CREATE THESE STRUCTS DYNAMICALLY
//THIS MAY BE ADDED INTO THE EXAMPLES FOLDER (CREATE LATER MAYBE)\

type User struct{
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}

type Address struct{
	City string
	State string
	Country string
	Pincode json.Number
}

func main(){
	dir := "./"
	db, err := New(dir, nil)
	if err != nil{
		fmt.Println(err)
		return
	}
	employees := []User{
		{"John","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
		{"John1","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
		{"John2","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
		{"John3","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
	}

	for _,items := range employees {
		db.Write("users", items.Name, User{
			Name: items.Name,
			Age: items.Age,
			Contact: items.Contact,
			Company: items.Company,
			Address: items.Address,
		})
	}

	records, err := db.ReadAll("users");
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _, f := range records{
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound);
		err != nil {
				fmt.Println("Error", err)
			}
		allusers = append(allusers, employeeFound)
	}
	fmt.Println((allusers))
}

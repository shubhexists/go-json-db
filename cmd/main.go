package main

import (
	"encoding/json"
	"fmt"
	. "github.com/shubhexists/go-json-db/models" 
)
//EVERY THING HERE SHOULD BE IN MODELS FOR THE USER TO CREATE THESE DYNAMICALLY, THESE ARE JUST THE EXAMPLES.
//ALSO THINK OF A BASIC WAY TO ALLOW USERS TO CREATE THESE STRUCTS DYNAMICALLY
//THIS MAY BE ADDED INTO THE EXAMPLES FOLDER (CREATE LATER MAYBE)


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
		{"John4","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
		{"John10","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
		{"John20","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
		{"John30","23","9585394030","Humanize",Address{
			"Delhi",
			"Delhi",
			"India",
			"110092",
		}},
	}
	// Move all the databases in a seperate folder to make it more clean, else multiple collections will create multiple folders polluting the code
	for _,items := range employees {
		db.Write("trial", items.Name, User{
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
}

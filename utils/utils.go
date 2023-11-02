package utils

import (
	"os"
	"reflect"
	"fmt"
)

//Utility Function to verify if a specific file exists or not..
func Stat(path string)(fi os.FileInfo, err error){
	if fi, err = os.Stat(path);
	os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return
}

//Utility fucntion to check for tag "db" with value main in the struct and returning the struct member name



// Utility function to check if datatype is Struct and if it is a struct, expand it

func ExpandStruct(s interface{}) {
	v := reflect.ValueOf(s)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if fieldType.Type.Kind() == reflect.Struct {
			fmt.Printf("Expanding struct field: %s\n", fieldType.Name)
			ExpandStruct(field.Interface())
		} else {
			fmt.Printf("Field: %s, Value: %v\n", fieldType.Name, field.Interface())
		}
	}
}
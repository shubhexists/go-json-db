package utils

import (
	"fmt"
	"os"
	"reflect"
)

// Utility Function to verify if a specific file exists or not..
func Stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

// Utility fucntion to check for tag "db" with value main in the struct and returning the struct member name
// Also Check if it is Unique, If not Return an Error - TODO
func CheckTag(s interface{}) (string, error) {
	v := reflect.ValueOf(s)
	t := v.Type()

	mainFieldFound := false
	mainField := ""

	for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if tagValue, ok := field.Tag.Lookup("db"); ok {
            if tagValue == "main" {
				if mainFieldFound {
					fmt.Println("Error: Multiple main fields found")
					return "", fmt.Errorf("multiple main fields found")
				} else {
					mainFieldFound = true
					mainField = v.Field(i).Interface().(string)
				}
            }
        }
    }
	if mainFieldFound {
		return mainField, nil
	} else {
		fmt.Println("Error: No main field found")
		return "", fmt.Errorf("no main field found")
	}
}

// Utility function to check if datatype is Struct and if it is a struct, expand it

func ExpandStruct(s interface{}) {
	v := reflect.ValueOf(s)
	t := v.Type()
	fmt.Println(t)

	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		fmt.Println(t.Kind())
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

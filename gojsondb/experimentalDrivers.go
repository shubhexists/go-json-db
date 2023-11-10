package gojsondb

import (
	"encoding/json"
	"fmt"
	"github.com/jcelliott/lumber"
	"github.com/shubhexists/go-json-db/utils"
	"os"
	"path/filepath"
)

// Update with only the required fields, TO COMPLETE
func (driver *Driver) Update(collection string, data string, v interface{}, newValues map[string]interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection - Unable To Update")
	}
	if data == "" {
		return fmt.Errorf("missing record - Unable To Update")
	}

	record := filepath.Join(driver.dir, collection, data)
	if _, err := utils.Stat(record); err != nil {
		return err
	}

	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	jsonData := make(map[string]interface{})
	err = json.Unmarshal(b, &jsonData)
	if err != nil {
		return err
	}

	//THIS CODE HAS ERROR SOMEWHERE - REST IS CORRECT

	// for key, value := range newValues{
	// 	if _, exists := jsonData[key]; exists {
	// 		jsonData[key] = value
	// 	} else {

	// 	}
	// }

	b, err = json.MarshalIndent(jsonData, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := os.WriteFile(record+".json", b, 0644); err != nil {
		return err
	}

	return nil
}

// TO CHECK IF IT WORKS IN NESTED STRUCTS(JSON) - Add to experimental maybe?
// ALSO THIS VALUE IS JUST PRINTED NOT RETURNED ( TO FIX)
// ADD Search by other values except primary key
func (driver *Driver) Search(collection string, searchField string, value string) error {
	if collection == "" {
		lumber.Error("Missing collection - Unable To Search")
		return fmt.Errorf("missing collection - Unable To Search")
	}
	if searchField == "" {
		lumber.Error("Please enter a valid Search Field")
		return fmt.Errorf("missing search field - Unable To Search")
	}
	if value == "" {
		lumber.Error("Please enter a valid value")
		return fmt.Errorf("missing value - Unable To Search")
	}
	dir := filepath.Join(driver.dir, collection)
	if _, err := utils.Stat(dir); err != nil {  
		return err
	}

	files, _ := os.ReadDir(dir)
	var records []string

	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			lumber.Error("Some error Occurred")
			return err
		}
		if err2 := json.Unmarshal(b, &records); err2 != nil {
			lumber.Error("Error while handling JSON Data, Kindly confirm if the JSON data is correct")
			return err2
		}
	}
	return nil
}

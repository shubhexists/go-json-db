package models

import (
	"encoding/json"
	"fmt"
	"github.com/shubhexists/go-json-db/utils"
	"io/ioutil"
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

	b, err := ioutil.ReadFile(record + ".json")
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
	if err := ioutil.WriteFile(record+".json", b, 0644); err != nil {
		return err
	}

	return nil
}

// TO CHECK IF IT WORKS IN NESTED STRUCTS(JSON) - Add to experimental maybe?
// ALSO THIS VALUE IS JUST PRINTED NOT RETURNED ( TO FIX)
// ADD Search by other values except primary key
func (driver *Driver) Search(collection string, data string, key string) error {
	return nil
}

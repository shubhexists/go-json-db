package models

import (
	"sync"
	// io/util is deprecated but I found no alternative on ChatGPT :), To be updated soon with newer one
	"io/ioutil"
	//Maybe later on shift to a faster library for json encoding and decoding
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/patrickmn/go-cache"
	"github.com/shubhexists/go-json-db/utils"
)

type (
	Driver struct {
		// db operations should be non synchronous
		// Got a reference from this - https://youtu.be/jkRN9zcLH1s?si=s5ec23U3tS5bi6EO
		mutex sync.Mutex
		// this map will be used to store mutexes for each collection
		mutexes map[string]*sync.Mutex
		dir     string
	}
)

/*
Current operations supported are -
1) Write
2) Read
3) ReadAll
4) Delete
5) Delete Collection
6) Update Record
*/

// CREATE A NEW DATABASE (COLLECTION)
func New(dir string) (*Driver, *cache.Cache, error) {
	//This checks for any incorrect filename and corrects it.
	dir = filepath.Clean(dir)

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
	}

	if _, err := os.Stat(dir); err == nil {
		lumber.Info("Database already exists")
		return &driver, StartCache(5, 10), nil
	}

	lumber.Info("Creating Database in directory %s", dir)
	return &driver, StartCache(5, 10), os.MkdirAll(
		dir,
		0755)
}

// MANAGE MUTEXES FOR EACH COLLECTION
func (driver *Driver) ManageMutex(collection string) *sync.Mutex {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()
	m, ok := driver.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		driver.mutexes[collection] = m
	}
	return m
}

// WRITE ANY RECORD TO A GIVEN COLLECTION
func (driver *Driver) Write(collection string, v interface{}) error {
	if collection == "" {
		lumber.Error("Missing collection - No place to save record!")
		return fmt.Errorf("missing collection - no place to save record")
	}

	data, err := utils.CheckTag(v)
	if err != nil {
		return err
	}

	mutex := driver.ManageMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, collection)
	tempdir := filepath.Join(dir, data)
	fnlPath := filepath.Join(dir, data + ".json")
	tmpPath := fnlPath + ".tmp"

	// Check if the file already exists
	if _, err := utils.Stat(tempdir); err == nil {
		fmt.Println("Record already exists!")
		return fmt.Errorf("record already exists")
	}
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

// READ ANY RECORD FROM A GIVEN COLLECTION
// Only From Primary Key
func (driver *Driver) Read(collection string, data string,c *cache.Cache, wantCache bool) (string, error) {
	if collection == "" {
		lumber.Error("Missing collection - No place to save record!")
		return "", fmt.Errorf("missing collection - Unable To Read")
	}

	if data == "" {
		lumber.Error("Missing data - No place to save record!")
		return "", fmt.Errorf("missing data - Unable To Read")
	}

	record := filepath.Join(driver.dir, collection, data)
	if _, err := utils.Stat(record); err != nil {
		return "", err
	}

	if(wantCache){
		if records, found := GetCache(c, record); found {
			lumber.Info("Fetching data from cache")
			return records.(string), nil
		}
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return "", err
	}

	if(wantCache){
		lumber.Info("Saved data to Cache! ")
		SetCache(c , record , string(b))
	}
	
	return string(b), nil
}

// READ ALL RECORDS FROM A GIVEN COLLECTION
// THIS WILL RETURN JSON ARRAY OF ALL THE RECORDS
func (driver *Driver) ReadAll(collection string, c *cache.Cache, wantCache bool) ([]string, error) {
	if collection == "" {
		lumber.Error("Missing collection - No place to save record!")
		return nil, fmt.Errorf("missing collection - Unable to Read Record")
	}
	dir := filepath.Join(driver.dir, collection)
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	if(wantCache){
		if records, found := GetCache(c, dir); found {
			lumber.Info("Fetching data from cache")
			return records.([]string), nil
		}
	}
	files, _ := ioutil.ReadDir(dir)
	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}
	if(wantCache){
		lumber.Info("Saved data to Cache")
		SetCache(c, dir, records)
	}

	return records, nil
}

// DELETE ANY RECORD FROM A GIVEN COLLECTION
func (driver *Driver) Delete(collection string, data string) error {
	if collection == "" {
		lumber.Error("Missing collection - Unable to Delete!")
		return fmt.Errorf("missing collection - Unable To Delete")
	}

	if data == "" {
		lumber.Error("Cannot Delete - Record Not Found")
		return fmt.Errorf("please enter a valid record")
	}

	path := filepath.Join(collection, data)
	mutex := driver.ManageMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, path)
	switch fi, err := utils.Stat(dir); {
	case fi == nil, err != nil:
		lumber.Error("Cannot Delete - Record %v Not Found", path)
		return fmt.Errorf("unable to find file or directory named %v", path)
	case fi.Mode().IsDir():
		lumber.Error("Cannot Delete - %v is a collection", path)
		return fmt.Errorf("this seems like a collection, kindly enter a record to delete or use db.DeleteCollection")
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

// Delete any collection (All records in that collection)
func (driver *Driver) DeleteCollection(collection string) error {
	path := filepath.Join(collection)
	mutex := driver.ManageMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, path)
	switch fi, err := utils.Stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("not a valid collection. Kindly enter a valid collection Name")
	default:
		return os.RemoveAll(dir)
	}
}

// Update any record from a given collection
// Currently we have to enter the entire User struct, UPDATE IT SO THAT WE CAN UPDATE ONLY THE REQUIRED FIELDS(Or Maybe make a new method for that?)
// Only Primary Key
func (driver *Driver) UpdateRecord(collection string, data string, v interface{}) error {
	if collection == "" {
		lumber.Error("Missing collection - No place to update record!")
		return fmt.Errorf("missing collection - Unable To Update")
	}

	if data == "" {
		lumber.Error("Missing Record - No place to update record!")
		return fmt.Errorf("missing data - Unable To Update")
	}

	mutex := driver.ManageMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	record := filepath.Join(driver.dir, collection, data)
	if _, err := utils.Stat(record); err != nil {
		return err
	}
	// "\t" is for indentation (Tab KEy)
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(record+".json", b, 0644); err != nil {
		return err
	}
	return nil
}

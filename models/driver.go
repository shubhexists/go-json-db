package models

import (
	"sync"
	// io/util is deprecated but I found no alternative on ChatGPT :), To be updated soon with newer one
	"io/ioutil" 
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
	"github.com/shubhexists/go-json-db/utils"
	"github.com/jcelliott/lumber"
)

type (
	Driver struct {
		// db operations should be non synchronous 
		// Got a reference from this - https://youtu.be/jkRN9zcLH1s?si=s5ec23U3tS5bi6EO
		mutex   sync.Mutex
		// this map will be used to store mutexes for each collection
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

//CREATE A NEW DATABASE (COLLECTION)
func New(dir string, options *Options)(*Driver, error){
	//This checks for any incorrect filename and corrects it.
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir: dir,
		mutexes: make(map[string]*sync.Mutex),
		log: opts.Logger,
	}

	if _,err := os.Stat(dir);
	err == nil{
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver,nil
	}
	opts.Logger.Debug("Creating the database at '%s' ...\n",dir)
	return &driver, os.MkdirAll(
		dir,
		0755)
}

//MANAGE MUTEXES FOR EACH COLLECTION
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

//WRITE ANY RECORD TO A GIVEN COLLECTION
func (driver *Driver) Write(collection string, data string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection - no place to save record")
	}

	if data == "" {
		return fmt.Errorf("missing data - Unable to save record (No Name)")
	}

	mutex := driver.ManageMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, collection)
	fnlPath := filepath.Join(dir, data+".json")
	tmpPath := fnlPath + ".tmp"
	
	if err := os.MkdirAll(dir,0755);
	err != nil{
		return err
	}
	
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tmpPath, b , 0644);
	err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

//READ ANY RECORD FROM A GIVEN COLLECTION
func (driver *Driver) Read(collection string, data string, v interface{})(string , error){
	if collection == ""{
		return "",fmt.Errorf("missing collection - Unable To Read")
	}
	
	if data == "" {
		return "",fmt.Errorf("missing data - Unable To Read")
	}

	record := filepath.Join(driver.dir, collection, data)
	if _,err := utils.Stat(record);
	err != nil{
		return "",err
	}

	b, err := ioutil.ReadFile(record+".json")
	if err != nil {
		return "",err
	}
	return string(b),nil
}

//READ ALL RECORDS FROM A GIVEN COLLECTION
//THIS WILL RETURN JSON ARRAY OF ALL THE RECORDS
func (driver *Driver) ReadAll(collection string)([]string, error){
	if collection == ""{
		return nil, fmt.Errorf("missing collection - Unable to Read Record")
	}
	dir := filepath.Join(driver.dir,collection)
	if _,err := utils.Stat(dir);
	err!= nil{
		return nil, err
	}

	files,_ := ioutil.ReadDir(dir)
	var records []string
	
	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}
	return records,nil
}

//DELETE ANY RECORD FROM A GIVEN COLLECTION
func (driver *Driver) Delete(collection string, data string) error {
	path := filepath.Join(collection,data)
	mutex := driver.ManageMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, path)
	switch fi, err := utils.Stat(dir);{
	case fi==nil, err!=nil:
		return fmt.Errorf("unable to find file or directory named %v",path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir+".json")
	}
	return nil
}

//Update any record from a given collection
//Currently we have to enter the entire User struct, UPDATE IT SO THAT WE CAN UPDATE ONLY THE REQUIRED FIELDS(Or Maybe make a new method for that?)
func (driver *Driver) UpdateRecord(collection string, data string, v interface{}) error {
	if collection == ""{
		return fmt.Errorf("missing collection - Unable To Update")
	}
	
	if data == "" {
		return fmt.Errorf("missing data - Unable To Update")
	}

	record := filepath.Join(driver.dir, collection, data)
	if _,err := utils.Stat(record);
	err != nil{
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(record+".json", b , 0644);
	err != nil {
		return err
	}
	return nil
}
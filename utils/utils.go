package utils

import (
	"os"
)

//Utility Function to verify if a specific file exists or not..
func Stat(path string)(fi os.FileInfo, err error){
	if fi, err = os.Stat(path);
	os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return
}

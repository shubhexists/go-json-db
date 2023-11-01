package utils

import (
	"os"
)

func Stat(path string)(fi os.FileInfo, err error){
	if fi, err = os.Stat(path);
	os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return
}

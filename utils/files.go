package utils

import (
	"io/ioutil"
	"log"
	"os"
)

func CreateTmpFileName(filePrefix string) string {
	tmpfile, err := ioutil.TempFile("", filePrefix)
	if err != nil {
		log.Printf("Cannot create tmp file: %s\n", err)
		return ""
	}

	fileName := tmpfile.Name()
	tmpfile.Close()
	os.Remove(tmpfile.Name())
	return fileName
}

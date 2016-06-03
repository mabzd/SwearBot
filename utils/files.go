package utils

import (
	"encoding/json"
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

func LoadJson(fileName string, in interface{}) error {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Cannot read JSON from file '%s': %v\n", fileName, err)
		return err
	}
	err = json.Unmarshal(bytes, in)
	if err != nil {
		log.Printf("Error when parsing JSON from file '%s': %v\n", fileName, err)
		return err
	}
	return nil
}

func SaveJson(fileName string, in interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		log.Printf("Error when marshaling JSON to file '%s': %v\n", fileName, err)
		return err
	}
	err = ioutil.WriteFile(fileName, bytes, 0666)
	if err != nil {
		log.Printf("Cannot write JSON to file '%s': %v\n", fileName, err)
		return err
	}
	return nil
}

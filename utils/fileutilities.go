package utils

import (
	"fmt"
	"os"
)

// CreateFile creates the new file in the "static" folder with the same name as the uploaded file
func CreateFile(filename string) (*os.File, error) {
	return os.Create(fmt.Sprintf("static/%s", filename))
}

// CreateStaticFolder creates the "static" folder if it doesn't exist
func CreateStaticFolder() error {
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		if err := os.Mkdir("static", os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

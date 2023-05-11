package api

import (
	"log"
	"os"
)

func FileSize(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Println("file err, ", err)
		return 0, err
	}
	return fileInfo.Size(), nil
}

package datalayer

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	MEMORY_DIRECTORY = "memories"
)

func GetStoreDirectory() string {

	// FILE_STORE_PATH is the folder where data is stored as json files
	FILE_STORE_PATH := os.Getenv("FILE_STORE_PATH")
	if FILE_STORE_PATH == "" {
		logrus.Fatal("env FILE_STORE_PATH not set")
	}

	return fmt.Sprintf("%s/%s", FILE_STORE_PATH, MEMORY_DIRECTORY)

}

func GetFileAbsolutePath(fileName string) string {

	// FILE_STORE_PATH is the folder where data is stored as json files
	storeDir := GetStoreDirectory()

	return fmt.Sprintf("%s/%s", storeDir, fileName)

}

func GetFileName(fileName string) string {

	return GetFileAbsolutePath(fileName)

}

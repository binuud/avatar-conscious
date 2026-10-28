package utils

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// save the image base 64 into uploads folder
// provide a unique name to the fileName, use UUID always
func SaveImageDataToFile(fileName string, rawBase64 string) error {

	// Fix Corruption: Locate and strip the "data:image/png;base64," prefix if it exists
	if idx := strings.Index(rawBase64, ","); idx != -1 {
		rawBase64 = rawBase64[idx+1:]
	}

	// 2. Decode the clean, sanitized base64 payload string
	decodedBytes, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		return status.Error(codes.InvalidArgument, "saving image to file, cannot decode file: "+err.Error())
	}

	if len(decodedBytes) == 0 {
		return status.Error(codes.InvalidArgument, "saving image to file, file empty")
	}

	destinationPath := GetEnv("FILE_STORE_PATH", ".") + "/uploads/" + fileName
	logrus.Infof("Uplaoded file to be saveed under -> %s", destinationPath)

	err = os.WriteFile(destinationPath, decodedBytes, 0644)
	if err != nil {
		return status.Error(codes.InvalidArgument, "cannot write file to server: "+err.Error())
	}
	return nil
}

package generator

import (
	"errors"
	"os"
	"path"
)

func getProcessedFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(errors.New("can't get current working directory"))
	}

	return path.Join(wd, os.Getenv("GOFILE"))
}

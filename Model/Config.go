package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		User     string `json:"user"`
		Port     string `json:"port"`
		DBName   string `json:"dbName"`
	} `json:"database"`
	AWS struct {
		S3 struct {
			AccessKeyID     string `json:"accessKeyID"`
			SecretAccessKey string `json:"secretAccessKey"`
			BucketName      string `json:"bucketName"`
		} `json:"S3"`
	} `json:"aws"`
	Etc struct {
		ScreenshotExt string `json:"screenshotExt"`
	} `json:"etc"`
}

func (c *Config) LoadConfiguration(file string) {

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(c)
}

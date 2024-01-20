package erpdb

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"upper.io/db.v3/mysql"
)

func loadEnv() error {
	err := godotenv.Load(".env") // Load environment variables from the .env file
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v", err)
	}
	return nil
}

var Settings mysql.ConnectionURL

func init() {
	if err := loadEnv(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Settings = mysql.ConnectionURL{
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_CONTAINER") + `:` + os.Getenv("DB_PORT"),
		Options: map[string]string{
			"parseTime": "true",
		},
	}
}

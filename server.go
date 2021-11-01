package main

import (
	"fmt"
	"gocv-example/app"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/middleware"
)

func sanityCheck() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("SERVER_ADDRESS") == "" ||
		os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Environment variable not defined...")
	}
}

func main() {
	sanityCheck()
	e := app.Routes()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Printf("Listening %s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))))
}

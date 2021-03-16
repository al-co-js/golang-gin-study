package main

import (
	"log"
	"os"

	"command/db"
	"command/routers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db.Connect()
	defer db.Disconnect()
	port := os.Getenv("PORT")

	app := routers.InitRouter()
	app.Run(":" + port)
}

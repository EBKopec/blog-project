package main

import (
	"github.com/blog-project/internal/handlers"
	"github.com/blog-project/internal/router"
	"github.com/blog-project/internal/storage"
	"github.com/blog-project/utils"

	"github.com/joho/godotenv"

	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}
}

func main() {

	dbConfig := storage.DBBlog{
		User:   utils.GetEnv("POSTGRES_USER"),
		Passwd: utils.GetEnv("POSTGRES_PASSWORD"),
		DBName: utils.GetEnv("POSTGRES_DB"),
	}

	db := storage.NewDBBlog(dbConfig)
	_, err := db.Open()
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.DB.Close()

	dbHandler := handlers.NewHandler(db)
	server := router.Router(dbHandler)
	if err := server.Run(); err != nil {
		log.Fatalln("error when server is initializing")
	}
}

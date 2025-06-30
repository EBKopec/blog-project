package main

import (
	"github.com/blog-project/internal/handlers"
	"github.com/blog-project/internal/router"
	"github.com/blog-project/internal/storage"
	"github.com/blog-project/utils"

	_ "github.com/blog-project/docs"
	"github.com/joho/godotenv"

	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}
}

// @title Blog API
// @version 1.0
// @description RESTful API for managing blog posts and comments
// @host localhost:8080
// @BasePath /
func main() {

	dbConfig := storage.DBBlog{
		User:   utils.GetEnv("POSTGRES_USER"),
		Passwd: utils.GetEnv("POSTGRES_PASSWORD"),
		DBName: utils.GetEnv("POSTGRES_DB"),
		DBHost: utils.GetEnv("POSTGRES_HOST"),
	}

	db := storage.NewDBBlog(dbConfig)
	clientDB, err := db.Open()
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer clientDB.DB.Close()

	dbHandler := handlers.NewHandler(clientDB)
	server := router.Router(dbHandler)
	if err := server.Run(); err != nil {
		log.Fatalln("error when server is initializing")
	}
}

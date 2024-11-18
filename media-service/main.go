package main

import (
	"context"
	"fmt"
	"log"
	"media_service/handler"
	"media_service/service"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, assuming variables set at system level")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	storageService, err := service.NewAzureBlobStorageService()
	if err != nil {
		log.Fatal(err)
	}

	fileHandler := handler.NewFileHandler(storageService)

	router := initializeRoutes(*fileHandler)

	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("Media Service listening on port %s...", port)
	log.Fatal(server.ListenAndServe())
}

func initializeRoutes(fh handler.FileHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /{roomId}/{mediaType}/{fileId}", http.HandlerFunc(fh.HandleMediaUpload))
	mux.Handle("POST /{roomId}/{mediaType}", http.HandlerFunc(fh.HandleMediaUpload))
	return mux
}

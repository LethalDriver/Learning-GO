package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	log.Println("Connected to MongoDB!")

	router := initializeRoutes() 

	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Println("Listening...")
	server.ListenAndServe() // Run the http server
}

func initializeRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /image/{imageId}", http.HandlerFunc(// TODO))
	return mux
}

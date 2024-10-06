package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/chat_app/user_service/handler"
	"example.com/chat_app/user_service/repository"
	"example.com/chat_app/user_service/service"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, assuming variables set at system level")
	}

	mongoUri := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")

	mongoClientOption := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(context.TODO(), mongoClientOption)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewMongoUserRepository(client, "chatdb", "users")

	authService, err := service.NewJwtService()
	if err != nil {
		log.Fatalf("Failed launching jwt sevice: %v", err)
	}
	userService := service.NewUserService(userRepo, authService)

	userHandler := handler.NewUserHandler(userService)

	router := initializeRoutes(userHandler) // configure routes

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Println("Listening...")
	server.ListenAndServe() // Run the http server
}

func initializeRoutes(u *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", u.HandleRegister)
	mux.HandleFunc("POST /api/login", u.HandleLogin)
	return mux
}

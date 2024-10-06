package main

import (
	"context"
	"log"
	"net/http"

	"example.com/myproject/api/handler"
	"example.com/myproject/repository"
	"example.com/myproject/service"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed loading env variables from .env: %v", err)
	}

	mongoClientOption := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), mongoClientOption)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	chatRoomRepo := repository.NewMongoChatRoomRepository(client, "chatdb", "chatrooms")
	userRepo := repository.NewMongoUserRepository(client, "chatdb", "users")

	authService, err := service.NewAuthService()
	if err != nil {
		log.Fatalf("Failed launching jwt sevice: %v", err)
	}
	userService := service.NewUserService(userRepo, authService)
	roomManager := service.NewRoomManager()
	roomService := service.NewChatService(chatRoomRepo, userRepo, roomManager)

	userHandler := handler.NewUserHandler(userService)
	wsHandler := handler.NewWebsocketHandler(roomService, userService)

	router := initializeRoutes(userHandler, wsHandler, authService) // configure routes

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Listening...")
	server.ListenAndServe() // Run the http server
}

func initializeRoutes(u *handler.UserHandler, ws *handler.WebsocketHandler, auth *service.AuthService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", u.HandleRegister)
	mux.HandleFunc("POST /api/login", u.HandleLogin)
	mux.Handle("GET /room/{roomId}", service.AuthMiddleware(auth, http.HandlerFunc(ws.HandleWebSocketUpgradeRequest)))
	return mux
}

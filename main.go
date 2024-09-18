package main

import (
	"context"
	"log"
	"net/http"

	"example.com/myproject/api/handler"
	"example.com/myproject/repository"
	"example.com/myproject/room"
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

	jwtService, err := service.NewJwtService()
	if err != nil {
		log.Fatalf("Failed launching jwt sevice: %v", err)
	}
	userService := service.NewUserService(userRepo, jwtService)
	roomManager := room.NewRoomManager(chatRoomRepo)

	userHandler := handler.NewUserHandler(userService)
	wsHandler := handler.NewWebsocketHandler(roomManager)


	router := initializeRoutes(userHandler, wsHandler) // configure routes
	  
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Listening...")
	server.ListenAndServe() // Run the http server
}
	  
func initializeRoutes(u *handler.UserHandler, ws *handler.WebsocketHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", u.HandleRegister)
	mux.HandleFunc("GET /room/{roomId}", ws.HandleWebSocketUpgradeRequest)
	return mux
}






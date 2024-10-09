package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/chat_app/chat_service/handler"
	"example.com/chat_app/chat_service/repository"
	"example.com/chat_app/chat_service/service"
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

	chatRoomRepo := repository.NewMongoChatRoomRepository(client, "chatdb", "chatrooms")

	roomManager := service.NewRoomManager()
	roomService := service.NewChatService(chatRoomRepo, roomManager)

	wsHandler := handler.NewWebsocketHandler(roomService)

	router := initializeRoutes(wsHandler) // configure routes

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("Chat Service listening on port %s...", port)
	log.Fatal(server.ListenAndServe())
}

func initializeRoutes(ws *handler.WebsocketHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /room/{roomId}", http.HandlerFunc(ws.HandleWebSocketUpgradeRequest))
	return mux
}

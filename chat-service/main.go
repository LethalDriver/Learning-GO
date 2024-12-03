package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/chat_app/chat_service/client"
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

	mongoClient, err := mongo.Connect(context.TODO(), mongoClientOption)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.TODO())

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	chatRoomRepo := repository.NewMongoChatRoomRepository(mongoClient, "chatdb", "chatrooms")
	mediaRepo := repository.NewMongoFileRepository(mongoClient, "chatdb", "mediafiles")

	mediaServiceClient, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	roomManager := service.NewRoomManager()
	chatService := service.NewChatService(chatRoomRepo, roomManager)
	roomService := service.NewRoomService(chatRoomRepo)
	mediaService := service.NewMediaService(mediaRepo, mediaServiceClient)

	wsHandler := handler.NewWebsocketHandler(chatService)
	roomHandler := handler.NewRoomHandler(roomService)
	mediaHandler := handler.NewMediaHandler(mediaService)

	router := initializeRoutes(wsHandler, roomHandler, mediaHandler) // configure routes

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("Chat Service listening on port %s...", port)
	log.Fatal(server.ListenAndServe())
}

func initializeRoutes(ws *handler.WebsocketHandler, rh *handler.RoomHandler, mh *handler.MediaHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /room/{roomId}/connect", http.HandlerFunc(ws.HandleWebSocketUpgradeRequest))
	mux.Handle("GET /room/{roomId}", http.HandlerFunc(rh.GetRoom))
	mux.Handle("POST /room", http.HandlerFunc(rh.CreateRoom))
	mux.Handle("DELETE /room", http.HandlerFunc(rh.DeleteRoom))
	mux.Handle("POST /room/{roomId}/users/add", http.HandlerFunc(rh.AddUsersToRoom))
	mux.Handle("PATCH /room/{roomId}/users/{userId}/promote", http.HandlerFunc(rh.PromoteUser))
	mux.Handle("PATCH /room/{roomId}/users/{userId}/demote", http.HandlerFunc(rh.DemoteUser))
	mux.Handle("DELETE /room/{roomId}/users/{userId}", http.HandlerFunc(rh.DeleteUserFromRoom))
	mux.Handle("DELETE /room/{roomId}/users/me", http.HandlerFunc(rh.LeaveRoom))
	mux.Handle("POST /media/upload", http.HandlerFunc(mh.UploadMedia))
	mux.Handle("GET /media/{mediaId}/download", http.HandlerFunc(mh.GetMediaFile))
	mux.Handle("GET /media/{mediaId}", http.HandlerFunc(mh.GetMediaMetadata))
	return mux
}

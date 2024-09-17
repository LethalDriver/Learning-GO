package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
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

	chatRoomRepo := NewMongoChatRoomRepository(client, "chatdb", "chatrooms")
	userRepo := NewMongoUserRepository(client, "chatdb", "users")

	userService := NewUserService(userRepo)
	roomManager := NewRoomManager(chatRoomRepo)

	userHandler := NewUserHandler(userService)
	wsHandler := NewWebsocketHandler(roomManager)


	router := initializeRoutes(userHandler, wsHandler) // configure routes
	  
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Listening...")
	server.ListenAndServe() // Run the http server
}
	  
func initializeRoutes(u *UserHandler, ws *WebsocketHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/register", u.HandleRegister)
	mux.HandleFunc("GET /room/{roomId}", ws.handleWebSocketUpgradeRequest)
	return mux
}






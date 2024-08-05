package main

import (
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type RoomManager struct {
	rooms map[string]*ChatRoomWebsocket
	lock sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*ChatRoomWebsocket),
	}
}

func (manager *RoomManager) GetOrCreateRoom(roomId string, repo ChatRoomRepository, conn *Connection) *ChatRoomWebsocket {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	var roomExists bool

	roomEntity, err := repo.GetRoom(roomId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Room: %s doesn't exist", roomId)
		} else {
			log.Printf("Error checking for existance room: %s", roomId)
		}
	} else {
		roomExists = true
	}

	if roomExists {
		if room, exists := manager.rooms[roomId]; exists {
			for _, message := range(roomEntity.Messages){
				conn.send <-[]byte(message.Content)
			}
			return room
		}
	} else {
		repo.CreateRoom(roomId)
	}

	newRoom := &ChatRoomWebsocket{
		Id: roomId,
		Members: make(map[*Connection]bool),
		Broadcast: make(chan []byte),
		Register: make(chan *Connection),
		Unregister: make(chan *Connection),
	}
	manager.rooms[roomId] = newRoom
	
	go newRoom.Run(repo)

	return newRoom
}
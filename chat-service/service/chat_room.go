package service

import (
	"context"
	"log"

	"example.com/chat_app/chat_service/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

// ChatRoom represents a chat room with its members and channels for various operations.
// It is a struct represetning a chat room instance in memory, different from the ChatRoomEntity in the repository package.
type ChatRoom struct {
	Id         string
	Members    map[*Connection]bool
	Text       chan structs.Message
	Seen       chan structs.SeenMessage
	Delete     chan structs.DeleteMessage
	Register   chan *Connection
	Unregister chan *Connection
}

// NewChatRoom creates a new instance of ChatRoom.
// It takes a room ID as a parameter and initializes the channels and members map.
func NewChatRoom(roomId string) *ChatRoom {
	return &ChatRoom{
		Id:         roomId,
		Members:    make(map[*Connection]bool),
		Text:       make(chan structs.Message),
		Seen:       make(chan structs.SeenMessage),
		Delete:     make(chan structs.DeleteMessage),
		Register:   make(chan *Connection),
		Unregister: make(chan *Connection),
	}
}

// Run starts the chat room's main loop, handling registration, unregistration, and message broadcasting.
func (r *ChatRoom) Run(service *ChatService) {
	ctx := context.Background()
	for {
		select {
		case conn := <-r.Register:
			log.Printf("Registering connection to room %s, address: %p", r.Id, conn)
			room, err := service.roomRepo.GetRoom(ctx, r.Id)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					log.Printf("No messages found for room %s", r.Id)
				} else {
					log.Printf("Error getting messages for room %s: %v", r.Id, err)
				}
				break
			}
			service.pumpExistingMessages(conn, room.Messages)
			r.Members[conn] = true

		case conn := <-r.Unregister:
			log.Printf("Unregistering connection from room %s", r.Id)
			if _, ok := r.Members[conn]; ok {
				delete(r.Members, conn)
				close(conn.sendMessage)
			}

		case message := <-r.Text:
			log.Printf("Broadcasting message to room %s: %s", r.Id, string(message.Content))
			message, err := service.processAndSaveMessage(ctx, r.Id, &message)
			if err != nil {
				log.Printf("Error saving message %q in room %s", string(message.Content), r.Id)
				break
			}
			for conn := range r.Members {
				select {
				case conn.sendMessage <- message:
				default:
					close(conn.sendMessage)
					delete(r.Members, conn)
				}
			}

		case seenMessage := <-r.Seen:
			log.Printf("Broadcasting seen update to room %s: %s", r.Id, seenMessage.MessageId)
			err := service.roomRepo.InsertSeenBy(ctx, r.Id, seenMessage.MessageId, seenMessage.SeenBy.Id)
			if err != nil {
				log.Printf("Error saving seen update for message %s in room %s", seenMessage.MessageId, r.Id)
				break
			}
			for conn := range r.Members {
				select {
				case conn.sendSeen <- seenMessage:
				default:
					close(conn.sendSeen)
					delete(r.Members, conn)
				}
			}

		case deleteMessage := <-r.Delete:
			log.Printf("Broadcasting delete message to room %s: %s", r.Id, deleteMessage.MessageId)
			err := service.roomRepo.DeleteMessage(ctx, r.Id, deleteMessage.MessageId)
			if err != nil {
				log.Printf("Error deleting message %s in room %s", deleteMessage.MessageId, r.Id)
				break
			}
			for conn := range r.Members {
				select {
				case conn.sendDelete <- deleteMessage:
				default:
					close(conn.sendDelete)
					delete(r.Members, conn)
				}
			}
		}
	}
}

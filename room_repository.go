package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type ChatRoomRepository interface {
	CreateRoom(id string) (ChatRoomEntity, error)
	GetRoom(id string) (ChatRoomEntity, error)
	DeleteRoom(id string) error
	RoomExists(id string) (bool, error)
	AddMessageToRoom(roomId string, content string) error
}

type MongoChatRoomRepository struct {
	collection *mongo.Collection
}

func NewMongoChatRoomRepository(client *mongo.Client, dbName, collectionName string) *MongoChatRoomRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoChatRoomRepository{collection: collection}
}

func (repo *MongoChatRoomRepository) CreateRoom(id string) (ChatRoomEntity, error) {
	newRoom := &ChatRoomEntity{
		Id: id,
		Messages: []MessageEntity{},
	}

	_, err := repo.collection.InsertOne(context.TODO(), newRoom)
	if err != nil {
		return ChatRoomEntity{}, err
	}
	return *newRoom, nil
}

func (repo *MongoChatRoomRepository) GetRoom(id string) (ChatRoomEntity, error) {
	var room ChatRoomEntity
	filter := bson.D{{Key: "id", Value: id}}
	err := repo.collection.FindOne(context.TODO(), filter).Decode(&room)
	if err != nil {
		return ChatRoomEntity{}, err
	}
	return room, nil
}

func (repo *MongoChatRoomRepository) DeleteRoom(id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := repo.collection.DeleteOne(context.TODO(), filter)
	return err
}

func (repo *MongoChatRoomRepository) RoomExists(id string) (bool, error) {
	filter := bson.M{"id": id}
	var result ChatRoomWebsocket
	err := repo.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil 
		}
		return false, err
	}
	return true, nil
}

func (repo *MongoChatRoomRepository) AddMessageToRoom(roomId string, content string) error {
	message := NewMessageEntity(content, roomId)
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$push": bson.M{
			"messages": message,
		},
	}
	_, err := repo.collection.UpdateOne(context.TODO(), filter, update) 
	return err
}
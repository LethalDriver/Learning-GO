package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type ChatRoomRepository interface {
	CreateRoom(id string) (*ChatRoomEntity, error)
	GetRoom(id string) (*ChatRoomEntity, error)
	DeleteRoom(id string) error
	AddMessageToRoom(roomId string, content string) error
}

type MongoChatRoomRepository struct {
	collection *mongo.Collection
}

func (repo *MongoChatRoomRepository) GetCollection() *mongo.Collection {
	return repo.collection
}

func NewMongoChatRoomRepository(client *mongo.Client, dbName, collectionName string) *MongoChatRoomRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoChatRoomRepository{collection: collection}
}

func (repo *MongoChatRoomRepository) CreateRoom(id string) (*ChatRoomEntity, error) {
	newRoom := &ChatRoomEntity{
		Id: id,
		Messages: []MessageEntity{},
	}

	_, err := repo.collection.InsertOne(context.TODO(), newRoom)
	if err != nil {
		return nil, err
	}
	return newRoom, nil
}

func (repo *MongoChatRoomRepository) GetRoom(id string) (*ChatRoomEntity, error) {
	return GetByKey[ChatRoomEntity, string]("id", id, repo)
}

func (repo *MongoChatRoomRepository) DeleteRoom(id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := repo.collection.DeleteOne(context.TODO(), filter)
	return err
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
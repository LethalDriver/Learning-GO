package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type ChatRoomRepository interface {
	SaveRoom(room ChatRoom) (ChatRoom, error)
	GetRoom(id string) (ChatRoom, error)
	UpdateRoom(id string, room ChatRoom) (ChatRoom, error)
	DeleteRoom(id string) error
}

type MongoChatRoomRepository struct {
	collection *mongo.Collection
}

func NewMongoChatRoomRepository(client *mongo.Client, dbName, collectionName string) *MongoChatRoomRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoChatRoomRepository{collection: collection}
}

func (repo *MongoChatRoomRepository) SaveRoom(room ChatRoom) (ChatRoom, error) {
	_, err := repo.collection.InsertOne(context.TODO(), room)
	if err != nil {
		return ChatRoom{}, err
	}
	return room, nil
}

func (repo *MongoChatRoomRepository) GetRoom(id string) (ChatRoom, error) {
	var room ChatRoom
	filter := bson.D{{Key: "_id", Value: id}}
	err := repo.collection.FindOne(context.TODO(), filter).Decode(&room)
	if err != nil {
		return ChatRoom{}, err
	}
	return room, nil
}

func (repo *MongoChatRoomRepository) UpdateRoom(id string, room ChatRoom) (ChatRoom, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: room}}
	_, err := repo.collection.UpdateOne(context.TODO(), filter, update)
	return room, err
}

func (repo *MongoChatRoomRepository) DeleteRoom(id string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := repo.collection.DeleteOne(context.TODO(), filter)
	return err
}
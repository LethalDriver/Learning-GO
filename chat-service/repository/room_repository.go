package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrRoomExists = errors.New("room already exists")

type ChatRoomRepository interface {
	CreateRoom(ctx context.Context, id string) (*ChatRoomEntity, error)
	GetRoom(ctx context.Context, id string) (*ChatRoomEntity, error)
	DeleteRoom(ctx context.Context, id string) error
	AddMessageToRoom(ctx context.Context, roomId string, message *Message) error
	InsertSeenBy(ctx context.Context, roomId string, messageId string, userId string) error
	DeleteMessage(ctx context.Context, roomId string, messageId string) error
}

type MongoChatRoomRepository struct {
	collection *mongo.Collection
}

func NewMongoChatRoomRepository(client *mongo.Client, dbName, collectionName string) *MongoChatRoomRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoChatRoomRepository{collection: collection}
}

func (repo *MongoChatRoomRepository) GetRoom(ctx context.Context, id string) (*ChatRoomEntity, error) {
	var room ChatRoomEntity
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (repo *MongoChatRoomRepository) CreateRoom(ctx context.Context, id string) (*ChatRoomEntity, error) {
	// Check if a room with the given ID already exists
	_, err := repo.GetRoom(ctx, id)
	if err == nil {
		return nil, ErrRoomExists
	} else if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("error checking for existing room: %w", err)
	}

	// Create a new room if it does not exist
	newRoom := &ChatRoomEntity{
		Id:       id,
		Messages: []Message{},
	}

	_, err = repo.collection.InsertOne(ctx, newRoom)
	if err != nil {
		return nil, fmt.Errorf("error inserting new room: %w", err)
	}
	return newRoom, nil
}


func (repo *MongoChatRoomRepository) DeleteRoom(ctx context.Context, id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := repo.collection.DeleteOne(ctx, filter)
	return err
}

func (repo *MongoChatRoomRepository) AddMessageToRoom(ctx context.Context, roomId string, message *Message) error {
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$push": bson.M{
			"messages": message,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

func (repo *MongoChatRoomRepository) InsertSeenBy(ctx context.Context, roomId string, messageId string, userId string) error {
	filter := bson.M{"id": roomId, "messages.id": messageId}
	update := bson.M{
		"$addToSet": bson.M{
			"messages.$.seenBy": userId,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

func (repo *MongoChatRoomRepository) DeleteMessage(ctx context.Context, roomId string, messageId string) error {
	filter := bson.M{"id": roomId, "messages.id": messageId}
	update := bson.M{
		"$pull": bson.M{
			"messages": bson.M{"id": messageId},
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

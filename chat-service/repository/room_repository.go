package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoomRepository interface {
	CreateRoom(ctx context.Context) (*ChatRoomEntity, error)
	GetRoom(ctx context.Context, id string) (*ChatRoomEntity, error)
	DeleteRoom(ctx context.Context, id string) error
	AddMessageToRoom(ctx context.Context, roomId string, message *Message) error
	InsertSeenBy(ctx context.Context, roomId string, messageId string, userId string) error
	DeleteMessage(ctx context.Context, roomId string, messageId string) error
	InsertUserIntoRoom(ctx context.Context, roomId string, user UserPermissions) error
	DeleteUserFromRoom(ctx context.Context, roomId string, userId string) error
	GetUsersPermissions(ctx context.Context, roomId string, userId string) (*UserPermissions, error)
	PromoteUserToAdmin(ctx context.Context, roomId string, userId string) error
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

func (repo *MongoChatRoomRepository) CreateRoom(ctx context.Context) (*ChatRoomEntity, error) {
	// Create a new room if it does not exist
	newRoom := &ChatRoomEntity{
		Id:       uuid.NewString(),
		Messages: []Message{},
		Users:    []UserPermissions{},
	}

	_, err := repo.collection.InsertOne(ctx, newRoom)
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

func (repo *MongoChatRoomRepository) InsertUserIntoRoom(ctx context.Context, roomId string, user UserPermissions) error {
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$addToSet": bson.M{
			"users": user,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

func (repo *MongoChatRoomRepository) PromoteUserToAdmin(ctx context.Context, roomId string, userId string) error {
	filter := bson.M{"id": roomId, "users.userId": userId}
	update := bson.M{
		"$set": bson.M{
			"users.$.role": Admin,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

func (repo *MongoChatRoomRepository) DeleteUserFromRoom(ctx context.Context, roomId string, userId string) error {
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$pull": bson.M{
			"users": bson.M{"userId": userId},
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

func (repo *MongoChatRoomRepository) GetUsersPermissions(ctx context.Context, roomId string, userId string) (*UserPermissions, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "id", Value: roomId}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$users"}}}},
		{{Key: "$match", Value: bson.D{{Key: "users.userId", Value: userId}}}},
		{{Key: "$project", Value: bson.D{{Key: "userPermissions", Value: "$users"}}}},
	}

	// Execute the aggregation pipeline
	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, mongo.ErrNoDocuments
	}

	var result struct {
		UserPermissions UserPermissions `bson:"userPermissions"`
	}
	if err := cursor.Decode(&result); err != nil {
		return nil, err
	}

	return &result.UserPermissions, nil
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

package repository

import (
	"context"
	"fmt"

	"example.com/chat_app/chat_service/structs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoChatRoomRepository provides methods to interact with the chat room collection in MongoDB.
type MongoChatRoomRepository struct {
	collection *mongo.Collection
}

// NewMongoChatRoomRepository creates a new instance of MongoChatRoomRepository.
// It takes a MongoDB client, database name, and collection name as parameters.
func NewMongoChatRoomRepository(client *mongo.Client, dbName, collectionName string) *MongoChatRoomRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoChatRoomRepository{collection: collection}
}

// GetRoom retrieves a chat room from the MongoDB collection by its ID.
func (repo *MongoChatRoomRepository) GetRoom(ctx context.Context, id string) (*structs.ChatRoomEntity, error) {
	var room structs.ChatRoomEntity
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// CreateRoom creates a new chat room in the MongoDB collection.
func (repo *MongoChatRoomRepository) CreateRoom(ctx context.Context) (*structs.ChatRoomEntity, error) {
	newRoom := &structs.ChatRoomEntity{
		Id:       uuid.NewString(),
		Messages: []structs.Message{},
		Users:    []structs.UserPermissions{},
	}

	_, err := repo.collection.InsertOne(ctx, newRoom)
	if err != nil {
		return nil, fmt.Errorf("error inserting new room: %w", err)
	}
	return newRoom, nil
}

// DeleteRoom removes a chat room from the MongoDB collection by its ID.
func (repo *MongoChatRoomRepository) DeleteRoom(ctx context.Context, id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := repo.collection.DeleteOne(ctx, filter)
	return err
}

// AddMessageToRoom adds a message to a chat room in the MongoDB collection.
func (repo *MongoChatRoomRepository) AddMessageToRoom(ctx context.Context, roomId string, message *structs.Message) error {
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$push": bson.M{
			"messages": message,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

// InsertUserIntoRoom adds a user to a chat room in the MongoDB collection.
func (repo *MongoChatRoomRepository) InsertUserIntoRoom(ctx context.Context, roomId string, user structs.UserPermissions) error {
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$addToSet": bson.M{
			"users": user,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

// ChangeUserRole changes the role of a user in a chat room in the MongoDB collection.
func (repo *MongoChatRoomRepository) ChangeUserRole(ctx context.Context, roomId string, userId string, role structs.Role) error {
	filter := bson.M{"id": roomId, "users.userId": userId}
	update := bson.M{
		"$set": bson.M{
			"users.$.role": role,
		},
	}
	_, err := repo.collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteUserFromRoom removes a user from a chat room in the MongoDB collection.
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

// GetUsersPermissions retrieves the permissions of a user in a chat room from the MongoDB collection.
func (repo *MongoChatRoomRepository) GetUsersPermissions(ctx context.Context, roomId string, userId string) (*structs.UserPermissions, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "id", Value: roomId}}}},
		bson.D{{Key: "$unwind", Value: "$users"}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "users.userId", Value: userId}}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "userPermissions", Value: "$users"}}}},
	}

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, mongo.ErrNoDocuments
	}

	var result struct {
		UserPermissions structs.UserPermissions `bson:"userPermissions"`
	}
	if err := cursor.Decode(&result); err != nil {
		return nil, err
	}

	return &result.UserPermissions, nil
}

// InsertSeenBy adds a user to the seenBy list of a message in a chat room in the MongoDB collection.
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

// DeleteMessage removes a message from a chat room in the MongoDB collection.
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

func (repo *MongoChatRoomRepository) GetUnseenMessages(ctx context.Context, roomId, userId string) ([]structs.Message, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "id", Value: roomId}}}},
		bson.D{{Key: "$unwind", Value: "$messages"}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "messages.seenBy", Value: bson.D{{Key: "$ne", Value: userId}}}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "messages.SentAt", Value: 1}}}}, // Sort messages by SentAt in ascending order
		bson.D{{Key: "$project", Value: bson.D{{Key: "message", Value: "$messages"}}}},
	}

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []structs.Message
	for cursor.Next(ctx) {
		var message structs.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil 
}

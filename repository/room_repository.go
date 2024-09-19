package repository

import (
	"context"

	"example.com/myproject/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoomRepository interface {
    CreateRoom(ctx context.Context, id string) (*entity.ChatRoomEntity, error)
    GetRoom(ctx context.Context, id string) (*entity.ChatRoomEntity, error)
    DeleteRoom(ctx context.Context, id string) error
    AddMessageToRoom(ctx context.Context, roomId string, content string) error
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

func (repo *MongoChatRoomRepository) CreateRoom(ctx context.Context, id string) (*entity.ChatRoomEntity, error) {
    newRoom := &entity.ChatRoomEntity{
        Id:       id,
        Messages: []entity.MessageEntity{},
    }

    _, err := repo.collection.InsertOne(ctx, newRoom)
    if err != nil {
        return nil, err
    }
    return newRoom, nil
}

func (repo *MongoChatRoomRepository) GetRoom(ctx context.Context, id string) (*entity.ChatRoomEntity, error) {
    return GetByKey[entity.ChatRoomEntity](ctx, "id", id, repo)
}

func (repo *MongoChatRoomRepository) DeleteRoom(ctx context.Context, id string) error {
    filter := bson.D{{Key: "id", Value: id}}
    _, err := repo.collection.DeleteOne(ctx, filter)
    return err
}

func (repo *MongoChatRoomRepository) AddMessageToRoom(ctx context.Context, roomId string, content string) error {
    message := entity.NewMessageEntity(content, roomId)
    filter := bson.M{"id": roomId}
    update := bson.M{
        "$push": bson.M{
            "messages": message,
        },
    }
    _, err := repo.collection.UpdateOne(ctx, filter, update)
    return err
}
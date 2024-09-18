package repository

import (
	"context"

	"example.com/myproject/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type ChatRoomRepository interface {
	CreateRoom(id string) (*entity.ChatRoomEntity, error)
	GetRoom(id string) (*entity.ChatRoomEntity, error)
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

func (repo *MongoChatRoomRepository) CreateRoom(id string) (*entity.ChatRoomEntity, error) {
	newRoom := &entity.ChatRoomEntity{
		Id: id,
		Messages: []entity.MessageEntity{},
	}

	_, err := repo.collection.InsertOne(context.TODO(), newRoom)
	if err != nil {
		return nil, err
	}
	return newRoom, nil
}

func (repo *MongoChatRoomRepository) GetRoom(id string) (*entity.ChatRoomEntity, error) {
	return GetByKey[entity.ChatRoomEntity]("id", id, repo)
}

func (repo *MongoChatRoomRepository) DeleteRoom(id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := repo.collection.DeleteOne(context.TODO(), filter)
	return err
}

func (repo *MongoChatRoomRepository) AddMessageToRoom(roomId string, content string) error {
	message := entity.NewMessageEntity(content, roomId)
	filter := bson.M{"id": roomId}
	update := bson.M{
		"$push": bson.M{
			"messages": message,
		},
	}
	_, err := repo.collection.UpdateOne(context.TODO(), filter, update) 
	return err
}
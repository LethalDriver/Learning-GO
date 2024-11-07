package repository

import (
	"context"
	"fmt"
	"media_service/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoFileRepository struct {
	collection *mongo.Collection
}

func NewMongoFileRepository(client *mongo.Client, dbName string) *MongoFileRepository {
	return &MongoFileRepository{
		collection: client.Database(dbName).Collection("media"),
	}
}

func (repo *MongoFileRepository) GetFile(ctx context.Context, id string, mediaType structs.MediaType) (*structs.MediaFile, error) {
	var file structs.MediaFile
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (repo *MongoFileRepository) DeleteFile(ctx context.Context, id string, mediaType structs.MediaType) error {
	filter := bson.M{"id": id}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	return nil
}

func (repo *MongoFileRepository) SaveFile(ctx context.Context, file *structs.MediaFile, mediaType structs.MediaType) error {
	_, err := repo.collection.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	return nil
}

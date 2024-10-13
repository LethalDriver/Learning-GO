package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	GetFile(ctx context.Context, id string, mediaType MediaType) (*MediaFile, error)
	DeleteFile(ctx context.Context, userId string, mediaType MediaType) error
	SaveFile(ctx context.Context, file *MediaFile, mediaType MediaType) error
}

type MongoFileRepository struct {
	collections map[MediaType]*mongo.Collection
}

func NewMongoFileRepository(client *mongo.Client, dbName, imageColl, videoColl, audioColl, othersColl string) *MongoFileRepository {
	return &MongoFileRepository{
		collections: map[MediaType]*mongo.Collection{
			Image: client.Database(dbName).Collection(imageColl),
			Video: client.Database(dbName).Collection(videoColl),
			Audio: client.Database(dbName).Collection(audioColl),
			Other: client.Database(dbName).Collection(othersColl),
		},
	}
}

func (repo *MongoFileRepository) getCollection(mediaType MediaType) (*mongo.Collection, error) {
	collection, exists := repo.collections[mediaType]
	if !exists {
		return nil, fmt.Errorf("unsupported media type: %v", mediaType)
	}
	return collection, nil
}

func (repo *MongoFileRepository) GetImage(ctx context.Context, id string, mediaType MediaType) (*MediaFile, error) {
	var file MediaFile
	filter := bson.M{"id": id}
	collection, err := repo.getCollection(mediaType)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (repo *MongoFileRepository) DeleteImage(ctx context.Context, id string, mediaType MediaType) error {
	filter := bson.M{"id": id}
	collection, err := repo.getCollection(mediaType)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	return nil
}

func (repo *MongoFileRepository) CreateImage(ctx context.Context, file *MediaFile, mediaType MediaType) error {
	collection, err := repo.getCollection(mediaType)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	return nil
}

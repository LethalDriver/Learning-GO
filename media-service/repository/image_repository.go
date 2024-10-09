package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageRepository interface {
	GetImage(ctx context.Context, id string) (*MediaFile, error)
	DeleteImage(ctx context.Context, userId string) error
	CreateImage(ctx context.Context, file *MediaFile) error
}

type MongoImageRepository struct {
	collection *mongo.Collection
}

func NewMongoImageRepository(client *mongo.Client, dbName, collectionName string) *MongoImageRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoImageRepository{collection: collection}
}

func (repo *MongoImageRepository) GetImage(ctx context.Context, id string) (*MediaFile, error) {
	var image MediaFile
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (repo *MongoImageRepository) DeleteImage(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}

func (repo *MongoImageRepository) CreateImage(ctx context.Context, file *MediaFile) error {
	_, err := repo.collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	return nil
}

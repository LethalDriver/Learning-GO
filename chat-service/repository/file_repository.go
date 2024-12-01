package repository

import (
	"context"
	"fmt"

	"example.com/chat_app/chat_service/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoFileRepository provides methods to interact with the file collection in MongoDB.
type MongoFileRepository struct {
	collection *mongo.Collection
}

// NewMongoFileRepository creates a new instance of MongoFileRepository.
// It takes a MongoDB client, database name, and collection name as parameters.
func NewMongoFileRepository(client *mongo.Client, dbName, collection string) *MongoFileRepository {
	return &MongoFileRepository{
		collection: client.Database(dbName).Collection(collection),
	}
}

// GetFile retrieves a file from the MongoDB collection by its ID.
func (repo *MongoFileRepository) GetFile(ctx context.Context, id string) (*structs.MediaFile, error) {
	var file structs.MediaFile
	filter := bson.M{"id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteFile removes a file from the MongoDB collection by its ID.
func (repo *MongoFileRepository) DeleteFile(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	return nil
}

// SaveFile inserts a new file into the MongoDB collection.
func (repo *MongoFileRepository) SaveFile(ctx context.Context, file *structs.MediaFile) error {
	_, err := repo.collection.InsertOne(ctx, file)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	return nil
}

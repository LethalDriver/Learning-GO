package repository

import (
	"context"

	"example.com/chat_app/user_service/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoUserRepository.
func NewMongoUserRepository(client *mongo.Client, dbName, collectionName string) *MongoUserRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoUserRepository{collection: collection}
}

// GetById retrieves a user by their ID.
func (repo *MongoUserRepository) GetById(ctx context.Context, id string) (*structs.UserEntity, error) {
	return getByKey[structs.UserEntity](ctx, "id", id, repo.collection)
}

// GetByUsername retrieves a user by their username.
func (repo *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*structs.UserEntity, error) {
	return getByKey[structs.UserEntity](ctx, "username", username, repo.collection)
}

// Save stores a user entity in the repository.
func (repo *MongoUserRepository) Save(ctx context.Context, user *structs.UserEntity) error {
	_, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// GetByKey filters the collection by the given key and returns the result.
func getByKey[T, V any](ctx context.Context, key string, value V, collection *mongo.Collection) (*T, error) {
	var entity T
	filter := bson.M{key: value}
	err := collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

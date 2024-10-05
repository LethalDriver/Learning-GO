package repository

import (
	"context"

	"example.com/myproject/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetById(ctx context.Context, id string) (*structs.UserEntity, error)
    GetByUsername(ctx context.Context, username string) (*structs.UserEntity, error)
    Save(ctx context.Context, user *structs.UserEntity) error
}

type MongoUserRepository struct {
    collection *mongo.Collection
}

func (repo *MongoUserRepository) GetCollection() *mongo.Collection {
    return repo.collection
}

func NewMongoUserRepository(client *mongo.Client, dbName, collectionName string) *MongoUserRepository {
    collection := client.Database(dbName).Collection(collectionName)
    return &MongoUserRepository{collection: collection}
}

func (repo *MongoUserRepository) GetById(ctx context.Context, id string) (*structs.UserEntity, error) {
    return GetByKey[structs.UserEntity](ctx, "id", id, repo)
}

func (repo *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*structs.UserEntity, error) {
    return GetByKey[structs.UserEntity](ctx, "username", username, repo)
}

func (repo *MongoUserRepository) Save(ctx context.Context, user *structs.UserEntity) error {
    _, err := repo.collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }
    return nil
}
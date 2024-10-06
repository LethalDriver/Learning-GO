package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntity struct {
	Id       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type UserRepository interface {
	GetById(ctx context.Context, id string) (*UserEntity, error)
	GetByUsername(ctx context.Context, username string) (*UserEntity, error)
	Save(ctx context.Context, user *UserEntity) error
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

func (repo *MongoUserRepository) GetById(ctx context.Context, id string) (*UserEntity, error) {
	return GetByKey[UserEntity](ctx, "id", id, repo)
}

func (repo *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*UserEntity, error) {
	return GetByKey[UserEntity](ctx, "username", username, repo)
}

func (repo *MongoUserRepository) Save(ctx context.Context, user *UserEntity) error {
	_, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

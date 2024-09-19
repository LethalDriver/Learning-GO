package repository

import (
	"context"

	"example.com/myproject/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetById(id string, ctx context.Context) (*entity.UserEntity, error)
	GetByUsername(username string, ctx context.Context) (*entity.UserEntity, error)
	Save(user *entity.UserEntity, ctx context.Context) error
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

func (repo *MongoUserRepository) GetById(id string, ctx context.Context) (*entity.UserEntity, error) {
	return GetByKey[entity.UserEntity]("id", id, repo, ctx)
}

func (repo *MongoUserRepository) GetByUsername(username string, ctx context.Context) (*entity.UserEntity, error) {
	return GetByKey[entity.UserEntity]("username", username, repo, ctx)
}

 func (repo *MongoUserRepository) Save(user *entity.UserEntity, ctx context.Context) error {
	_, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
 }
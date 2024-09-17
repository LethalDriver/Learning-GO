package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetById(id string) (*UserEntity, error)
	GetByUsername(username string) (*UserEntity, error)
	Save(user *UserEntity) error
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

func (repo *MongoUserRepository) GetById(id string) (*UserEntity, error) {
	return GetByKey[UserEntity, string]("id", id, repo)
}

func (repo *MongoUserRepository) GetByUsername(username string) (*UserEntity, error) {
	return GetByKey[UserEntity, string]("username", username, repo)
}

 func (repo *MongoUserRepository) Save(user *UserEntity) error {
	_, err := repo.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
 }
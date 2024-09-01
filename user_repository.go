package main

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetUserById(id string) (UserEntity, error)
	GetUserByUsername(username string) (UserEntity, error)
	RegisterUser(username string, password string)
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

func (repo *MongoUserRepository) GetUserById(id string) (*UserEntity, error) {
	return GetByKey[UserEntity, string]("id", id, repo)
}

func (repo *MongoUserRepository) GetByUsername(username string) (*UserEntity, error) {
	return GetByKey[UserEntity, string]("username", username, repo)
}
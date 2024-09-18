package repository

import (
	"context"

	"example.com/myproject/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetById(id string) (*entity.UserEntity, error)
	GetByUsername(username string) (*entity.UserEntity, error)
	Save(user *entity.UserEntity) error
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

func (repo *MongoUserRepository) GetById(id string) (*entity.UserEntity, error) {
	return GetByKey[entity.UserEntity]("id", id, repo)
}

func (repo *MongoUserRepository) GetByUsername(username string) (*entity.UserEntity, error) {
	return GetByKey[entity.UserEntity]("username", username, repo)
}

 func (repo *MongoUserRepository) Save(user *entity.UserEntity) error {
	_, err := repo.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
 }
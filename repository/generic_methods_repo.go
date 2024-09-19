package repository

import (
	"context"

	"example.com/myproject/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository interface to abstract the MongoDB collection access
type Repository interface {
    GetCollection() *mongo.Collection
}

// Comparable is a constraint that ensures the type supports == and != operators.
type Comparable interface {
    comparable
}

// GetByKey filters the collection by the given key and returns the result.
func GetByKey[T entity.Entity, V Comparable](ctx context.Context, key string, value V, repo Repository) (*T, error) {
    var entity T
    filter := bson.D{{Key: key, Value: value}}
    err := repo.GetCollection().FindOne(ctx, filter).Decode(&entity)
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// DeleteById deletes a document from the collection by its ID.
func DeleteById[V Comparable](ctx context.Context, id V, repo Repository) error {
    filter := bson.D{{Key: "id", Value: id}}
    _, err := repo.GetCollection().DeleteOne(ctx, filter)
    return err
}
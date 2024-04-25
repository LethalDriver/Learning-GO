package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Message struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	ChannelId int    `json:"channelId"`
}

type MessageRepository interface {
	FindAll() ([]*Message, error)
	FindById(id int) (*Message, error)
	Create(message *Message) error
}

type RedisMessageRepository struct {
	redisClient *redis.Client
}

func NewRedisMessageRepository(redisClient *redis.Client) *RedisMessageRepository {
	return &RedisMessageRepository{redisCleint: redisClient}
}

func (r *RedisMessageRepository) FindAll() ([]*Message, error) {
	var messages []*Message
	keys, err := r.redisClient.Keys("message:*").Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		message, err := r.FindById(key)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *RedisMessageRepository) FindById(id int) (*Message, error) {
	message := &Message{}
	err := r.redisClient.Get(fmt.Sprintf("message:%d", id)).Scan(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *RedisMessageRepository) Create(message *Message) error {
	return r.redisClient.Set(fmt.Sprintf("message:%d", message.Id), message, 0).Err()
}
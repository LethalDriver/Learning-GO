package main

type Message struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	ChannelId int    `json:"channelId"`
}

type Room struct {
	Id   int    `json:"id"`
	Messages [Message] `json:"messages"`
}
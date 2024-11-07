package structs

type RoomDto struct {
	Id      string    `json:"id"`
	Members []UserDto `json:"members"`
}

type UserDto struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

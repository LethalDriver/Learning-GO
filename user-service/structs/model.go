package structs

type UserEntity struct {
	Id       string `bson:"id" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

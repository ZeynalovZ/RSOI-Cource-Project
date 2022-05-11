package models

type Register struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
}

type Auth struct {
	Id       string `json:"id" bson:"id"`
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
}
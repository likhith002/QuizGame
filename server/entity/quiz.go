package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name"`
	Questions []QuizQuestion     `json:"questions"`
}

type QuizQuestion struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Choices []QuizChoice `json:"choices"`
}

type QuizChoice struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}

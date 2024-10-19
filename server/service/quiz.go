package service

import (
	"og_ed/entity"
	"og_ed/internal/collection"
)

type QuizService struct {
	quizCollection *collection.QuizCollection
}

func Quiz(quizCollection *collection.QuizCollection) *QuizService {

	return &QuizService{
		quizCollection: quizCollection,
	}
}

func (quizService QuizService) GetQuizzes() ([]entity.Quiz, error) {

	return quizService.quizCollection.GetQuizzes()

}

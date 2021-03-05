package service

import (
	"awesomeProject/repository"
)

type Service struct {
	concurrentParsers int
	database          *repository.Service
}

func NewService(concurrentParsers, dbBufferSize int) (*Service, error) {
	dbService, err := repository.NewService(dbBufferSize)
	if err != nil {
		return nil, err
	}

	return &Service{
		concurrentParsers: concurrentParsers,
		database:          dbService}, nil
}

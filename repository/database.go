package repository

import (
	"log"
	"time"
)

type DBClient interface {
	initialize() error
	insertMovieBatch([]*Movie, chan error)
}

type Service struct {
	dbClient   DBClient
	bufferSize int
}

func NewService(bufferSize int) (*Service, error) {
	// TODO: improve dependency injection
	mongoClient := &MongoClient{}
	err := mongoClient.initialize()
	if err != nil {
		return nil, err
	}

	return &Service{
		dbClient:   mongoClient,
		bufferSize: bufferSize}, nil
}

func (s *Service) BufferedInsert(input <-chan *Movie) chan error {
	batchChannel := make(chan []*Movie, 5)
	errorChannel := make(chan error)

	go s.insert(batchChannel, errorChannel)

	go func() {
		count := s.bufferSize - 1
		buffer := make([]*Movie, s.bufferSize)
		for {
			var tempOutput chan []*Movie
			var tempInput <-chan *Movie
			if count < 0 {
				tempOutput = batchChannel
			} else {
				tempInput = input
			}
			select {
			case movie := <-tempInput: // tempInput is nil: disable the case
				buffer[count] = movie
				count--
			case tempOutput <- buffer: // tempOutput is nil: disable the case
				count = s.bufferSize - 1
				buffer = make([]*Movie, s.bufferSize)
			case <-time.After(time.Millisecond * 50):
				if len(buffer) > 0 {
					batchChannel <- buffer
				}
				close(batchChannel)

				log.Print("DONE INSERTING")
				return
			}
		}
	}()

	return errorChannel
}

func (s *Service) insert(movies chan []*Movie, errorChan chan error) {
	for batch := range movies {
		s.dbClient.insertMovieBatch(batch, errorChan)
	}

	close(errorChan)
}

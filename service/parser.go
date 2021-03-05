package service

import (
	"awesomeProject/repository"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (s *Service) Process(reader io.Reader) error {
	rawChannel := make(chan []string)

	go s.digester(reader, rawChannel)
	parserChannels := make([]<-chan *repository.Movie, s.concurrentParsers)
	for i := 0; i < s.concurrentParsers; i++ {
		parserChannels[i] = s.parser(rawChannel)
	}

	errorChannel := s.database.BufferedInsert(s.fanIn(parserChannels))
	for err := range errorChannel {
		if err != nil {
			log.Printf("error inserting: %v", err)
			close(rawChannel)
			return err
		}
	}

	return nil
}

func (s *Service) digester(reader io.Reader, output chan []string) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("interrupted reading")
		}
	}()

	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1
	for {
		line, err := r.Read()
		if err != nil && err != io.EOF {
			log.Printf("error reading line: %v", err)
			continue
		}

		if err == io.EOF {
			break
		}

		if len(line) != 24 {
			log.Printf("error, cannot process record with %d values, expected 24", len(line))
			continue
		}

		if line != nil {
			output <- line
		}
	}

	log.Print("DONE reading...")
	close(output)
}

func specialUnmarshal(raw string, in interface{}) error {
	if len(raw) == 0 {
		return nil
	}

	processed := raw
	r := regexp.MustCompile(`(".*)(')(.*")`)
	for r.MatchString(processed) {
		processed = r.ReplaceAllString(processed, "$1$3")
	}

	r = regexp.MustCompile(`\B'|'\B`)

	processed = r.ReplaceAllString(processed, `"`)
	processed = strings.Replace(processed, ": None", `: ""`, -1)

	return json.Unmarshal([]byte(processed), in)
}

func MovieFromRecord(record []string) *repository.Movie {
	movie := &repository.Movie{
		Adult:            record[0] == "True",
		Budget:           StrToIntPtr(record[2]),
		Homepage:         &record[4],
		Id:               StrToIntPtr(record[5]),
		ImdbId:           &record[6],
		OriginalLanguage: &record[7],
		OriginalTitle:    &record[8],
		Overview:         &record[9],
		Popularity:       StrToFloat32Ptr(record[10]),
		PosterPath:       &record[11],
		ReleaseDate:      &record[14],
		Revenue:          StrToIntPtr(record[15]),
		Runtime:          StrToFloat32Ptr(record[16]),
		Status:           &record[18],
		Tagline:          &record[19],
		Title:            &record[20],
		Video:            record[21] == "True",
		VoteAverage:      StrToFloat32Ptr(record[22]),
		VoteCount:        StrToIntPtr(record[23])}

	var err error
	collection := &repository.Collection{}
	err = specialUnmarshal(record[1], collection)
	if err != nil {
		log.Printf("error unmarshalling 'BelongsToCollection' field: %v", err)
	} else {
		movie.BelongsToCollection = collection
	}

	genres := &repository.Genres{}
	err = specialUnmarshal(record[3], genres)
	if err != nil {
		log.Printf("error unmarshalling 'Genres' field: %v", err)
	} else {
		movie.Genres = genres
	}

	companies := &repository.Companies{}
	err = specialUnmarshal(record[12], companies)
	if err != nil {
		log.Printf("error unmarshalling 'ProductionCompanies' field: %v", err)
	} else {
		movie.ProductionCompanies = companies
	}

	countries := &repository.Countries{}
	err = specialUnmarshal(record[13], countries)
	if err != nil {
		log.Printf("error unmarshalling 'ProductionCountries' field: %v", err)
	} else {
		movie.ProductionCountries = countries
	}

	languages := &repository.Languages{}
	err = specialUnmarshal(record[17], languages)
	if err != nil {
		log.Printf("error unmarshalling 'SpokenLanguages' field: %v", err)
	} else {
		movie.SpokenLanguages = languages
	}

	return movie
}

func (s *Service) parser(input <-chan []string) chan *repository.Movie {
	c := make(chan *repository.Movie)
	go func() {
		for record := range input {
			c <- MovieFromRecord(record)
		}
		close(c)
	}()

	return c
}

func (s *Service) fanIn(input []<-chan *repository.Movie) chan *repository.Movie {
	mainChannel := make(chan *repository.Movie)
	for _, channel := range input {
		go func(channel <-chan *repository.Movie) {
			for movie := range channel {
				mainChannel <- movie
			}
		}(channel)
	}

	return mainChannel
}

func StrToIntPtr(raw string) *int {
	val, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}

	return &val
}

func StrToFloat32Ptr(raw string) *float32 {
	val, err := strconv.ParseFloat(raw, 32)
	if err != nil {
		return nil
	}

	fValue := float32(val)
	return &fValue
}

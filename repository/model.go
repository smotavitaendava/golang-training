package repository

type (
	Movie struct {
		ID                  string      `bson:"_id,omitempty"`
		Adult               bool        `bson:"adult"`
		BelongsToCollection *Collection `bson:"belongs_to_collection"`
		Budget              *int        `bson:"budget"`
		Genres              *Genres     `bson:"genres"`
		Homepage            *string     `bson:"homepage"`
		Id                  *int        `bson:"imported_id"`
		ImdbId              *string     `bson:"imbd_id"`
		OriginalLanguage    *string     `bson:"original_language"`
		OriginalTitle       *string     `bson:"original_title"`
		Overview            *string     `bson:"overview"`
		Popularity          *float32    `bson:"popularity"`
		PosterPath          *string     `bson:"poster_path"`
		ProductionCompanies *Companies  `bson:"prod_companies"`
		ProductionCountries *Countries  `bson:"prod_countries"`
		ReleaseDate         *string     `bson:"release_date"`
		Revenue             *int        `bson:"revenue"`
		Runtime             *float32    `bson:"runtime"`
		SpokenLanguages     *Languages  `bson:"spoken_languages"`
		Status              *string     `bson:"status"`
		Tagline             *string     `bson:"tag_line"`
		Title               *string     `bson:"title"`
		Video               bool        `bson:"video"`
		VoteAverage         *float32    `bson:"vote_average"`
		VoteCount           *int        `bson:"vote_count"`
	}

	Collection struct {
		ID           int    `json:"id" bson:"collection_id"`
		Name         string `json:"name" bson:"name"`
		PosterPath   string `json:"poster_path" bson:"poster_path"`
		BackdropPath string `json:"backdrop_path" bson:"backdrop_path"`
	}

	Genre struct {
		ID   int    `json:"id" bson:"genre_id"`
		Name string `json:"name" bson:"name"`
	}

	Genres []Genre

	Company struct {
		ID   int    `json:"id" bson:"company_id"`
		Name string `json:"name" bson:"name"`
	}

	Companies []*Company

	Country struct {
		ISOCode string `json:"iso_3166_1" bson:"iso_3166_1_code"`
		Name    string `json:"name" bson:"name"`
	}

	Countries []*Country

	Language struct {
		ISOCode string `json:"iso_639_1" bson:"iso_639_1_code"`
		Name    string `json:"name" bson:"name"`
	}

	Languages []*Language
)

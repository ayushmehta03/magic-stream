package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct{
	GenreID int
	GenreName string
}

type Ranking struct{
 RankingValue int
 RankingName string
}




type Movie struct{
	ID bson.ObjectID `bson:"_id" json:"_id"`
	ImdbID string  `bson:"imdb_id" json:"imdb_id"`
	Title string  `bson:"title" json:"title"`
	PosterPath string `bson:"poster_path" json:"poster_path"`
	YoutubeID string  `bson:"youtube_id" json:"youtube_id"`
	Genre []Genre `bson: "genre" json:"genre"`
	AdminReview string `bson:"admin_review" json:"admin_review"`
	Ranking Ranking  `bson:"ranking" json:"ranking"`

}
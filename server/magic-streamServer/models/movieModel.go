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
	ID bson.ObjectID 
	ImdbID string 
	Title string 
	PosterPath string
	YoutubeID string
	Genre []Genre
	AdminReview string
	Ranking Ranking

}
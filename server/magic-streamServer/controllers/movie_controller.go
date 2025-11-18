package controllers

import (
	"context"
	"net/http"
	"time"

	database "github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
	models "github.com/ayushmehta03/magic-stream/server/magic-streamServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/go-playground/validator/v10"
)

var movieCollection *mongo.Collection=database.OpenCollection("movies")


// create validator object
var validate =validator.New()

// get all the movies present inside the database


func GetMovies() gin.HandlerFunc{
	return func (c *gin.Context){
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)


		defer cancel()

		var movies [] models.Movie

		cursor,err:=movieCollection.Find(ctx,bson.M{})
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch movies"})


			
		}
		defer cursor.Close(ctx)

			if err=cursor.All(ctx,&movies);err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to decode movies"})
			}

		c.JSON(http.StatusOK,movies)
	}
}

// get a single movie details with params as imdb id

func GetMovie() gin.HandlerFunc {
	return func (c *gin.Context){

		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()


		movieId:=c.Param("imdb_id")

		if movieId==""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"movie id is required"})
			return 
		}

		var movie models.Movie

		err:=movieCollection.FindOne(ctx,bson.M{"imdb_id":movieId}).Decode(&movie)

		if err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"movie not found"})
			return 
		}

		c.JSON(http.StatusOK,movie)






	}
}



// add movie to the database

func AddMovie() gin.HandlerFunc{
	return func(c *gin.Context){
	
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel();

		var movie models.Movie

		if err:=c.ShouldBindJSON(&movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Inavlid input"})
			return
		}

		if err:=validate.Struct(movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Validation failed","details":err.Error()})
			return 
		}


		result,err:=movieCollection.InsertOne(ctx,movie);

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to add movie"})
			return 
		}



		c.JSON(http.StatusCreated,result)










	}
}


// admin review langchain go


func AdminReviewUpdate() gin.HandlerFunc{
	return func(c *gin.Context){
		movieId:=c.Param("imdb_id")


		if movieId==""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Movie id required"})
			return
		}

		var req struct{
			AdminReview string `json:"admin_review"`

		}

		var resp struct{
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}


		if err:=c.ShouldBind(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request body"})
		}


	}
}
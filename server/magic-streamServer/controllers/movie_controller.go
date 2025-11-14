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
)

var movieCollection *mongo.Collection=database.OpenCollection("movies")


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
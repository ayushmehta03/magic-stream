package routes

import (
	controllers "github.com/ayushmehta03/magic-stream/server/magic-streamServer/controllers"
	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)



func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client){
	router.Use(middleware.AuthMiddlware())


	router.GET("/movie/:imdb_id",controllers.GetMovie(client))

	router.POST("/addmovie",controllers.AddMovie(client))

	router.GET("/recommendedmovies",controllers.GetRecommendedMovies(client))

	router.PATCH("/updatereview/:imdb_id",controllers.AdminReviewUpdate(client))


	
}




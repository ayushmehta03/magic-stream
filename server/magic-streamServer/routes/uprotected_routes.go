package routes

import (
	controllers "github.com/ayushmehta03/magic-stream/server/magic-streamServer/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)



func SetupUnProtectedRoutes(router *gin.Engine,client *mongo.Client){

router.GET("/movies",controllers.GetMovies(client))

		
	router.POST("/register",controllers.RegisterUser(client))

	router.POST("/login",controllers.LogInUser(client))

	router.POST("logout",controllers.LogoutHandler(client))
}





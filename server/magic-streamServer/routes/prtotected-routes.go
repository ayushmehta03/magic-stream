package routes


import(
	controllers "github.com/ayushmehta03/magic-stream/server/magic-streamServer/controllers"
	 "github.com/ayushmehta03/magic-stream/server/magic-streamServer/middleware"
	 	"github.com/gin-gonic/gin"


)



func SetupProtectedRoutes(router *gin.Engine){
	router.Use(middleware.AuthMiddlware())


	router.GET("/movie/:imdb_id",controllers.GetMovie())

	router.POST("/addmovie",controllers.AddMovie())

	router.GET("/recommendedmovies",controllers.GetRecommendedMovies())

	router.PATCH("/updatereview/:imdb_id",controllers.AdminReviewUpdate())


	
}




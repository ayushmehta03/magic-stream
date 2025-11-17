package routes


import(
	controllers "github.com/ayushmehta03/magic-stream/server/magic-streamServer/controllers"
	 	"github.com/gin-gonic/gin"


)



func SetupUnProtectedRoutes(router *gin.Engine){

router.GET("/movies",controllers.GetMovies())

		
	router.POST("/register",controllers.RegisterUser())

	router.POST("/login",controllers.LogInUser())

}





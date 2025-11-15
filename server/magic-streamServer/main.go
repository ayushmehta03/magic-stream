package main

import (
	"fmt"
	controllers "github.com/ayushmehta03/magic-stream/server/magic-streamServer/controllers"
	"github.com/gin-gonic/gin"
)

func main(){
	// main function


	router:=gin.Default()



	router.GET("/hello",func(c *gin.Context){
		c.String(200,"hello, magic-stream-movies")
	})

	router.GET("/movies",controllers.GetMovies())

		router.GET("/movie/:imdb_id",controllers.GetMovie())

	router.POST("/addmovie",controllers.AddMovie())


	if err:=router.Run(":8080");err!=nil{
		fmt.Println("failed to start server",err)
	}
	
}
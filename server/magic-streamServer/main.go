package main

import (
	"fmt"
		"github.com/ayushmehta03/magic-stream/server/magic-streamServer/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	// main function


	router:=gin.Default()



	router.GET("/hello",func(c *gin.Context){
		c.String(200,"hello, magic-stream-movies")
	})

	routes.SetupUnProtectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	

	if err:=router.Run(":8080");err!=nil{
		fmt.Println("failed to start server",err)
	}
	
}
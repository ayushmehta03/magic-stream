package main

import (
	"fmt"
		"github.com/ayushmehta03/magic-stream/server/magic-streamServer/routes"
		"go.mongodb.org/mongo-driver/v2/mongo"
			database "github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"


	"github.com/gin-gonic/gin"

)

func main(){
	// main function


	router:=gin.Default()



	router.GET("/hello",func(c *gin.Context){
		c.String(200,"hello, magic-stream-movies")
	})

	 var client *mongo.Client=database.Connect()


	routes.SetupUnProtectedRoutes(router,client)
	routes.SetupProtectedRoutes(router,client)

	

	if err:=router.Run(":8080");err!=nil{
		fmt.Println("failed to start server",err)
	}
	
}
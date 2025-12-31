package main

import (
	"context"
	"fmt"
	"log"

	database "github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/routes"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
)

func main(){
	// main function


	router:=gin.Default()


	router.GET("/hello",func(c *gin.Context){
		c.String(200,"hello, magic-stream-movies")
	})

	err:=godotenv.Load(".env")
	if err!=nil{
		log.Println("Warning: unable to find .env file")
	}

	 var client *mongo.Client=database.Connect()

	if err:=client.Ping(context.Background(),nil);err!=nil{
		log.Fatalf("Failed to reach server: %v",err)
	}

	defer func(){
		err:=client.Disconnect(context.Background())
		if err!=nil{
			log.Fatalf("Failed to disconnect from MongoDB %v",err)
		}
	}()




	routes.SetupUnProtectedRoutes(router,client)
	routes.SetupProtectedRoutes(router,client)

	

	if err:=router.Run(":8080");err!=nil{
		fmt.Println("failed to start server",err)
	}
	
}
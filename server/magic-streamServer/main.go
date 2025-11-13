package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main(){
	// main function


	router:=gin.Default()



	router.GET("/hello",func(c *gin.Context){
		c.String(200,"hello, magic-stream-movies")
	})

	if err:=router.Run(":8080");err!=nil{
		fmt.Println("failed to start server",err)
	}
	
}
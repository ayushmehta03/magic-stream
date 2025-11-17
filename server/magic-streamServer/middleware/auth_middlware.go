package middleware

import (
	"net/http"

	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/utils"
	"github.com/gin-gonic/gin"
)


func AuthMiddlware() gin.HandlerFunc{
 return func(c *gin.Context){

	token,err:=utils.GetAccessToken(c)

	if err!=nil{
		c.JSON(http.StatusUnauthorized,gin.H{"error":err.Error()})
		c.Abort()
	}

	if token==""{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"No token provided "})
		c.Abort()
	}



 }
}
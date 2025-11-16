package controllers

import (
	"context"
	"net/http"
	"time"

	database "github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
	models "github.com/ayushmehta03/magic-stream/server/magic-streamServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
		"go.mongodb.org/mongo-driver/v2/bson"


)

var userCollection *mongo.Collection= database.OpenCollection("users")




func Hashpassword(password string) (string,error){

	HashPassword,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	if err!=nil{
		return "",err
	}

	return string(HashPassword),nil
}


func RegisterUser() gin.HandlerFunc{
	return func(c *gin.Context){


		var user models.User


		if err:=c.ShouldBindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input data"})
			return
		}

		validate:=validator.New();

		if err:=validate.Struct(user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Validation failed","details":err.Error()})
			return
		}


		hashedPassword,err:=Hashpassword(user.Password)

		



		if err!=nil{

			c.JSON(http.StatusBadRequest,gin.H{"error":"Internal server error"});
			return

		}


		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
				defer cancel()


		count,err:=userCollection.CountDocuments(ctx,bson.M{"email":user.Email})

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to check existing user"})
			return
		}

		if count>0{
			c.JSON(http.StatusConflict,gin.H{"error":"user already exists"})
			return
		}

		user.UserId=bson.NewObjectID().Hex()

		user.CreatedAt=time.Now()
		user.UpdatedAt=time.Now()
		user.Password=hashedPassword



		result,err:=userCollection.InsertOne(ctx,user)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to register "})
			return
		}


		c.JSON(http.StatusCreated,result)


















	}
}



func LogInUser() gin.HandlerFunc{
	return func(c *gin.Context){

		var userLogIn models.UserLogin

		if err:=c.ShouldBindJSON(&userLogIn);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid input data"})
			return
		}

		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()

		var foundUser models.User

		err:=userCollection.FindOne(ctx,bson.M{"email":userLogIn.Email}).Decode(&foundUser)

		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Inavlid email or password"})
			return
		}


		err=bcrypt.CompareHashAndPassword([]byte(foundUser.Password),[]byte(userLogIn.Password))

		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid email or password"})
			return
		}

		






	}
}

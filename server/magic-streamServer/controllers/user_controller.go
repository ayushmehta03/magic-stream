package controllers

import (
	"context"
	"net/http"
	"time"
	"fmt"

	database "github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
	models "github.com/ayushmehta03/magic-stream/server/magic-streamServer/models"
	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
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


		token,refreshToken,err:=utils.GenerateAllTokens(foundUser.Email,foundUser.FirstName,foundUser.LastName,foundUser.Role,foundUser.UserId)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to generate tokens"})
			return
		}

		err=utils.UpdateAllTokens(foundUser.UserId,token,refreshToken)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update tokens"})
			return
		}

		c.JSON(http.StatusOK,models.UserResponse{
			UserId:foundUser.UserId,
			FirstName: foundUser.FirstName,
			LastName: foundUser.LastName,
			Email: foundUser.Email,
			Role:foundUser.Role,
			FavouriteGenres: foundUser.FavouriteGenres,
			Token: token,
			RefreshToken: refreshToken,

		})




	}
}
func LogoutHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the access_token cookie

		var UserLogout struct {
			UserId string `json:"user_id"`
		}

		err := c.ShouldBindJSON(&UserLogout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Println("User ID from Logout request:", UserLogout.UserId)

		err = utils.UpdateAllTokens(UserLogout.UserId, "", "", client) // Clear tokens in the database
		// Optionally, you can also remove the user session from the database if needed

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
			return
		}
		// c.SetCookie(
		// 	"access_token",
		// 	"",
		// 	-1, // MaxAge negative → delete immediately
		// 	"/",
		// 	"localhost", // Adjust to your domain
		// 	true,        // Use true in production with HTTPS
		// 	true,        // HttpOnly
		// )
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: "",
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		// // Clear the refresh_token cookie
		// c.SetCookie(
		// 	"refresh_token",
		// 	"",
		// 	-1,
		// 	"/",
		// 	"localhost",
		// 	true,
		// 	true,
		// )
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct{
	Email string
	FirstName string
	LastName string
	Role string
	UserId string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection=database.OpenCollection("users")


var SECRET_KEY string=os.Getenv("SECRET_KEY")
 var SECRET_REFRESH_KEY string=os.Getenv("SECRET_REFRESH_KEY")

func GenerateAllTokens(email,firstName,lastName,role,userId string)(string,string,error){
	claims:=&SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "magic-stream",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),

		},
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	signedToken,err:=token.SignedString([]byte(SECRET_KEY))

	if err!=nil{
		return "","",err
	}


	// refresh token 
	refreshClaims:=&SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "magic-stream",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*7*time.Hour)),

		},
	}

	refreshToken:=jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims)

	refreshSignedToken,err:=refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	if err!=nil{
		return "","",err
	}


	return signedToken,refreshSignedToken,nil





}



func UpdateAllTokens(userId,token,refreshToken string) (err error){
	var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
	defer cancel()

	updatedAt,_:=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))


	updateData:=bson.M{
		"$set":bson.M{
			"token":token,
			"refresh_token":refreshToken,
			"update_at":updatedAt,
		},
	}


	_,err=userCollection.UpdateOne(ctx,bson.M{"user_id":userId},updateData)

	if err!=nil{
		return err
	}

	return nil
}


func GetAccessToken(c *gin.Context)(string,error){


	authHeader:=c.Request.Header.Get("Authorization")

	if authHeader==""{
		return "",errors.New("Authorization header is required")
	}

tokenString := authHeader[len("Bearer "):]


	if tokenString==""{
		return "",errors.New("Bearer token is required ")
	}

	return tokenString,nil
}

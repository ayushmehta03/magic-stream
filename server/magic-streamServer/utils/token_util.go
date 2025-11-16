package utils

import (
	"os"
	"time"

	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),

		},
	}

	refreshToken:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	refreshSignedToken,err:=token.SignedString([]byte(SECRET_KEY))

	if err!=nil{
		return "","",err
	}

	



}


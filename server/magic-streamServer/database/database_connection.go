package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// connecting the db 
func DBConnect() *mongo.Client{
	err:=godotenv.Load(".env")
	if err!=nil{
		log.Println("warning: unable to find .env file")
	}

	MongoDb:=os.Getenv("MONGODB_URI")

	if MongoDb==""{
		log.Fatal("MONGODB_URI not set !")
	}

	fmt.Println("MongoDB URI: ",MongoDb)




	clientOptions:=options.Client().ApplyURI(MongoDb)

		client,err:=mongo.Connect(clientOptions)

		if err!=nil{
			return nil
		}

		return client

}


 var Client *mongo.Client=DBConnect()


	// function for running all the queries 

func OpenCollection(collectionName string) *mongo.Collection{

	err:=godotenv.Load(".env")

	if err!=nil{
		log.Println("warning: unable to find .env file")
	}


	databaseName:=os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME: ",databaseName)

	collection:=Client.Database(databaseName).Collection(collectionName)

	if collection==nil{
		return nil
	}

	return collection
}




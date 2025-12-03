package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	database "github.com/ayushmehta03/magic-stream/server/magic-streamServer/database"
	models "github.com/ayushmehta03/magic-stream/server/magic-streamServer/models"
	"github.com/ayushmehta03/magic-stream/server/magic-streamServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection=database.OpenCollection("movies")
var rankingCollection *mongo.Collection=database.OpenCollection("rankings")


// create validator object
var validate =validator.New()

// get all the movies present inside the database


func GetMovies() gin.HandlerFunc{
	return func (c *gin.Context){
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)


		defer cancel()

		var movies [] models.Movie

		cursor,err:=movieCollection.Find(ctx,bson.M{})
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch movies"})


			
		}
		defer cursor.Close(ctx)

			if err=cursor.All(ctx,&movies);err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to decode movies"})
			}

		c.JSON(http.StatusOK,movies)
	}
}

// get a single movie details with params as imdb id

func GetMovie() gin.HandlerFunc {
	return func (c *gin.Context){

		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()


		movieId:=c.Param("imdb_id")

		if movieId==""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"movie id is required"})
			return 
		}

		var movie models.Movie

		err:=movieCollection.FindOne(ctx,bson.M{"imdb_id":movieId}).Decode(&movie)

		if err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"movie not found"})
			return 
		}

		c.JSON(http.StatusOK,movie)






	}
}



// add movie to the database

func AddMovie() gin.HandlerFunc{
	return func(c *gin.Context){
	
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel();

		var movie models.Movie

		if err:=c.ShouldBindJSON(&movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Inavlid input"})
			return
		}

		if err:=validate.Struct(movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Validation failed","details":err.Error()})
			return 
		}


		result,err:=movieCollection.InsertOne(ctx,movie);

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to add movie"})
			return 
		}



		c.JSON(http.StatusCreated,result)










	}
}


// admin review langchain go


func AdminReviewUpdate() gin.HandlerFunc{
	return func(c *gin.Context){
		movieId:=c.Param("imdb_id")


		if movieId==""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Movie id required"})
			return
		}

		var req struct{
			AdminReview string `json:"admin_review"`

		}

		var resp struct{
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}


		if err:=c.ShouldBind(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request body"})
			return 
		}


		sentiment,rankVal,err:=GetReviewRanking(req.AdminReview)


		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error getting the admin review"})
		}


		filter:=bson.M{"imdb_id":movieId}

		update:=bson.M{
			"$set":bson.M{
				"admin_review":req.AdminReview,
				"ranking":bson.M{
					"ranking_value":rankVal,
					"ranking_name":sentiment,
				},
			},
		}


		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		result,err:=movieCollection.UpdateOne(ctx,filter,update)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error updating movie"})
			return
		}

		if result.MatchedCount==0{
			c.JSON(http.StatusNotFound,gin.H{"error":"Movie not found"})
			return 
		}

		resp.RankingName=sentiment
		resp.AdminReview=req.AdminReview

		c.JSON(http.StatusOK,resp)

		



	}
}


func GetReviewRanking(admin_review string)(string,int,error){

	rankings,err:=GetRankings()
	if err!=nil{
		return "",0,err
	}


	sentimentDelimited:="";


	for _,ranking:=range rankings{
		if ranking.RankingValue!=999{
			sentimentDelimited=sentimentDelimited+ranking.RankingName +","
		}
	}

	sentimentDelimited=strings.Trim(sentimentDelimited,",")



	err=godotenv.Load(".env")

	if err!=nil{
		log.Println("Warning: .env file not found")
	}


	OpenAiApiKey:=os.Getenv("OPENAI_API_KEY")


	if OpenAiApiKey==""{
		return "",0,errors.New("could not read OPENAI_API_KEY")
	}



	llm,err:=openai.New(openai.WithToken(OpenAiApiKey))

	if err!=nil{
		return "",0,err;
	}



	base_prompt_temp:=os.Getenv("BASE_PROMPT_TEMPLATE")

	base_prompt:=strings.Replace(base_prompt_temp,"{rankings}",sentimentDelimited,1)

	response,err:=llm.Call(context.Background(),base_prompt+admin_review)

	if err!=nil{
		return "",0,err
	}

	rankValue:=0

	for _, ranking :=range rankings{
		if ranking.RankingName==response{
			rankValue=ranking.RankingValue
			break
		}
	}

	return  response,rankValue,nil

}


func GetRankings()([]models.Ranking,error){

	var rankings[] models.Ranking


	var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
	
	defer cancel()
	cursor,err:=rankingCollection.Find(ctx,bson.M{})

	if err!=nil{
		return nil,err
	}

	defer cursor.Close(ctx)


	if err:=cursor.All(ctx,&rankings);err!=nil{
		return nil,err
	}

	return rankings,nil

}

func GetRecommendedMovies() gin.HandlerFunc{
	return func(c *gin.Context){

	userId,err:=utils.GetUserIdFromContext(c)

	if err!=nil{
 	c.JSON(http.StatusBadRequest,gin.H{"error":"user id not found in context"})
		return 
	}

	




	}
}
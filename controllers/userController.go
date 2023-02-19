package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-authentication/database"
	"go-authentication/helpers"
	"go-authentication/models"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

const timeLayout = time.RFC3339

var timeString = time.Now().Format(timeLayout)

func HashPassword()

func VerifyPassword()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationError := validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": " error occured while checking for the user email",
			})
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while checking for the user phone",
			})
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "this email or phone number exited",
			})
		}
		user.CreatedAt, _ = time.Parse(timeLayout, timeString)
		user.UpdatedAt, _ = time.Parse(timeLayout, timeString)
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.UserType, user.UserID)
		user.Token = token
		user.RefreshToken = refreshToken

		//insert to database
		resultInsertNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "User item was not created",
			})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertNumber)
	}
}

func Login()

func GetUsers()

func GetUser() gin.HandlerFunc {

	return func(c *gin.Context) {
		//find by user_id
		userId := c.Param("user_id")

		//verify user is admin or not
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}
}

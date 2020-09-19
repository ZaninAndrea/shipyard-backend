package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nbutton23/zxcvbn-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// load .env file
	godotenv.Load()

	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("CONNECTION_URI")))
	if err != nil {
		panic(err)
	}

	// disconnect on quit
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	userCollection := client.Database("production").Collection("users")
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		// read query parameters
		email, providedEmail := c.Request.URL.Query()["email"]
		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedEmail || !providedPassword {
			c.JSON(400, gin.H{
				"error": "You need to pass an email and password in the query",
			})
			return
		}

		// check if user exists
		filter := bson.M{"email": email[0]}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, filter)

		if err != nil {
			panic(err)
		} else if count == 0 {
			c.JSON(400, gin.H{"error": "User with this email not found"})
			return
		}

		// load user
		userFound := loadUserByEmail(email[0], userCollection)

		// check that the password is correct
		hash := userFound.Password
		match := CheckPasswordHash(password[0], hash)
		if !match {
			c.JSON(500, gin.H{
				"error": "Wrong password",
			})
			return
		}

		c.JSON(200, gin.H{
			"token": GenerateToken(userFound.ID.Hex()),
		})
	})

	r.GET("/user", func(c *gin.Context) {
		// check authentication
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		userFound := loadUserByID(parsedToken.UserID, userCollection)

		// parse the bson data into JSON saved as []byte
		jsonBytes, err := bson.MarshalExtJSON(userFound.Data, true, true)

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, string(jsonBytes))
	})

	r.POST("/user", func(c *gin.Context) {
		email, providedEmail := c.Request.URL.Query()["email"]
		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedEmail || !providedPassword {
			c.JSON(400, gin.H{
				"error": "You need to pass an email and password in the query",
			})
			return
		}

		passwordStrength := zxcvbn.PasswordStrength(password[0], []string{email[0]})
		if passwordStrength.Score < 2 {
			c.JSON(400, gin.H{
				"error": "The password is too weak",
			})
			return
		}

		// check if a user with the same email exists
		filter := bson.M{"email": email[0]}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, filter)

		if err != nil {
			panic(err)
		} else if count > 0 {
			c.JSON(400, gin.H{"error": "Another user with this email already exists"})
			return
		}

		// parse json body to bson
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		var initialData interface{}
		err = bson.UnmarshalExtJSON(jsonData, true, &initialData)

		// create new user
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		passwordHash, err := HashPassword(password[0])

		if err != nil {
			panic(err)
		}

		res, err := userCollection.InsertOne(ctx, bson.M{"email": email[0], "password": passwordHash, "data": initialData})
		id := res.InsertedID.(primitive.ObjectID).Hex()

		c.JSON(200, gin.H{
			"token": GenerateToken(id),
		})
	})

	r.PUT("/user", func(c *gin.Context) {
		// Check authorization
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		// parse json body to bson
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		var updateQuery interface{}
		err = bson.UnmarshalExtJSON(jsonData, true, &updateQuery)

		// update database
		objID, _ := primitive.ObjectIDFromHex(parsedToken.UserID)
		filter := bson.M{"_id": objID}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		update := bson.M{"$set": bson.M{"data": updateQuery}}

		err = userCollection.FindOneAndUpdate(ctx, filter, update).Err()
		if err != nil {
			panic(err)
		}

		c.String(200, "")
	})

	r.DELETE("/user", func(c *gin.Context) {
		// check authentication
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		// update database
		objID, _ := primitive.ObjectIDFromHex(parsedToken.UserID)
		filter := bson.M{"_id": objID}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = userCollection.FindOneAndDelete(ctx, filter).Err()
		if err != nil {
			panic(err)
		}

		c.String(200, "")
	})

	r.Run()
}

package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	jsonpatchtomongo "github.com/ZaninAndrea/json-patch-to-mongo"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
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

	// create server and allow CORS from all origins
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(location.Default())

	r.POST("/login", func(c *gin.Context) {
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

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
			c.JSON(400, gin.H{
				"error": "Wrong password",
			})
			return
		}

		c.JSON(200, gin.H{
			"token": GenerateToken(userFound.ID.Hex()),
		})
	})

	r.POST("/changePassword", func(c *gin.Context) {
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

		email, providedEmail := c.Request.URL.Query()["email"]
		oldPassword, providedOldPassword := c.Request.URL.Query()["oldPassword"]
		newPassword, providedNewPassword := c.Request.URL.Query()["newPassword"]
		if !providedOldPassword || !providedNewPassword || !providedEmail {
			c.JSON(400, gin.H{
				"error": "You need to pass email, oldPassword and newPassword in the query",
			})
			return
		}
		passwordStrength := zxcvbn.PasswordStrength(newPassword[0], []string{email[0]})
		if passwordStrength.Score < 2 {
			c.JSON(400, gin.H{
				"error": "The new password is too weak",
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
		match := CheckPasswordHash(oldPassword[0], hash)
		if !match {
			c.JSON(400, gin.H{
				"error": "Wrong password",
			})
			return
		}

		passwordHash, err := HashPassword(newPassword[0])
		if err != nil {
			panic(err)
		}
		update := bson.M{"$set": bson.M{"password": passwordHash}}

		err = userCollection.FindOneAndUpdate(ctx, filter, update).Err()
		if err != nil {
			panic(err)
		}

		c.String(200, "")
		config.sendPasswordChangedEmail(email[0])
	})

	r.GET("/user", func(c *gin.Context) {
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

		// check authentication
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse authorization token"})
			return
		}

		var _projection bson.M
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			_projection = bson.M{}
		} else {
			err = bson.UnmarshalExtJSON(jsonData, true, &_projection)
		}

		projection := bson.M{}
		for k, v := range _projection {
			projection["data."+k] = v
		}

		userFound := loadUserByID(parsedToken.UserID, userCollection, projection)

		// parse the bson data into JSON saved as []byte
		jsonBytes, err := bson.MarshalExtJSON(userFound.Data, false, true)

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, string(jsonBytes))
	})

	r.POST("/user", func(c *gin.Context) {
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

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
			c.JSON(400, gin.H{"error": "Failed to parse body"})
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

	r.PATCH("/user", func(c *gin.Context) {
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

		// Check authorization
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse token"})
			return
		}

		rawPatch, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to read body"})
			return
		}
		_updateQuery, err := jsonpatchtomongo.ParsePatchesWithPrefix(rawPatch, "data.")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// apply the update query to the authenticated user
		objID, _ := primitive.ObjectIDFromHex(parsedToken.UserID)
		filter := bson.M{"_id": objID}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = userCollection.FindOneAndUpdate(ctx, filter, _updateQuery).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to apply patch to database"})
			return
		}

		c.String(200, "")
	})

	r.PUT("/user", func(c *gin.Context) {
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

		// Check authorization
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse authorization token"})
			return
		}

		// parse json body to bson
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to read body"})
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
		config, err := GetServerConfig(c, client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		userCollection := config.UserCollection

		// check authentication
		parsedToken, err := parseBearer(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse authorization token"})
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

	SetupAdminRoute(r, client)

	if os.Getenv("TLS_CERT_DIR") == "" {
		r.Run()
	} else {
		dir := os.Getenv("TLS_CERT_DIR")
		r.RunTLS(":8080", dir+"/fullchain.pem", dir+"/privkey.pem")
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupAdminRoute(r *gin.Engine, client *mongo.Client) {
	adminDomain := os.Getenv("ADMIN_DOMAIN")
	r.GET("/admin/checkPassword", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	r.GET("/admin/configs", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		configs, err := GetAllServerConfigs(client)
		if err != nil {
			panic(err)
		}

		jsonBytes, err := json.Marshal(configs)
		if err != nil {
			panic(err)
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, string(jsonBytes))
	})

	r.POST("/admin/configs", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		// parse json body to DatabaseConfig
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse body"})
			return
		}
		var configData DatabaseConfigNoID
		err = json.Unmarshal(jsonData, &configData)
		if err != nil {
			c.JSON(400, gin.H{"error": "Configuration passed is invalid"})
			return
		}

		// Check that the DatabaseConfig is valid
		if configData.Domain == "" {
			c.JSON(400, gin.H{"error": "You must pass a domain"})
			return
		}

		// check if a configuration with the same domain exists
		filter := bson.M{"domain": configData.Domain}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		count, err := client.Database("administration").Collection("servers").CountDocuments(ctx, filter)
		if err != nil {
			panic(err)
		} else if count > 0 {
			c.JSON(400, gin.H{"error": "Another server with this domain already exists"})
			return
		}

		// Create new server configuration
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Database("administration").Collection("servers").InsertOne(
			ctx,
			configData,
		)
		if err != nil {
			panic(err)
		}

		c.String(200, "")
	})

	r.PUT("/admin/configs", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		// parse json body to DatabaseConfig
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse body"})
			return
		}
		var configData DatabaseConfigNoInternals
		err = json.Unmarshal(jsonData, &configData)
		if err != nil {
			c.JSON(400, gin.H{"error": "Configuration passed is invalid"})
			return
		}

		// check if a configuration with the same domain exists
		filter := bson.M{"_id": configData.ID}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		count, err := client.Database("administration").Collection("servers").CountDocuments(ctx, filter)
		if err != nil {
			panic(err)
		} else if count == 0 {
			c.JSON(400, gin.H{"error": "No server exists with the passed id"})
			return
		}

		// Create new server configuration
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Database("administration").Collection("servers").UpdateOne(
			ctx,
			filter,
			bson.M{
				"$set": configData,
			},
		)
		if err != nil {
			panic(err)
		}

		c.String(200, "")
	})

	r.DELETE("/admin/configs/:id", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Passed an invalid id"})
			return
		}

		// check if a configuration with the same domain exists
		filter := bson.M{"_id": id}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.Database("administration").Collection("servers").DeleteOne(ctx, filter)
		fmt.Println((*res).DeletedCount)
		fmt.Println(filter)
		if err != nil {
			panic(err)
		}

		c.String(200, "")
	})

	r.GET("/admin/configs/:configId/users", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		_limit, providedLimit := c.Request.URL.Query()["limit"]
		var limit int64 = 30
		if providedLimit {
			i1, err := strconv.Atoi(_limit[0])
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Limit must be an integer",
				})
				return
			}
			if i1 < 30 {
				limit = int64(i1)
			}
		}
		_offset, providedOffset := c.Request.URL.Query()["offset"]
		var offset int64 = 0
		if providedOffset {
			i1, err := strconv.Atoi(_offset[0])
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Offset must be an integer",
				})
				return
			}
			if i1 > 0 {
				offset = int64(i1)
			}
		}

		id, err := primitive.ObjectIDFromHex(c.Param("configId"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Passed an invalid id"})
			return
		}

		// load the configuration
		filter := bson.M{"_id": id}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var config DatabaseConfig
		err = client.Database("administration").Collection("servers").FindOne(ctx, filter).Decode(&config)
		if err != nil {
			c.JSON(400, gin.H{"error": "Couldn't load the server configuration matching the passed id"})
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		options := options.Find()
		options.SetLimit(limit)
		options.SetSkip(offset)
		options.SetProjection(bson.M{"data": 0})
		cursor, err := client.Database("generic_"+config.ID.Hex()).Collection("users").Find(ctx, bson.D{}, options)
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not load the users"})
			return
		}
		users := make([]User, 0)
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = cursor.All(ctx, &users)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed parsing the users"})
			return
		}

		// Parse to JSON and return it
		jsonBytes, err := json.Marshal(users)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed parsing users to json"})
			fmt.Println(err)
			return
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, string(jsonBytes))
	})

	r.GET("/admin/configs/:configId/users/:userId", func(c *gin.Context) {
		url := location.Get(c)
		if url.Hostname() != adminDomain {
			c.JSON(400, gin.H{
				"error": "This route is not available",
			})
			return
		}

		password, providedPassword := c.Request.URL.Query()["password"]
		if !providedPassword {
			c.JSON(400, gin.H{
				"error": "You must specify the password query field",
			})
			return
		}
		if !CheckPasswordHash(password[0], "$2a$14$"+os.Getenv("ADMIN_PASSWORD_HASH")) {
			c.JSON(400, gin.H{
				"error": "Passed password is wrong",
			})
			return
		}

		configId, err := primitive.ObjectIDFromHex(c.Param("configId"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Passed an invalid configuration id"})
			return
		}

		userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Passed an invalid user id"})
			return
		}

		// load the configuration
		filter := bson.M{"_id": configId}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var config DatabaseConfig
		err = client.Database("administration").Collection("servers").FindOne(ctx, filter).Decode(&config)
		if err != nil {
			c.JSON(400, gin.H{"error": "Couldn't load the server configuration matching the passed id"})
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var result User
		err = client.
			Database("generic_"+config.ID.Hex()).
			Collection("users").
			FindOne(ctx, bson.M{"_id": userId}).Decode(&result)
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not load the user"})
			return
		}

		// Parse to JSON and return it
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not marshal the user to JSON"})
			return
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, string(jsonBytes))
	})

}

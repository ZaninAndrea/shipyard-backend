package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

}

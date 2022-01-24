package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	internal "github.com/ZaninAndrea/shipyard-backend/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupApiServer(client *mongo.Client) *gin.Engine {

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

	internal.SetupUserRoute(r, client)
	internal.SetupAdminRoute(r, client)
	return r
}

func SetupStaticServer() *gin.Engine {
	websitesPath := os.Getenv("WEBSITES_DIRECTORY")

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

	// Serve the requested file from disk
	r.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		if path == "" {
			path = "index.html"
		}

		url := location.Get(c)
		c.File(filepath.Join(websitesPath, url.Hostname(), c.Param("path")))
	})

	return r
}

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
	apiServer := SetupApiServer(client)
	staticServer := SetupStaticServer()

	var wg sync.WaitGroup
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		if os.Getenv("TLS_CERT_DIR") == "" {
			apiServer.Run(":8080")
		} else {
			dir := os.Getenv("TLS_CERT_DIR")
			apiServer.RunTLS(":8080", dir+"/fullchain.pem", dir+"/privkey.pem")
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		if os.Getenv("TLS_CERT_DIR") == "" {
			staticServer.Run(":80")
		} else {
			dir := os.Getenv("TLS_CERT_DIR")
			staticServer.RunTLS(":80", dir+"/fullchain.pem", dir+"/privkey.pem")
		}
	}(&wg)

	wg.Wait()
}

// func test(r *gin.Engine, client *mongo.Client) {
// 	objID, _ := primitive.ObjectIDFromHex("613b7135a161e522ea5d5575")
// 	filter := bson.M{"_id": objID}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_updateQuery, err := jsonpatchtomongo.ParsePatchesWithPrefix([]byte(`[
//   		{ "op": "replace", "path": "/abc", "value": "test" }
// 	]`), "data.")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	userCollection := client.Database("generic_613b7122a161e522ea5d5574").Collection("users")
// 	err = userCollection.FindOneAndUpdate(ctx, filter, bson.A{_updateQuery}).Err()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

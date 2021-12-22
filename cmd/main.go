package main

import (
	"context"
	"os"
	"time"

	internal "github.com/ZaninAndrea/shipyard-backend/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	internal.SetupUserRoute(r, client)
	internal.SetupAdminRoute(r, client)

	// test(r, client)
	if os.Getenv("TLS_CERT_DIR") == "" {
		r.Run()
	} else {
		dir := os.Getenv("TLS_CERT_DIR")
		r.RunTLS(":8080", dir+"/fullchain.pem", dir+"/privkey.pem")
	}
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

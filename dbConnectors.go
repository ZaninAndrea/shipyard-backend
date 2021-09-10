package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gomail "gopkg.in/gomail.v2"
)

// User is a representation of a document from the users collection in MongoDB
type User struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Email    string
	Password string
	Data     bson.M
}

func loadUserByEmail(email string, collection *mongo.Collection) User {
	filter := bson.M{"email": email}

	var result User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		panic(err)
	}

	return result
}
func loadUserByID(id string, collection *mongo.Collection, projection bson.M) User {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	var result User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)

	if err != nil {
		panic(err)
	}

	return result
}

type DatabaseConfig struct {
	ID             primitive.ObjectID `bson:"_id, omitempty"`
	Domain         string
	UserCollection *mongo.Collection `bson:"-" json:"-"`
	App            struct {
		Name        string
		LogoLink    string
		Link        string
		HeaderColor string
	}
	Company struct {
		Name    string
		Address string
	}
	Smtp struct {
		Server      string
		Port        int
		Username    string
		Password    string
		EmailDialer *gomail.Dialer `bson:"-" json:"-"`
	}
}
type DatabaseConfigNoID struct {
	Domain         string
	UserCollection *mongo.Collection `bson:"-" json:"-"`
	App            struct {
		Name        string
		LogoLink    string
		Link        string
		HeaderColor string
	}
	Company struct {
		Name    string
		Address string
	}
	Smtp struct {
		Server      string
		Port        int
		Username    string
		Password    string
		EmailDialer *gomail.Dialer `bson:"-" json:"-"`
	}
}
type DatabaseConfigNoInternals struct {
	ID     primitive.ObjectID `bson:"_id, omitempty"`
	Domain string
	App    struct {
		Name        string
		LogoLink    string
		Link        string
		HeaderColor string
	}
	Company struct {
		Name    string
		Address string
	}
	Smtp struct {
		Server   string
		Port     int
		Username string
		Password string
	}
}

func GetServerConfig(c *gin.Context, client *mongo.Client) (*DatabaseConfig, error) {
	url := location.Get(c)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var config DatabaseConfig
	err := client.Database("administration").Collection("servers").FindOne(ctx, bson.M{
		"domain": url.Hostname(),
	}).Decode(&config)

	if err != nil {
		return nil, fmt.Errorf("Could not find a server configuration associated to the domain %s", url.Hostname())
	}

	config.UserCollection = client.Database("generic_" + config.ID.Hex()).Collection("users")
	config.Smtp.EmailDialer = gomail.NewDialer(
		config.Smtp.Server, config.Smtp.Port, config.Smtp.Username, config.Smtp.Password,
	)

	return &config, nil
}

func GetAllServerConfigs(client *mongo.Client) ([]DatabaseConfig, error) {
	configs := []DatabaseConfig{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := client.Database("administration").Collection("servers").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	cursor.All(ctx, &configs)

	return configs, nil
}

package main

// 👈 Require the packages
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/amrilsyaifa/go_mongodb_rest_api/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 👈 Create required variables that we'll re-assign later
var (
	server 		*gin.Engine
	ctx			context.Context
	mongoClient *mongo.Client
	redisClient	*redis.Client
)

func init() {
	// 👇 Load the .env variables
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environtment variables", err)
	}

	/**
	 How difrent TODO and Background

	 	* Background
		Background returns a non-nil, empty Context. It is never canceled, has no
		values, and has no deadline. It is typically used by the main function,
		initialization, and tests, and as the top-level Context for incoming requests.
		eg: ctx = context.Background()

		TODO
		TODO returns a non-nil, empty Context. Code should use context.TODO when
		it's unclear which Context to use or it is not yet available (because the
		surrounding function has not yet been extended to accept a Context parameter).
		eg: ctx = context.TODO()
	*/
	// 👇 Create a context
	ctx = context.TODO()

	// 👇 Connect to MongoDB
	mongoConn := options.Client().ApplyURI((config.DBUri))
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		panic(err)
	}

	if err:= mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB succesfully connected...")


	// 👇 Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisClient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB",0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis client connected successfully...")

	// 👇 Create the Gin Engine instance
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	value, err := redisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	log.Fatal(server.Run(":" + config.Port))
}
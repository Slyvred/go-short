package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type shortenedUrl struct {
	ID             string    `json:"id"`
	Original       string    `json:"original"`
	Shortened      string    `json:"shortened"` // Will be used as a key
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	LastAccessedAt time.Time `json:"lastAccessedAt" bson:"lastAccessedAt"`
	AccessCount    uint      `json:"accessCount" bson:"accessCount"`
}

func connectToMongo() *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file") // Exit is handled below, because the MONGO_URI is set as a secret on fly (= not in dotenv)
	}

	connectionString := os.Getenv("MONGO_URI")
	if connectionString == "" {
		log.Fatal("Missing MONGO_URI environment variable")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB")
	return client.Database("slyvred").Collection("go-shorten")
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

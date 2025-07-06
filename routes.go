package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func hello(c *gin.Context) {

	body := `Welcome to Go-short, a minimalist URL shortener I built in Go in order to learn the language.

Routes:
POST /shorten - Create a shortened URL
		Form data: url (the original URL to shorten)
		Example response: {"shortened": "abc12345", "original": "https://example.com"}

GET /:short - Redirect to original URL
		Example: GET /abc12345 redirects to the original URL

GET /stats/:short - Get URL statistics
		Example response: {"original": "https://example.com", "accessCount": 5, "lastAccessed": "2025-07-06T10:30:00Z"}

Note: Shortened urls that haven't been clicked in the last 60 days are automatically deleted.`

	c.String(http.StatusOK, body)
}

func postCreateShortenUrl(c *gin.Context, coll *mongo.Collection) {
	originalUrl := c.PostForm("url")
	c.Bind(&originalUrl)

	// Avoid recreating url if it already exists in db
	var existing shortenedUrl
	err := coll.FindOne(context.TODO(), bson.M{"original": originalUrl}).Decode(&existing)
	if err == nil {
		c.IndentedJSON(http.StatusFound, gin.H{
			"shortened": existing.Shortened,
			"original":  existing.Original,
		})
		return
	}

	// Created shortened url
	short := strings.ToLower(generateRandomString(8))

	// Check for collisions in the database
	for {
		res := coll.FindOne(context.TODO(), bson.M{"shortened": short})
		if res.Err() == mongo.ErrNoDocuments {
			break
		}
		log.Println("Collision found, generating new url")
		short = strings.ToLower(generateRandomString(8))
	}

	newUrl := shortenedUrl{
		ID:           uuid.New().String(),
		Original:     originalUrl,
		Shortened:    short,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		AccessCount:  0,
	}

	// Insert url in db
	_, err = coll.InsertOne(context.TODO(), newUrl)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Mongo insert failed"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"shortened": newUrl.Shortened,
		"original":  newUrl.Original,
	})
}

func getUrlStats(c *gin.Context, coll *mongo.Collection) {
	var short = c.Param("short")

	// Check if the requested url exists in db
	var result shortenedUrl
	err := coll.FindOne(context.TODO(), bson.M{"shortened": short}).Decode(&result)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Shortened URL not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"original":     result.Original,
		"accessCount":  result.AccessCount,
		"lastAccessed": result.LastAccessed,
	})
}

func getShortenedUrl(c *gin.Context, coll *mongo.Collection) {
	var short = c.Param("short")

	// Check if the requested url exists in db
	var result shortenedUrl
	err := coll.FindOne(context.TODO(), bson.M{"shortened": short}).Decode(&result)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Shortened URL not found"})
		return
	}

	// Update object
	result.AccessCount++
	result.LastAccessed = time.Now()

	coll.UpdateOne(context.TODO(), bson.M{"shortened": short}, bson.M{"$set": bson.M{
		"accessCount":  result.AccessCount,
		"lastAccessed": result.LastAccessed,
	}})

	c.Redirect(http.StatusFound, result.Original)
}

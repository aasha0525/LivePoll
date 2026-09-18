package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"livepoll/db"
	"livepoll/models"
)

type createPollInput struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

// CreatePoll validates and stores a new poll. Requires auth.
func CreatePoll(c *gin.Context) {
	var input createPollInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Backend validation (this is the part the assignment cares about)
	if len(input.Question) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question is required"})
		return
	}

	// Drop empty option strings before validating count
	cleaned := make([]string, 0, len(input.Options))
	for _, o := range input.Options {
		if len(o) > 0 {
			cleaned = append(cleaned, o)
		}
	}
	if len(cleaned) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 2 options required"})
		return
	}

	userIDHex := c.GetString("userID")
	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	options := make([]models.Option, len(cleaned))
	for i, text := range cleaned {
		options[i] = models.Option{Text: text, Votes: 0}
	}

	poll := models.Poll{
		Question:  input.Question,
		Options:   options,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.Polls().InsertOne(ctx, poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create poll"})
		return
	}

	poll.ID = res.InsertedID.(bson.ObjectID)
	c.JSON(http.StatusCreated, poll)
}

// GetPoll fetches a single poll by id. Public — no auth required, so the
// person receiving the shared link can view and vote on it.
func GetPoll(c *gin.Context) {
	id := c.Param("id")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var poll models.Poll
	if err := db.Polls().FindOne(ctx, bson.M{"_id": objID}).Decode(&poll); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

// MyPolls lists polls created by the logged-in user, for the dashboard.
func MyPolls(c *gin.Context) {
	userIDHex := c.GetString("userID")
	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.Polls().Find(ctx, bson.M{"created_by": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch polls"})
		return
	}
	defer cursor.Close(ctx)

	polls := []models.Poll{}
	if err := cursor.All(ctx, &polls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not decode polls"})
		return
	}

	c.JSON(http.StatusOK, polls)
}

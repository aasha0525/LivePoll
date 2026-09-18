package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"livepoll/db"
)

type voteInput struct {
	OptionIndex int `json:"option_index"`
}

type voteUpdate struct {
	PollID      string `json:"poll_id"`
	OptionIndex int    `json:"option_index"`
	Votes       int64  `json:"votes"`
}

// Vote is the most important part of the assignment: it bumps the count in
// Redis for an instant response, publishes the change over Redis Pub/Sub so
// every open results page hears about it immediately via WebSocket, and
// persists the new count to MongoDB in the background so it survives restarts.
func Vote(c *gin.Context) {
	pollID := c.Param("id")
	objID, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poll id"})
		return
	}

	var input voteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Confirm the poll and option exist before touching Redis.
	var poll struct {
		Options []struct {
			Text  string `bson:"text"`
			Votes int64  `bson:"votes"`
		} `bson:"options"`
	}
	if err := db.Polls().FindOne(ctx, bson.M{"_id": objID}).Decode(&poll); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
		return
	}
	if input.OptionIndex < 0 || input.OptionIndex >= len(poll.Options) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid option"})
		return
	}

	// 1. Fast increment in Redis (this is the "live" counter).
	redisKey := "poll:" + pollID + ":option:" + itoa(input.OptionIndex)
	newCount, err := db.Rdb.Incr(db.Ctx, redisKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "vote failed"})
		return
	}

	// 2. Publish the update so every connected results page updates instantly.
	update := voteUpdate{PollID: pollID, OptionIndex: input.OptionIndex, Votes: newCount}
	payload, _ := json.Marshal(update)
	if err := db.Rdb.Publish(db.Ctx, db.PollChannel(pollID), payload).Err(); err != nil {
		log.Printf("redis publish error: %v", err)
	}

	// 3. Persist to MongoDB in the background so the count survives restarts.
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		field := "options." + itoa(input.OptionIndex) + ".votes"
		_, err := db.Polls().UpdateOne(bgCtx, bson.M{"_id": objID}, bson.M{"$inc": bson.M{field: 1}})
		if err != nil {
			log.Printf("mongo vote persist error: %v", err)
		}
	}()

	c.JSON(http.StatusOK, update)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

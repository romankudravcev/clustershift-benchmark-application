package api

import (
	"benchmarker/db"
	"benchmarker/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func StartBenchmark(c *gin.Context) {
	err := db.DeleteAllMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Benchmark started - all messages deleted"})
}

func PostMessage(c *gin.Context) {
	var requestBody struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := models.Message{
		Content:   requestBody.Content,
		CreatedAt: time.Now().UTC(),
	}

	id, err := db.SaveMessage(&message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetMessages(c *gin.Context) {
	messages, err := db.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// Additional helper endpoint to get a single message
func GetMessage(c *gin.Context) {
	id := c.Param("id")
	messageID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := db.GetMessage(messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if message == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

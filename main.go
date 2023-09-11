package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type FollowResult struct {
	Handle string `json:"handle"`
	Output string `json:"output"`
}

var ctx = context.Background()

func main() {
	r := gin.Default()

	// GET endpoint that prints "Hello, World!"
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// POST endpoint for login and follow handles
	r.POST("/follow", func(c *gin.Context) {
		redisUrl := os.Getenv("REDIS_URL")
		opt, _ := redis.ParseURL(redisUrl)
		client := redis.NewClient(opt)

		handle := c.PostForm("handle")
		password := c.PostForm("password")
		handlesToFollow := c.PostForm("follow")

		handle = strings.TrimSpace(handle)
		if !strings.HasSuffix(handle, ".bsky.social") {
			handle += ".bsky.social"
		}

		handleList := strings.Split(handlesToFollow, ",")

		// Define the path to your CLI binary
		cliPath := "./bsky"
		loginCmd := exec.Command(cliPath, "login", handle, password)
		loginErr := loginCmd.Run()

		if loginErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"login error": loginErr.Error()})
			return
		}

		// Create an array to store follow results
		followResults := []FollowResult{}

		// Loop through the handles to follow and execute the CLI follow command for each
		for _, h := range handleList {
			h = strings.TrimSpace(h)
			if !strings.HasSuffix(h, ".bsky.social") {
				h += ".bsky.social"
			}
			followCmd := exec.Command(cliPath, "follow", h)
			followOutput, followErr := followCmd.CombinedOutput()
			if followErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"follow error": followErr.Error()})
				return
			}
			fmt.Println("[FOLLOWED] " + h)
			client.Incr(ctx, "count")
			followResults = append(followResults, FollowResult{Handle: h, Output: string(followOutput)})
		}

		// Return the follow results as JSON
		c.JSON(http.StatusOK, followResults)
	})

	r.Run(":8080")
}

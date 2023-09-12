package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	// NOTE: This is not needed in deployed environment
	// "github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type FollowResult struct {
	Handle string `json:"handle"`
	Output string `json:"output"`
}

type UserInput struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
	Follow   string `json:"follow"`
}

var ctx = context.Background()

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	r.POST("/follow", func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
			return
		}

		var user UserInput

		// Bind the JSON request body to the UserInput struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create separate variables for handle, password, and follow
		handle := user.Handle
		password := user.Password
		follow := user.Follow

		handle = strings.TrimSpace(handle)
		if !strings.HasSuffix(handle, ".bsky.social") {
			handle += ".bsky.social"
		}

		handleList := strings.Split(follow, ",")

		cliPath := "./bsky"

		loginCmd := exec.Command(cliPath, "login", handle, password)
		loginErr := loginCmd.Run()

		if loginErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"login error": loginErr.Error()})
			return
		}

		// NOTE: This is not needed in deployed environment
		// envErr := godotenv.Load()
		// if envErr != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"env error": envErr.Error()})
		// 	return
		// }

		redisUrl := os.Getenv("REDIS_URL")
		opt, _ := redis.ParseURL(redisUrl)
		client := redis.NewClient(opt)

		followResults := []FollowResult{}

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

		c.JSON(http.StatusOK, followResults)
	})

	r.Run(":8080")
}

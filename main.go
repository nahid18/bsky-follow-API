package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type FollowResult struct {
	Handle string `json:"handle"`
	Output string `json:"output"`
}

func main() {
	r := gin.Default()

	// GET endpoint that prints "Hello, World!"
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// POST endpoint for login and follow handles
	r.POST("/follow", func(c *gin.Context) {
		// Parse the handle, password, and handles to follow from the request
		handle := c.PostForm("handle")
		password := c.PostForm("password")
		handlesToFollow := c.PostForm("follow")

		handle = strings.TrimSpace(handle)
		if !strings.HasSuffix(handle, ".bsky.social") {
			handle += ".bsky.social"
		}

		// Split the comma-separated handles into an array
		handleList := strings.Split(handlesToFollow, ",")

		// Define the path to your CLI binary
		cliPath := "./bsky"

		// Execute the login CLI command
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
			// Trim leading and trailing whitespace
			h = strings.TrimSpace(h)

			// Check if the handle has ".bsky.social" at the end
			if !strings.HasSuffix(h, ".bsky.social") {
				h += ".bsky.social"
			}

			fmt.Println("Following " + h)

			followCmd := exec.Command(cliPath, "follow", h)
			followOutput, followErr := followCmd.CombinedOutput()
			if followErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"follow error": followErr.Error()})
				return
			}

			followResults = append(followResults, FollowResult{Handle: h, Output: string(followOutput)})
		}

		// Return the follow results as JSON
		c.JSON(http.StatusOK, followResults)
	})

	r.Run(":8080")
}

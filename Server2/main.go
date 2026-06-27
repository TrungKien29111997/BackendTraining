package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/user/:userID", func(c *gin.Context) {
		userId := c.Param("userID")
		atk := c.Query("atk")
		c.JSON(200, gin.H{
			"userID": userId,
			"atk":    atk,
		})
	})
	r.Run(":8080")
}

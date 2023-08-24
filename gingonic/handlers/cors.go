package handlers

import (
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://jopit.com.ar")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, authorization, headerauthorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Vary", "Access-Control-Request-Method")
		c.Writer.Header().Set("Vary", "Access-Control-Request-Headers")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	corsConfig = cors.Config{
		AllowOrigins:     []string{"https://jopit.com.ar"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
)

func CORSMiddleware() gin.HandlerFunc {
	// Crea el middleware de CORS utilizando la configuración
	corsMiddleware := cors.New(corsConfig)

	return func(c *gin.Context) {
		// Aplica el middleware de CORS a la solicitud actual
		corsMiddleware(c)
		c.Next()
	}
}

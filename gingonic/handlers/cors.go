package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	corsConfig = cors.Config{
		AllowAllOrigins:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "token"},
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

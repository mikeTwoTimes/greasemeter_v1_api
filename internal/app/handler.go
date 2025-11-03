package app

import (
	"Greasemeter-rest-api/internal/middlewares"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *App) handler() http.Handler {
	g := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders:     []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	
	g.Use(cors.New(config))

	v1 := g.Group("/v1")
	auth := v1.Group("/")
	usersStore := users.NewStore(a.db)

	auth.Use(middlewares.Auth(a.jwtSecret, usersStore.UserExists))

	usersHandler := users.NewHandler(usersStore, a.jwtSecret, a.mailer)
	usersHandler.RegisterRoutes(v1, auth)

	return g
}

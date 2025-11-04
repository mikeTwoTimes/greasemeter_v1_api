package app

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/middlewares"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/services/bookmarks"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/services/places"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/services/recommendations"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/services/reports"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/services/reviews"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/services/users"
)

func (a *App) handler() http.Handler {
	g := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders:     []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
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

	placesStore := places.NewStore(a.db)
	placesHandler := places.NewHandler(placesStore)
	placesHandler.RegisterRoutes(v1)

	reviewsStore := reviews.NewStore(a.db)
	reviewsHandler := reviews.NewHandler(reviewsStore)
	reviewsHandler.RegisterRoutes(v1, auth)

	bookmarksStore := bookmarks.NewStore(a.db)
	bookmarksHandler := bookmarks.NewHandler(bookmarksStore)
	bookmarksHandler.RegisterRoutes(auth)

	recommendationsStore := recommendations.NewStore(a.db)
	recommendationsHandler := recommendations.NewHandler(
		recommendationsStore,
	)
	recommendationsHandler.RegisterRoutes(auth)

	reportsStore := reports.NewStore(a.db)
	reportsHandler := reports.NewHandler(reportsStore)
	reportsHandler.RegisterRoutes(auth)
	
	return g
}

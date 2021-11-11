package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)

func adminRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Get("/info", handlers.GetUserInfo)
		user.Delete("/logout", handlers.Logout)
		user.Post("/add-restaurant",handlers.AddRestaurant)
	})
}
package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)

func adminRoutes(r chi.Router) {
	r.Group(func(admin chi.Router) {
		admin.Get("/info", handlers.GetUserInfo)
		admin.Delete("/logout", handlers.Logout)
		admin.Post("/add-restaurant",handlers.AddRestaurant)
		admin.Post("/add-dish",handlers.AddDish)
		admin.Post("/add-user",handlers.AddUser)
		admin.Get("/get-users",handlers.GetUsers)
		admin.Get("/get-all",handlers.GetAll)
		admin.Post("/get-dish",handlers.GetDish)
		admin.Post("/get-distance",handlers.GetDistance)
	})
}
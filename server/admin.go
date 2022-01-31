package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)

func adminRoutes(r chi.Router) {
	r.Group(func(admin chi.Router) {
		admin.Get("/info", handlers.GetUserInfo)
		admin.Delete("/logout", handlers.Logout)
		admin.Post("/restaurant",handlers.AddRestaurant)
		admin.Post("/dish",handlers.AddDish)
		admin.Post("/user",handlers.AddUser)
		admin.Get("/users",handlers.GetUsers)
		admin.Get("/all",handlers.GetAll)
		admin.Post("/dishes",handlers.GetDish)
		admin.Post("/distance",handlers.GetDistance)
	})
}
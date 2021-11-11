package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)

func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Get("/info", handlers.GetUserInfo)
		user.Delete("/logout", handlers.Logout)
		user.Post("/add-address",handlers.AddAddress)
		user.Get("/get-all",handlers.GetAll)
		user.Get("/get-dish",handlers.GetDish)
		 user.Get("/get-all",handlers.GetAll)
	})
}

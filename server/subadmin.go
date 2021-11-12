package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)



func subadminroutes(r chi.Router) {
	r.Group(func(subadmin chi.Router) {
		subadmin.Get("/info", handlers.GetUserInfo)
		subadmin.Delete("/logout", handlers.Logout)
		subadmin.Post("/add-restaurant",handlers.AddSubAdminRestaurant)
		subadmin.Post("/add-dish",handlers.AddSubAdminDish)
		subadmin.Post("/add-user",handlers.AddSubAdminUser)
		subadmin.Get("/get-users",handlers.GetSubAdminUsers)

	})
}

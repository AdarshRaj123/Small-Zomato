package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)



func subadminroutes(r chi.Router) {
	r.Group(func(subadmin chi.Router) {
		subadmin.Get("/info", handlers.GetUserInfo)
		subadmin.Delete("/logout", handlers.Logout)
		subadmin.Post("/restaurant",handlers.AddSubAdminRestaurant)
		subadmin.Post("/dish",handlers.AddSubAdminDish)
		subadmin.Post("/user",handlers.AddSubAdminUser)
		subadmin.Get("/users",handlers.GetSubAdminUsers)
	})
}

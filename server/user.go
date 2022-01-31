package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)

func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Get("/info", handlers.GetUserInfo)
		user.Delete("/logout", handlers.Logout)
		user.Post("/address",handlers.AddAddress)
		user.Get("/all",handlers.GetAll)
		user.Get("/dish",handlers.GetDish)
		user.Post("/distance",handlers.GetDistance)
		user.Post("/upload-file",handlers.UploadFile)

	})
}

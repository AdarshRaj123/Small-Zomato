package server

import (
	"SmallZomato/handlers"
	"github.com/go-chi/chi/v5"
)



func commonroutes(r chi.Router) {
	r.Group(func(common chi.Router) {
		common.Get("/info", handlers.GetDetails)





	})
}

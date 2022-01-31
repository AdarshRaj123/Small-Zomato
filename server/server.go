package server

import (
	"SmallZomato/handlers"
	"SmallZomato/middlewares"
	"SmallZomato/utils"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

// SetupRoutes provides all the routes that can be used
func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/small-zomato", func(v1 chi.Router) {
		v1.Use(middlewares.CommonMiddlewares()...)
		v1.Post("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})
		v1.Route("/", func(public chi.Router) {
			public.Post("/register", handlers.RegisterUser)
			public.Post("/login", handlers.LoginUser)

		})
		v1.Route("/user", func(user chi.Router) {
			user.Use(middlewares.AuthMiddleware)
			user.Use(middlewares.UserCheck())
			user.Group(userRoutes)
		})

		v1.Route("/internal", func(internal chi.Router) {
			internal.Use(middlewares.AuthMiddleware)
			internal.Route("/admin", func(admin chi.Router) {
				admin.Use(middlewares.AdminCheck())
				admin.Group(adminRoutes)
			})
			internal.Route("/sub-admin", func(subAdmin chi.Router) {
				subAdmin.Use(middlewares.SubAdminCheck())
				subAdmin.Group(subadminroutes)
			})
			internal.Route("/users", func(users chi.Router) {
				users.Use(middlewares.CommonCheck())
				users.Group(commonroutes)
			})
		})
	})
	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}

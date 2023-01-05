package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

// Handler - stores pointers to the router, order service and server
type Handler struct {
	Router  *chi.Mux
	Service OrderService
	Server  *http.Server
}

// Response object
type Response struct {
	Message string `json:"message"`
}

// NewHandler - returnes a pointer to a Handler
func NewHandler(service OrderService) *Handler {
	log.Info("Setting handler")
	h := &Handler{
		Service: service,
	}
	h.Router = chi.NewRouter()
	//setting json formatting middleware
	h.Router.Use(JSONMiddleware)
	//setting logger for all incoming request
	h.Router.Use(middleware.Logger)
	//timeout all requests that take longer than 15 seconds
	h.Router.Use(middleware.Timeout(15 * time.Second))
	//set up the routes
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.Router,
	}
	return h
}

// mapRoutes - sets up all routes in app
func (h *Handler) mapRoutes() {
	h.Router.Get("/api/v1/orders/{id}", h.GetOrder)
	h.Router.Get("/api/v1/orders", h.GetAllOrders)
	h.Router.Post("/api/v1/orders", h.PostOrder)
	h.Router.Put("/api/v1/orders/{id}", h.UpdateOrder)
	h.Router.Delete("/api/v1/orders/{id}", h.DeleteOrder)

	h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}

// Serve - gracefully serves handler
func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	//create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("server shutting down gracefully")
	return nil
}

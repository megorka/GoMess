package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Host string `yaml:"HTTP_HOST" env:"HTTP_HOST" env-default:"localhost"`
	Port string `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8080"`
}

type Router struct {
	config  Config
	Router  *mux.Router
	Handler Handler
}

func NewRouter(cfg Config, h *Handler) *Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auth/signup", h.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST")
	return &Router{
		config: cfg,
		Router: r,
	}
}

func (r *Router) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", r.config.Host, r.config.Port), r.Router))
}

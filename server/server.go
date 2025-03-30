package server

import (
	"YuriyMishin/metrics/handlers"
	"YuriyMishin/metrics/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	storage := storage.NewMemStorage()
	metricHandlers := handlers.NewMetricHandlers(storage)

	r := mux.NewRouter()
	r.HandleFunc("/", metricHandlers.RootHandler).Methods("GET")
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", metricHandlers.UpdateHandler).Methods("POST")
	r.HandleFunc("/value/{metricType}/{metricName}", metricHandlers.ValueHandler).Methods("GET")

	return &Server{router: r}
}

func (s *Server) Start(addr string) error {
	fmt.Printf("Server is running on http://%s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

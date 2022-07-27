package main

import (
	"net/http"
	"sync"

	"github.com/Infoblox-CTO/heka-ui/pkg/backend"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Server struct represents the backend API server
type Server struct {
	wg      *sync.WaitGroup
	done    chan bool
	address string
}

// Name returns the name of this package
func (s *Server) Name() string {
	return "Heka ui backend server"
}

//
func (s *Server) Init(w *sync.WaitGroup) error {
	s.wg = w
	s.address = "all the backend server API address"
	return nil
}

//
func (s *Server) Start() error {
	//create a new router
	router := mux.NewRouter()
	go func() {
		//specify endpoints, handler functions and HTTP method
		router.HandleFunc("/health-check", backend.HealthCheck).Methods("GET")
		router.HandleFunc("/cloud_manifest", backend.CloudManifest).Methods("GET")
		router.HandleFunc("/grid_data", backend.GridData).Methods("GET")
		router.HandleFunc("/node_data/{uniq_object}/{uniq_id}", backend.NodeData).Methods("GET")
		router.HandleFunc("/list_nodes", backend.ListNodes).Methods("GET")
		router.HandleFunc("/login", backend.Login).Methods("POST")
		router.HandleFunc("/logout", backend.Logout).Methods("POST")
		http.Handle("/", router)

		//start and listen to requests
		serverAdd := viper.GetString("server.address")
		logger.Infoln("Started backend server on", serverAdd)
		http.ListenAndServe(serverAdd, router)
	}()
	return nil
}

//
func (s *Server) Check() bool {
	// Check is heka enabled
	return true
}

//
func (s *Server) Stop() {
	s.done <- true
}

//
func (s *Server) IsUP() bool {
	return true
}

func ServerFactory() Handle {
	return &Server{}
}

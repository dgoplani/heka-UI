package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"context"
	"sync"
	"time"

	"github.com/Infoblox-CTO/heka-ui/pkg/backend"
	"github.com/Infoblox-CTO/heka-ui/utils"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Server struct represents the backend API server
type Server struct {
	wg      *sync.WaitGroup
	srv     *http.Server
	active  bool
	done    chan bool
	address string
}

// Name returns the name of this package
func (s *Server) Name() string {
	return "Heka ui backend server"
}

// SysStat used to get Grid UUid and version details
type SysStat struct {
	GridUUID    string `json:"grid_uuid"`
	NiosVersion string `json:"nios_version"`
}

func (s *Server) Init(w *sync.WaitGroup) error {
	s.wg = w
	s.active = true
	s.done = make(chan bool)
	s.address = "all the backend server API address"
	go func() {
		fetchGridData()
	}()
	go s.triggerControldScript()
	return nil
}

func (s *Server) Start() error {
	//create a new router
	router := mux.NewRouter()
	serverAdd := viper.GetString("server.address")
	srv := &http.Server{
		Addr:    serverAdd,
		Handler: router,
	}
	s.srv = srv

	go func() {
		//specify endpoints, handler functions and HTTP method
		router.HandleFunc("/health-check", backend.HealthCheck).Methods("GET")
		router.HandleFunc("/cloud_manifest", backend.CloudManifest).Methods("GET")
		router.HandleFunc("/grid_data", backend.GridData).Methods("GET")
		router.HandleFunc("/node_data/{uniq_object}/{uniq_id}", backend.NodeData).Methods("GET")
		router.HandleFunc("/list_nodes", backend.ListNodes).Methods("GET")
		router.HandleFunc("/login", backend.Login).Methods("POST")
		router.HandleFunc("/logout", backend.Logout).Methods("POST")
		router.HandleFunc("/is_ready", backend.IsReady).Methods("GET")

		//start and listen to requests

		logger.Infoln("Started backend server on", serverAdd)

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}

func (s *Server) Check() bool {
	// Check is heka enabled
	if isHekaUiValidToRun() && isHekaEnabled() {
		// Check whether heka is enabled
		return true
	}

	return false

}

func (s *Server) Stop() {
	if s.active {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		if err := s.srv.Shutdown(ctx); err != nil {
			logger.Infoln("Failed to Shutt down http heka-ui server")
		}
		logger.Infoln("Shutting down http heka-ui server")
		s.done <- true
	}
	s.active = false
}

func (s *Server) IsUP() bool {
	return s.active
}

func ServerFactory() Handle {
	return &Server{}
}

// fetchGridData sets the value for grid_uuid and nios_version
func fetchGridData() {
	var data SysStat
	var err error
	maxRetryPeriodres := 10 * time.Minute
	retryIntervalres := 1 * time.Minute
	for {
		data, err = getNiosDataFromScript()
		if err != nil {
			if retryIntervalres > maxRetryPeriodres {
				time.Sleep(5 * time.Second)
				logger.Warningln("Not able to collect grid UUId and Nios version")
				break
			} else {
				time.Sleep(retryIntervalres)
				retryIntervalres *= 2
				logger.Errorf("Grid data collection script failed retrying for %v Error: %v", retryIntervalres, err)
				continue
			}
		} else {
			viper.Set("grid_uuid", data.GridUUID)
			viper.Set("nios_version", data.NiosVersion)
			break
		}

	}
}

// getNiosDataFromScript collects grid_uuid and nios_version from Nios Db
func getNiosDataFromScript() (SysStat, error) {
	var stat SysStat
	status, resp, err := utils.MakePostRequest(200, "collect_grid_data", []string{})
	if err != nil {
		logger.Errorf("System Check HTTP POST request failed : %v with status code : %v", err, status)
		return stat, err
	}

	if status != http.StatusOK {
		err = fmt.Errorf("System Check HTTP Post request failed. Received HTTP Response code : %v", status)
		return stat, err
	}
	err = json.Unmarshal(resp, &stat)
	if err != nil {
		logger.Errorf("Error while Unmarshalling System Status HTTP Response: %v", err)
		return stat, err
	}
	return stat, nil
}

// isHekaUiValidToRunchecks whether valid to run check request on a particular grid
func isHekaUiValidToRun() bool {
	x := viper.GetStringMap("labels")

	if (x["nios_role"] == "GM" || x["nios_role"] == "standalone") && (x["ha_status"] == "Active" || x["ha_status"] == "Not Configured") {
		logger.Debugln("Heka-ui valid to run")
		return true
	}
	return false
}

// isHekaEnabled checks whether heka_optin flag is enabled
func isHekaEnabled() bool {

	if viper.GetBool("labels.heka_optin") {
		return true
	}
	return false
}

// triggerControldScript collects hotfix data for every hotfix interval
func (s *Server) triggerControldScript() {
	for {
		select {
		case <-s.done:
			logger.Infoln("Closing hotfix manager to collect data")
			return
		default:
			utils.CollectHotfixData("--trigger")
			wait := viper.GetDuration("hotfix.interval")
			logger.Infoln("Collecing hotfix data from all nodes for every ", wait)
			time.Sleep(wait)
		}
	}
}

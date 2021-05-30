package znet

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type healthStatus int

const (
	Healthy healthStatus = iota
	Unhealthy
)

func (s healthStatus) String() string {
	return [...]string{
		"healthy",
		"unhealthy",
	}[s]
}

type statusCheckHandler struct {
	server *Server
}

func (s *statusCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := make(map[string]interface{})

	status["status"] = Healthy.String()
	errors := []string{}

	if s.server != nil {
		if s.server.grpcServer == nil {
			status["status"] = Unhealthy.String()
			errors = append(errors, ErrNoGRPCServices.Error())
		} else {
			grpcInfo := s.server.grpcServer.GetServiceInfo()
			status["grpcServices"] = len(grpcInfo)
		}
	}

	if len(errors) > 0 {
		status["errors"] = errors
	}

	// payload, err := json.MarshalIndent(status, "", "  ")
	payload, err := json.Marshal(status)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Fprint(w, string(payload))
}

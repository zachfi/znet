package znet

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type statusCheckHandler struct {
	server *Server
}

func (s *statusCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := make(map[string]interface{})

	status["status"] = "healthy"

	grpcInfo := s.server.grpcServer.GetServiceInfo()

	status["grpcServices"] = len(grpcInfo)

	payload, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Error(err)
	}

	fmt.Fprint(w, string(payload))
}

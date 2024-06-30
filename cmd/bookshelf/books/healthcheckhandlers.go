package books

import (
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/rest"
)

type HealthCheckMessage struct {
	Status string `json:"status"`
}

func (m *Module) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	healthCheckMessage := HealthCheckMessage{
		Status: "available",
	}

	m.logger.Info("writing response", "response", healthCheckMessage)
	rest.Respond(w, r, http.StatusOK, healthCheckMessage, nil)
}

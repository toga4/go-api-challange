package handler

import (
	"net/http"

	"github.com/toga4/go-api-challange/log"
	"github.com/unrolled/render"
)

type ChallangeHandler interface {
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
}

type challangeHandler struct {
	Render *render.Render
}

func NewChallangeHandler() ChallangeHandler {
	return &challangeHandler{
		Render: render.New(),
	}
}

func (ch *challangeHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	log.R(r).Info("health check")
	ch.Render.Text(w, http.StatusOK, "alive!")
}

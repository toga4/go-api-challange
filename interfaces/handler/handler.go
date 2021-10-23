package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/toga4/go-api-challange/log"
	"github.com/toga4/go-api-challange/usecase"
	"github.com/unrolled/render"
)

type ChallangeHandler interface {
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
	HandleHello(w http.ResponseWriter, r *http.Request)
	HandleDelegate(w http.ResponseWriter, r *http.Request)
}

type challangeHandler struct {
	Render           *render.Render
	ChallangeUsecase usecase.ChallangeUsecase
}

func NewChallangeHandler(cu usecase.ChallangeUsecase) ChallangeHandler {
	return &challangeHandler{
		Render:           render.New(),
		ChallangeUsecase: cu,
	}
}

func (ch *challangeHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	log.R(r).Info("health check")
	ch.Render.Text(w, http.StatusOK, "alive!")
}

func (ch *challangeHandler) HandleHello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		b = []byte("(read error)")
	}
	log.R(r).WithValues(
		"headers", r.Header,
		"body", string(b),
	).Info("hello")
	ch.Render.Text(w, http.StatusOK, "hello!")
}

func (ch *challangeHandler) HandleDelegate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.R(r).Info("delegating")
	if err := ch.ChallangeUsecase.DelegateRequest(ctx); err != nil {
		log.R(r).Error(err, "HandleDelegate: DelegateRequest: err")
		ch.Render.Text(w, http.StatusInternalServerError, err.Error())
	}

	ch.Render.Text(w, http.StatusOK, "delegated!")
}

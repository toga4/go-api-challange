package usecase

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/toga4/go-api-challange/log"
	"github.com/toga4/go-api-challange/middleware"
)

type ChallangeUsecase interface {
	DelegateRequest(ctx context.Context) error
}

type challangeUsecase struct {
	HttpClient *http.Client
	HostURI    string
}

func NewChallangeUsecase(hostURI string) ChallangeUsecase {
	httpClient := &http.Client{
		Transport: &middleware.GCPTraceTransport{},
	}
	return &challangeUsecase{
		HttpClient: httpClient,
		HostURI:    hostURI,
	}
}

func (u *challangeUsecase) DelegateRequest(ctx context.Context) error {
	method := http.MethodGet
	url := u.HostURI + "/"
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	res, err := u.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode < 200 || 300 <= res.StatusCode {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			b = []byte("read error: " + err.Error())
		}
		log.C(ctx).Error(nil, "error response", "response_body", string(b))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		b = []byte("read error: " + err.Error())
	}
	log.C(ctx).Info("request successful", "response_body", string(b))

	return nil
}

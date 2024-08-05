package nodeshift

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/devilcove/httpclient"
	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/models"
)

type request struct {
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}

type response struct{}

const (
	backendHostProduction  = "app.nodeshift.com"
	backendHostStaging     = "app.nodeshift.co"
	backendHostDevelopment = "app.nodeshift.local"
)

func Notify(host, uuid, token string) error {
	backendHost := backendHostDevelopment
	if strings.HasSuffix(host, "nodeshift.network") {
		backendHost = backendHostProduction
	} else if strings.HasSuffix(host, "nodeshift.co") {
		backendHost = backendHostStaging
	}

	api := httpclient.JSONEndpoint[response, models.ErrorResponse]{
		URL:    "https://" + backendHost,
		Route:  "/api/vpc/register",
		Method: http.MethodPost,
		Data: request{
			Uuid:  uuid,
			Token: token,
		},
		Response:      response{},
		ErrorResponse: models.ErrorResponse{},
	}

	_, errData, err := api.GetJSON(response{}, models.ErrorResponse{})
	if err != nil {
		if errors.Is(err, httpclient.ErrStatus) {
			logger.FatalLog("error registering with server", strconv.Itoa(errData.Code), errData.Message)
		}

		return err
	}

	return nil
}

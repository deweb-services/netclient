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
	Server string `json:"server"`
	Uuid   string `json:"uuid"`
}

type response struct{}

const (
	backendHostProduction  = "app.nodeshift.com"
	backendHostStaging     = "app.nodeshift.co"
	backendHostDevelopment = "app.nodeshift.local"
)

var unknownServerTypeErr = errors.New("unknown server type")

func Notify(event models.HostUpdate) error {
	if event.Action != models.JoinHostToNetwork {
		return nil
	}

	backendHost, server, err := getServerHost(event.Node.Server)
	if err != nil {
		return unknownServerTypeErr
	}

	api := httpclient.JSONEndpoint[response, models.ErrorResponse]{
		URL:    "https://" + backendHost,
		Route:  "/api/vpc/register",
		Method: http.MethodPost,
		Data: request{
			Uuid:   event.Node.ID.String(),
			Server: server,
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

func getServerHost(server string) (string, string, error) {
	if strings.HasSuffix(server, "nodeshift.network") {
		return backendHostProduction, strings.TrimSuffix(server, ".nodeshift.network"), nil
	} else if strings.HasSuffix(server, "nodeshift.co") {
		return backendHostStaging, strings.TrimSuffix(server, ".nodeshift.co"), nil
	} else if strings.HasSuffix(server, "nodeshift.cloud") {
		return backendHostDevelopment, strings.TrimSuffix(server, ".nodeshift.cloud"), nil
	}

	return "", "", unknownServerTypeErr
}

package requester

import (
	"fmt"
	"github.com/hanchon-live/autostake-bot/internal/util"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: 2 * time.Second,
}

var settings util.Config

func init() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic("Error reading the config")
	}
	settings = config
}

func MakeGetRequest(endpointType string, url string) (string, error) {
	var endpoints []string
	if endpointType == "rest" {
		endpoints = settings.Rest
	} else if endpointType == "jrpc" {
		endpoints = settings.Jrpc
	} else if endpointType == "web3" {
		endpoints = settings.Web3
	} else {
		return "", fmt.Errorf("Invalid endpoint type")
	}

	for _, endpoint := range endpoints {
		var sb strings.Builder
		sb.WriteString(endpoint)
		sb.WriteString(url)

		resp, err := client.Get(sb.String())

		if err != nil || resp.StatusCode == 429 {
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil || len(string(body)) == 0 {
			continue
		}

		return string(body), nil
	}

	return "", fmt.Errorf("All endpoints are down")
}

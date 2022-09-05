package requester

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hanchon-live/autostake-bot/internal/util"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

var settings util.Config

func init() {
	config, err := util.LoadConfig()
	if err != nil {
		fmt.Println("Error reading the config, using localnet values!")
	}
	settings = config
}

func MakeGetRequest(endpointType string, url string) (string, error) {
	var endpoints []string
	if endpointType == "rest" {
		endpoints = settings.Rest
	} else {
		return "", fmt.Errorf("Invalid endpoint type")
	}

	for _, endpoint := range endpoints {
		var sb strings.Builder
		sb.WriteString(endpoint)
		sb.WriteString(url)

		resp, err := client.Get(sb.String())

		if err != nil {
			continue
		}
		if resp.StatusCode == 404 {
			return "", fmt.Errorf("Element not found")
		}

		if resp.StatusCode != 200 {
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

func MakePostRequest(endpointType string, url string, param []byte) (string, error) {
	body := bytes.NewBuffer(param)

	var endpoints []string
	if endpointType == "rest" {
		endpoints = settings.Rest
	} else {
		return "", fmt.Errorf("Invalid endpoint type")
	}

	for _, endpoint := range endpoints {
		var sb strings.Builder
		sb.WriteString(endpoint)
		sb.WriteString(url)

		resp, err := client.Post(sb.String(), "application/json", body)

		if err != nil || resp.StatusCode == 429 {
			continue
		}

		defer resp.Body.Close()

		bodyResponse, err := ioutil.ReadAll(resp.Body)

		if err != nil || len(string(bodyResponse)) == 0 {
			continue
		}

		return string(bodyResponse), nil
	}

	return "", fmt.Errorf("All endpoints are down")
}

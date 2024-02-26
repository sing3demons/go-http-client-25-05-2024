package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpClient interface {
	Get(url string, headers map[string]string) (string, error)
	Post(url string, payload any, headers map[string]string) (string, error)
}

type httpClient struct {
	Client http.Client
}

func NewHttpClient() HttpClient {
	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "10"
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeout = 10
	}

	return &httpClient{
		Client: http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

func (client *httpClient) makeRequest(method, url string, payload io.Reader, headers map[string]string) (string, error) {
	req, err := client.buildRequest(method, url, payload, headers)
	if err != nil {
		return "", err
	}

	response, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func (client *httpClient) buildRequest(method, url string, payload io.Reader, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (client *httpClient) Get(url string, headers map[string]string) (string, error) {
	return client.makeRequest(http.MethodGet, url, nil, headers)
}

func (client *httpClient) Post(url string, payload any, headers map[string]string) (string, error) {

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return client.makeRequest(http.MethodPost, url, bytes.NewReader(payloadData), headers)
}

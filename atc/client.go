package atc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type VersionClient struct {
	baseURL string
}

func NewVersionClient(baseURL string) *VersionClient {
	return &VersionClient{baseURL: baseURL}
}

func (c *VersionClient) GetServerVersion() (string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	uri, err := url.Parse(c.baseURL)

	if err != nil {
		return "", err
	}

	uri.Path = "/api/v1/info"

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)

	if err != nil {
		return "", err
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	type ServerInfo struct {
		Version string `json:"version"`
	}

	info := ServerInfo{}
	err = json.Unmarshal(body, &info)

	if err != nil {
		return "", err
	}

	return info.Version, nil
}

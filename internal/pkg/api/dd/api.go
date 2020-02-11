package dd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const endpointv1 = "https://api.datadoghq.com/api/v1"

type API struct {
	apikey   string
	appkey   string
	timeout  time.Duration
	endpoint string
}

func NewAPI(endpoint string, apikey string, appkey string) API {
	if endpoint == "" {
		endpoint = endpointv1
	}
	return API{
		apikey:   apikey,
		appkey:   appkey,
		endpoint: endpoint,
		timeout:  15 * time.Second,
	}
}
func (a API) postMetrics(metrics []*Metric) error {

	post := map[string][]*Metric{
		"series": metrics,
	}
	endpoint := fmt.Sprintf("%s/series?api_key=%s", a.endpoint, a.apikey)

	return write(endpoint, post, a.timeout)
}

// writes a json blob
func write(endpoint string, data interface{}, timeout time.Duration) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	body := bytes.NewReader(raw)
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			err = urlErr.Err
		}
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusCreated:
		return nil
	}

	return fmt.Errorf("http status %v: %s", resp.StatusCode, string(responseBody))
}
func (c *Client) Watch(duration time.Duration) {
	go c.watch(duration)
}

func (c *Client) watch(duration time.Duration) {
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			if err := c.Flush(); err != nil {
				logrus.Warn("Failed to flush metrics from client: {}", err)

			}
		case msg := <-c.Stop:
			logrus.Error("Datadog Client stopped: {}", msg)
			ticker.Stop()
			return
		}
	}
}

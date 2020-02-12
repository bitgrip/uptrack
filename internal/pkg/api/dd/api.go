package dd

import (
	cons "bitbucket.org/bitgrip/uptrack/internal/pkg"
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
func (a API) postSeries(series []*Metric) error {

	post := map[string][]*Metric{
		"series": series,
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
func (c *Client) Watch(freq time.Duration) {
	go c.watch(freq)
}

func (c *Client) watch(freq time.Duration) {
	ticker := time.NewTicker(freq)

	for {
		select {
		case <-ticker.C:
			if err := c.Flush(); err != nil {
				logrus.Warn("Failed to flush metricsMap from client: {}", err)

			}
		case msg := <-c.Stop:
			logrus.Error("Datadog Client stopped: {}", msg)
			ticker.Stop()
			return
		}
	}
}

type DDTags map[string]string

func (t DDTags) ToTagList() []string {
	out := make([]string, 0)
	for k, v := range t {
		if k != cons.Host {
			out = append(out, k+":"+v)
		}
	}
	return unique(out)
}

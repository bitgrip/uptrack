package dd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	cons "github.com/bitgrip/uptrack/internal/pkg"
	"github.com/sirupsen/logrus"
)

type API struct {
	apikey   string
	appkey   string
	timeout  time.Duration
	endpoint string
}

func NewAPI(endpoint string, apikey string, appkey string) API {
	return API{
		apikey:   apikey,
		appkey:   appkey,
		endpoint: endpoint,
		timeout:  15 * time.Second,
	}
}
func (a API) postSeries(series []*Metric) error {

	timeout := a.timeout

	post := map[string][]*Metric{
		"series": series,
	}
	raw, err := json.Marshal(post)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	body := bytes.NewReader(raw)
	req, err := http.NewRequest(http.MethodPost, a.endpoint, body)

	q := req.URL.Query()
	q.Add("api_key", a.apikey)
	req.URL.RawQuery = q.Encode()

	//set Headers
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}
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
				logrus.Warn(fmt.Sprintf("Failed to flush metricsMap from client: %s", err))
			}
		case msg := <-c.Stop:
			logrus.Error(fmt.Sprintf("Datadog Client stopped: %s", msg))
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

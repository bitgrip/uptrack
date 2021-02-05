package ctl

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

//traces the request, callbacks are called for relevant metrics within.
func trace(registry metrics.Registry, name string) httptrace.ClientTrace {
	var connect time.Time
	start := time.Now()
	return httptrace.ClientTrace{

		ConnectStart: func(network, addr string) {
			connect = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			registry.SetConnectTime(name, float64(time.Since(connect).Milliseconds()))
		},

		GotFirstResponseByte: func() {
			registry.SetTTFB(name, float64(time.Since(start).Milliseconds()))

		},
	}
}

//the purpose of a handwritten Transport is to have a clean uncached request in each iteration.
func transport() http.Transport {
	return http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
func getAccessToken(upJob *job.UpJob) (string, error) {
	bearerToken := ""
	oauthServerResponse := &upJob.OauthServerResponse

	refreshToken := false
	if oauthServerResponse.ExpiresAt != (time.Time{}) {
		if oauthServerResponse.ExpiresAt.Sub(time.Now()).Seconds() < 10 {
			refreshToken = true
		}

	}
	if oauthServerResponse.AccessToken == "" {
		refreshToken = true
	}
	if oauthServerResponse.Refresh {
		refreshToken = true
		oauthServerResponse.Refresh = false
	}

	if refreshToken {
		values := url.Values{}
		//load params from Oauth
		for k, v := range upJob.Oauth.Params {
			values.Set(k, v)
		}
		if oauthServerResponse.RefreshToken != "" {
			values.Set("refresh_token", oauthServerResponse.RefreshToken)
			values.Set("request_type", "refresh_token")
		}

		authReq, _ := http.NewRequest("POST", upJob.Oauth.AuthUrl, strings.NewReader(values.Encode()))

		for k, v := range upJob.Oauth.Headers {
			authReq.Header.Add(k, v)
		}

		client := &http.Client{}
		authResp, err := client.Do(authReq)
		if err != nil {
			logrus.Fatal(err)
			return "", err
		}
		if authResp.StatusCode != http.StatusOK {
			errorMsg := fmt.Sprintf("Failed authentication on auth_url: '%s'", upJob.Oauth.AuthUrl)
			logrus.Error(errorMsg)
			return "", fmt.Errorf(errorMsg)

		}
		bytes, _ := ioutil.ReadAll(authResp.Body)
		oauthResponse := job.OauthResponse{}

		err = yaml.Unmarshal(bytes, &oauthResponse)

		*oauthServerResponse = oauthResponse
		oauthServerResponse.ExpiresAt = time.Now().Add(time.Second * time.Duration(oauthServerResponse.ExpiresIn))
	}
	bearerToken = oauthServerResponse.AccessToken

	return bearerToken, nil
}

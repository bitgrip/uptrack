package ctl

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"github.com/sirupsen/logrus"
)

func doUpChecks(registry metrics.Registry, upJob *job.UpJob) {
	jobName := upJob.Name
	url := upJob.URL

	//prepare request
	req, _ := http.NewRequest(string(upJob.Method), url, upJob.Body())
	clientTrace := trace(registry, jobName)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &clientTrace))

	//config := clientcredentials.Config{}
	//config.Client()
	//https://www.oauth.com/oauth2-servers/access-tokens/client-credentials/
	if upJob.Oauth.AuthUrl != "" {
		//Perform authentication via oauth client credentials flow
		bearerToken, err := getAccessToken(upJob)
		if err != nil {
			logrus.Error(fmt.Sprintf("Failed to receive Bearer Token for url: '%s' and auth_url: '%s'", upJob.URL, upJob.Oauth.AuthUrl))
			logrus.Error(err)
			return
		}
		//Add Bearer Token in Authorization Header
		if bearerToken != "" {
			req.Header.Add("Authorization", "Bearer "+bearerToken)
		}
	}

	for k, v := range upJob.Headers {
		req.Header.Add(k, v)
	}

	t := transport()
	//measure request time
	startReq := time.Now()
	resp, err := t.RoundTrip(req)
	if err != nil {
		logrus.Warn(fmt.Sprintf("Failed Request for job '%s'. msg: '%s' ", jobName, err))
		registry.IncCanNotConnect(jobName)
		return
	}
	endReq := time.Since(startReq)
	registry.SetRequestTime(jobName, float64(endReq.Milliseconds()))

	//time until expiry of ssl certs
	if upJob.CheckSSL {
		hours := time.Until(resp.TLS.PeerCertificates[0].NotAfter).Hours()
		registry.SetSSLDaysLeft(jobName, hours/24)
	}

	up := resp.StatusCode == upJob.Expected
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	} else {
		registry.SetBytesReceived(jobName, float64(len(bodyBytes)+len(resp.Header)))
	}
	if upJob.ContentMatch != "" {
		up = up && matchResponse(upJob.ContentMatch, bodyBytes, upJob.ReverseContentMatch)
	}
	if up {

		//count successful connections
		registry.IncCanConnect(jobName)

	} else {
		//count failed connections
		registry.IncCanNotConnect(jobName)
	}

}

//TODO Move this somewwhere else
func getAccessToken(upJob *job.UpJob) (string, error) {
	bearerToken := ""
	oauthServerResponse := upJob.OauthServerResponse

	refreshToken := false
	if oauthServerResponse.ExpiresAt != (time.Time{}) {
		if oauthServerResponse.ExpiresAt.Sub(time.Now()).Seconds() < 10 {
			refreshToken = true
		}

	}
	if oauthServerResponse.AccessToken == "" {
		refreshToken = true
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

		oauthServerResponse = oauthResponse
		oauthServerResponse.ExpiresAt = time.Now().Add(time.Second * time.Duration(oauthServerResponse.ExpiresIn))
	} else {
		bearerToken = oauthServerResponse.AccessToken
	}
	return bearerToken, nil
}

func matchResponse(pattern string, bytes []byte, reverseMatch bool) bool {
	regex, _ := regexp.Compile(pattern)
	doesMatch := regex.Match(bytes)
	if reverseMatch {
		return !doesMatch
	}
	return doesMatch
}

func doDnsChecks(registry metrics.Registry, dnsJob job.DnsJob) {

	jobName := dnsJob.Name
	actIps, err := net.LookupIP(dnsJob.FQDN)

	if err != nil {
		logrus.Warn(fmt.Sprintf("Failed Request for job '%s'. msg: '%s' ", jobName, err))
		return
	}

	actIpsS := make([]string, 0)
	for _, v := range actIps {
		actIpsS = append(actIpsS, v.String())
	}
	expIps := dnsJob.IPs

	intersecIps := GetIntersecting(expIps, actIpsS)
	if len(expIps) == 0 {
		registry.SetIpsRatio(jobName, 0)
	}
	registry.SetIpsRatio(jobName, float64(len(intersecIps)/len(expIps)))
}

func GetIntersecting(expIps []string, actIps []string) []string {
	intersecIps := make([]string, 0)
	for _, expIp := range expIps {
		for _, actIp := range actIps {
			if actIp == expIp {
				intersecIps = append(intersecIps, expIp)
			}
		}

	}
	return intersecIps
}

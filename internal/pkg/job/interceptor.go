package job

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (job *UpJob) init() {
	if job.Oauth.AuthUrl != "" {
		interceptor := OauthInterceptor{}
		job.RequestInterceptors = append(job.RequestInterceptors, interceptor)
	}
}

type OauthInterceptor struct {
}

func (interceptor OauthInterceptor) Intercept(req *http.Request, upJob *UpJob) *http.Request {
	//Perform authentication via oauth client credentials flow
	bearerToken, err := getAccessToken(upJob)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to receive Bearer Token for url: '%s' and auth_url: '%s'", upJob.URL, upJob.Oauth.AuthUrl))
		logrus.Error(err)
		return req
	}
	//Add Bearer Token in Authorization Header
	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	}
	return req

}

func getAccessToken(upJob *UpJob) (string, error) {
	bearerToken := ""
	//fetch last oauthResponse from context, to soo, if token has to be refreshed
	oauthServerResponse := &upJob.Context.OauthResponse
	refreshToken := false
	//if no Access Token provided yet
	if oauthServerResponse.AccessToken == "" {
		refreshToken = true
	}
	//check expiry time, if expires soon, do refresh
	if oauthServerResponse.ExpiresAt != (time.Time{}) {
		if oauthServerResponse.ExpiresAt.Sub(time.Now()).Seconds() < 10 {
			refreshToken = true
		}

	}

	//check if last response suggests to refresh
	if oauthServerResponse.Refresh {
		refreshToken = true
		oauthServerResponse.Refresh = false
	}

	//Request for a new Token, otherwise use the one stored in the context
	if refreshToken {
		values := url.Values{}
		//load params from Oauth
		for k, v := range upJob.Oauth.Params {
			values.Set(k, v)
		}
		//In case, a Refresh token was provided in the previous token request, send a refresh_token request
		if oauthServerResponse.RefreshToken != "" {
			values.Set("refresh_token", oauthServerResponse.RefreshToken)
			values.Set("request_type", "refresh_token")
		}

		authReq, _ := http.NewRequest("POST", upJob.Oauth.AuthUrl, strings.NewReader(values.Encode()))

		//Set headers
		for k, v := range upJob.Oauth.Headers {
			authReq.Header.Add(k, v)
		}

		client := &http.Client{}
		authResp, err := client.Do(authReq)
		if err != nil {
			logrus.Fatal(err)
			*oauthServerResponse = OauthResponse{}
			return "", err
		}
		if authResp.StatusCode != http.StatusOK {
			errorMsg := fmt.Sprintf("Failed authentication on auth_url: '%s'", upJob.Oauth.AuthUrl)
			logrus.Error(errorMsg)
			*oauthServerResponse = OauthResponse{}
			return "", fmt.Errorf(errorMsg)

		}
		bytes, _ := ioutil.ReadAll(authResp.Body)
		oauthResponse := OauthResponse{}

		err = yaml.Unmarshal(bytes, &oauthResponse)

		*oauthServerResponse = oauthResponse
		oauthServerResponse.ExpiresAt = time.Now().Add(time.Second * time.Duration(oauthServerResponse.ExpiresIn))
	}
	bearerToken = oauthServerResponse.AccessToken

	return bearerToken, nil
}

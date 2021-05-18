package ctl

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
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
func transport(job *job.UpJob) ChecksTripper {
	t := http.Transport{
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

	return ChecksTripper{Tripper: &t, Job: job}

}

// This type implements the http.RoundTripper interface
type ChecksTripper struct {
	Tripper http.RoundTripper
	Job     *job.UpJob
}

func (rt ChecksTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {

	for _, interceptor := range rt.Job.RequestInterceptors {
		req = interceptor.Intercept(req, rt.Job)
	}

	fmt.Printf("Sending request to %v\n", req.URL)

	// Send the request, get the response (or the error)
	res, e = rt.Tripper.RoundTrip(req)

	// Handle the result.
	if e != nil {
		fmt.Printf("Error: %v", e)
	} else {
		fmt.Printf("Received %v response\n", res.Status)
	}
	return
}

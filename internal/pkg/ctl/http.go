package ctl

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
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

//the pupose of a handwritten Transport is to have a clean uncached request in each iteration.
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

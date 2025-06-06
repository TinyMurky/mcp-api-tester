// Package netclient is the package that can use to post stuff
// Deprecate
package netclient

// import (
// 	"net"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// var (
// 	netClientInstance *NetClient
// 	once              sync.Once // For Client Singleton
// )

// const (
// 	// DialerTimeout determine how long will http.client try to connect target server
// 	DialerTimeout = 5 * time.Second

// 	// TSLHandshakeTimeout determine how long http.client will do TSL handshake and create secure online interaction
// 	TSLHandshakeTimeout = 5 * time.Second

// 	// ResponseHeaderTimeout determine how long http.client will wait server send back response header,
// 	// indicate that server atleast start running
// 	ResponseHeaderTimeout = 5 * time.Second

// 	// ClientTimeout determine how long http.client will wait for server to reponse
// 	ClientTimeout = 10 * time.Second

// 	// IdleConnTimeout is the maximum amount of time an idle
// 	// (keep-alive) connection will remain idle before closing
// 	// itself.
// 	IdleConnTimeout = 60 * time.Second
// )

// // NetClient wrap
// type NetClient struct {
// 	*http.Client
// 	url string
// }

// // GetClient return a singleton instance of NetClient
// func GetClient() *NetClient {
// 	once.Do(func() {
// 		netClientInstance = initNetClitent()
// 	})
// 	return netClientInstance
// }

// // initNetClitent create a new instance of NetClient
// func initNetClitent() *NetClient {

// 	dailer := &net.Dialer{
// 		Timeout: DialerTimeout,
// 	}

// 	transport := &http.Transport{
// 		DialContext:           dailer.DialContext,
// 		TLSHandshakeTimeout:   TSLHandshakeTimeout,
// 		ResponseHeaderTimeout: ResponseHeaderTimeout,
// 		IdleConnTimeout:       IdleConnTimeout,
// 	}

// 	httpClient := &http.Client{
// 		Transport: transport,
// 		Timeout:   ClientTimeout,
// 	}

// 	return &NetClient{Client: httpClient}
// }

// // SetURL set the URL for the NetClient
// func SetURL(url string) {
// 	nc := GetClient()
// 	nc.url = url
// }

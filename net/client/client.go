// Package netclient is the package that can use to post stuff
package netclient

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	netClientInstance *NetClient
	once              sync.Once // For Client Singleton
)

const (
	// DialerTimeout determine how long will http.client try to connect target server
	DialerTimeout = 2 * time.Second

	// TSLHandshakeTimeout determine how long http.client will do TSL handshake and create secure online interaction
	TSLHandshakeTimeout = 3 * time.Second

	// ResponseHeaderTimeout determine how long http.client will wait server send back response header,
	// indicate that server atleast start running
	ResponseHeaderTimeout = 5 * time.Second

	// ClientTimeout determine how long http.client will wait for server to reponse
	ClientTimeout = 10 * time.Second

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	IdleConnTimeout = 60 * time.Second
)

// NetClient wrap
type NetClient struct {
	*http.Client
}

func initNetClitent() *NetClient {

	dailer := &net.Dialer{}

	trasport := &http.Transport{
		DialContext: (&net.Dialer{}).DialContext,
	}
}

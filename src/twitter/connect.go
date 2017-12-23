package twitter

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	httpClient *http.Client
	conn       net.Conn
	reader     io.ReadCloser
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		authorizeClient()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			},
		}
	})
	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	authClient.SetAuthorizationHeader(req.Header, creds, "POST", req.URL, params)
	return httpClient.Do(req)
}

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, time.Second*5)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

// CloseConn closes the current connection to the twitter Stream
func CloseConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

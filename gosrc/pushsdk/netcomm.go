package pushsdk

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"rdpalert/embedded"
)

const (
	postJSONContentType = "application/json; charset=utf-8"
)

var (
	ErrHttpRequestFailed = errors.New("http request failed")
	customUserAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Patrick-RDPAlert-Assist/" + embedded.CurVersionStr
)

type customUserAgentRT struct {
	UserAgent string `json:"userAgent"`
}

func (cuart *customUserAgentRT) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", cuart.UserAgent)
	return http.DefaultTransport.RoundTrip(req)
}

func SendHttpPostJSON(url string, body []byte) ([]byte, int, error) {
	http.DefaultClient.Transport = &customUserAgentRT{UserAgent: customUserAgent}
	buf := bytes.NewBuffer(body)
	respD, err := http.Post(url, postJSONContentType, buf)
	if err != nil {
		return nil, -1, err
	}
	resp, err := io.ReadAll(respD.Body)
	defer func() { _ = respD.Body.Close() }()
	if err != nil {
		return nil, respD.StatusCode, err
	}
	if respD.StatusCode != http.StatusOK {
		return resp, respD.StatusCode, ErrHttpRequestFailed
	}
	return resp, respD.StatusCode, nil
}

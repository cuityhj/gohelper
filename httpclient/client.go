package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var gHttpClient *HttpClient

func init() {
	gHttpClient = &HttpClient{&http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			TLSHandshakeTimeout:   5 * time.Second,
			MaxIdleConns:          1,
			IdleConnTimeout:       5 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
				DualStack: true,
			}).DialContext,
		},
	}}
}

type HttpClient struct {
	*http.Client
}

func GetHttpClient() *HttpClient {
	return gHttpClient
}

type HttpContext struct {
	URL      string
	Request  interface{}
	Response interface{}
}

func (cli *HttpClient) Post(ctx *HttpContext) error {
	return cli.request(http.MethodPost, ctx.URL, ctx.Request, ctx.Response)
}

func (cli *HttpClient) Get(ctx *HttpContext) error {
	return cli.request(http.MethodGet, ctx.URL, ctx.Request, ctx.Response)
}

func (cli *HttpClient) Put(ctx *HttpContext) error {
	return cli.request(http.MethodPut, ctx.URL, ctx.Request, ctx.Response)
}

func (cli *HttpClient) Delete(ctx *HttpContext) error {
	return cli.request(http.MethodDelete, ctx.URL, ctx.Request, ctx.Response)
}

func (cli *HttpClient) request(httpMethod, url string, req, resp interface{}) error {
	var httpReqBody io.Reader
	if req != nil {
		reqBody, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("marshal request failed: %s", err.Error())
		}

		httpReqBody = bytes.NewBuffer(reqBody)
	}

	httpReq, err := http.NewRequest(httpMethod, url, httpReqBody)
	if err != nil {
		return fmt.Errorf("new http request failed: %s", err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := cli.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send http request failed: %s", err.Error())
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read http response body failed: %s", err.Error())
	}

	if len(body) > 0 && resp != nil {
		if err := json.Unmarshal(body, resp); err != nil {
			return fmt.Errorf("unmarshal http response failed: %s", err.Error())
		}
	}

	return nil
}

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
	URL           string
	Request       interface{}
	Response      interface{}
	ResponseError error
}

func (cli *HttpClient) Post(ctx *HttpContext) error {
	return cli.request(http.MethodPost, ctx)
}

func (cli *HttpClient) Get(ctx *HttpContext) error {
	return cli.request(http.MethodGet, ctx)
}

func (cli *HttpClient) Put(ctx *HttpContext) error {
	return cli.request(http.MethodPut, ctx)
}

func (cli *HttpClient) Delete(ctx *HttpContext) error {
	return cli.request(http.MethodDelete, ctx)
}

func (cli *HttpClient) request(httpMethod string, ctx *HttpContext) error {
	var httpReqBody io.Reader
	if ctx.Request != nil {
		reqBody, err := json.Marshal(ctx.Request)
		if err != nil {
			return fmt.Errorf("marshal request failed: %s", err.Error())
		}

		httpReqBody = bytes.NewBuffer(reqBody)
	}

	httpReq, err := http.NewRequest(httpMethod, ctx.URL, httpReqBody)
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

	if checkRespStatusCodeValid(httpResp.StatusCode) == false {
		if ctx.ResponseError != nil {
			if err := json.Unmarshal(body, ctx.ResponseError); err != nil {
				return fmt.Errorf("unmarshal http response error failed: %s", err.Error())
			} else {
				return fmt.Errorf("handle http request failed: %d %s",
					httpResp.StatusCode, ctx.ResponseError.Error())
			}
		} else {
			return fmt.Errorf("handle http request failed with status code %d",
				httpResp.StatusCode)
		}
	}

	if len(body) > 0 && ctx.Response != nil {
		if err := json.Unmarshal(body, ctx.Response); err != nil {
			return fmt.Errorf("unmarshal http response failed: %s", err.Error())
		}
	}

	return nil
}

func checkRespStatusCodeValid(code int) bool {
	switch code {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return true
	default:
		return false
	}
}

package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RestClient struct {
	Addr       string
	restClient *http.Client
}

func NewRestClient() *RestClient {
	return &RestClient{
		restClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false,
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
			},
			Timeout: time.Second * 300,
		},
	}
}

func (self *RestClient) SetAddr(addr string) *RestClient {
	self.Addr = addr
	return self
}

func (self *RestClient) SetRestClient(restClient *http.Client) *RestClient {
	self.restClient = restClient
	return self
}

func (self *RestClient) SendPostRequest(addr string, data []byte) ([]byte, error) {
	resp, err := self.restClient.Post(addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rest response body error:%s", err)
	}
	return body, nil
}

func (self *RestClient) SendGetRequest(addr string, params map[string]string) ([]byte, error) {
	encodeParam := "?"
	for name, value := range params {
		encodeParam += name + "=" + value + "&"
	}
	encodeParam = encodeParam[:len(encodeParam)-1]
	resp, err := self.restClient.Get(addr + encodeParam)
	if err != nil {
		return nil, fmt.Errorf("http get params: %s, error: %s", encodeParam, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rest response body error:%s", err)
	}
	return body, nil
}

package spvchannels

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ClientConfig struct {
	Insecure bool // equivalent curl -k
	BaseURL  string
	Version  string
	User     string
	Passwd   string
	Token    string
}

type Client struct {
	cfg        ClientConfig
	HTTPClient HTTPClient
}

func NewClient(c ClientConfig) *Client {
	httpClient := http.Client{
		Timeout: time.Minute,
	}

	if c.Insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		cfg:        c,
		HTTPClient: &httpClient,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	if c.cfg.Token == "" {
		req.SetBasicAuth(c.cfg.User, c.cfg.Passwd)
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.Token))
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := successResponse{
		Code: res.StatusCode,
		Data: v,
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&fullResponse.Data); err != nil {
			return err
		}
	}

	return nil
}

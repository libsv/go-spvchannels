package spvchannels

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	ws "github.com/gorilla/websocket"
)

// ClientConfig hold configuration for rest api connection
type ClientConfig struct {
	Insecure bool // equivalent curl -k
	BaseURL  string
	Version  string
	User     string
	Passwd   string
	Token    string
}

// Client hold rest api configuration and http connection
type Client struct {
	cfg        ClientConfig
	HTTPClient HTTPClient
}

// NewClient create a new rest Client
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

// errorResponse hold structure of error rest call
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// successResponse hold structure of success rest call
type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// sendRequest send the http request and retreive the response
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Websocket client is listening to the stream of notifications, which notifies new messages. It does not receive the
// content of the message itself (even though, the notification itself is a text)
// User has to write an separate engine to pull the new (unread) message content when it receive a notification.
// This can be easily done through the existing endpoint Messages in the provided rest api

// ProcessMessage is a callback to process websocket messages
//    t   : message type
//    msg : message content
//    err : message error
type ProcessMessage = func(t int, msg []byte, err error) error

// WSConfig hold the configuration for a websocket connection
type WSConfig struct {
	Insecure  bool // skip ssl certification check
	BaseURL   string
	Version   string
	ChannelId string
	Token     string // This token should be different then the one used to write messages
}

// WSClient is the structure holding the
//    - websocket configuration
//    - websocket connection
//    - message processor callback
//    - optionally max number of received messages (maxNotified)
type WSClient struct {
	cfg         WSConfig
	ws          *ws.Conn
	procces     ProcessMessage
	nbNotified  uint64
	maxNotified uint64
}

// NewWSClient create a new websocket client. Providing the
//    - WSConfig          : base configuration of the weboscket connection
//    - ProcessMessage    : callback to process the received messages
//    - uint64 (optional) : Max number of messages to process. After having received enough, the websocket exit
func NewWSClient(c WSConfig, p ProcessMessage, m ...uint64) *WSClient {
	max := uint64(0)
	if len(m) > 0 {
		max = m[0]
	}

	return &WSClient{
		cfg:         c,
		ws:          nil,
		procces:     p,
		nbNotified:  0,
		maxNotified: max,
	}
}

// urlPath return the path part of the connection URL
func (c *WSClient) urlPath() string {
	return fmt.Sprintf("/api/%s/channel/%s/notify", c.cfg.Version, c.cfg.ChannelId)
}

// NbNotified return the number of processed messages
func (c *WSClient) NbNotified() uint64 {
	return c.nbNotified
}

// Run establishes the connection and start listening the notification stream
// process the notification if a callback is provided
func (c *WSClient) Run() error {
	u := url.URL{
		Scheme: "wss",
		Host:   "localhost:5010",
		Path:   c.urlPath(),
	}

	q := u.Query()
	q.Set("token", c.cfg.Token)
	u.RawQuery = q.Encode()

	d := ws.DefaultDialer
	if c.cfg.Insecure {
		d = &ws.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		}
	}

	conn, _, err := d.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	c.ws = conn

	for c.nbNotified < c.maxNotified {
		t, msg, err := c.ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("%v. Total processed %d messages", err, c.nbNotified)
		}

		c.nbNotified += 1
		if c.procces != nil {
			err := c.procces(t, msg, err)
			if err != nil {
				return fmt.Errorf("%v. Total processed %d messages", err, c.nbNotified)
			}
		}
	}

	return nil
}

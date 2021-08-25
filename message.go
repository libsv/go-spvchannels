package spvchannels

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type (

	// Message contains the structure of a channel message.
	Message struct {
		Sequence    int       `json:"sequence"`
		Received    time.Time `json:"received"`
		ContentType string    `json:"content_type"`
		Payload     string    `json:"payload"`
	}

	// Sequence contains information about reading.
	Sequence struct {
		Read bool `json:"read"`
	}
)

// HeadMessage calls head on service message endpoint.
func (c *Client) HeadMessage(ctx context.Context, channelId string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://"+c.cfg.BaseURL+fmt.Sprintf(message, c.cfg.Version, channelId), nil)
	if err != nil {
		return false, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (c *Client) GetMessage(ctx context.Context, channelId string) (*[]*Message, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+c.cfg.BaseURL+fmt.Sprintf(message, c.cfg.Version, channelId), nil)
	if err != nil {
		return nil, err
	}

	res := []*Message{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) WriteMessage(ctx context.Context, channelId string) (*Message, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://"+c.cfg.BaseURL+fmt.Sprintf(message, c.cfg.Version, channelId), nil)
	if err != nil {
		return nil, err
	}

	res := Message{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) WriteMessageSequence(ctx context.Context, channelId, sequenceId string, older bool) (*Sequence, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://"+c.cfg.BaseURL+fmt.Sprintf(messageSequence, c.cfg.Version, channelId, sequenceId), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("older", strconv.FormatBool(older))
	req.URL.RawQuery = q.Encode()

	res := Sequence{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteMessageSequence(ctx context.Context, channelId, sequenceId string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "https://"+c.cfg.BaseURL+fmt.Sprintf(message, c.cfg.Version, channelId), nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}

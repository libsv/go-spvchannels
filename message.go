package spvchannels

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

func (c *Client) getMessageBaseEndpoint() string {
	return fmt.Sprintf("https://%s/api/%s", c.cfg.BaseURL, c.cfg.Version)
}

// MessageHeadRequest head request object.
type MessageHeadRequest struct {
	ChannelId string `json:"channelid"`
}

// MessageWriteRequest request struct for creating a new message.
type MessageWriteRequest struct {
	ChannelId string `json:"channelid"`
	Message   string `json:"message"`
}

// MessageWriteReply is the object returned when creating a new message.
type MessageWriteReply struct {
	Sequence   int64  `json:"sequence"`
	Received   string `json:"received"`
	ContenType string `json:"content_type"`
	Payload    string `json:"payload"`
}

// MessagesRequest argument for creating a message request.
type MessagesRequest struct {
	ChannelId string `json:"channelid"`
	UnRead    bool   `json:"unread"`
}

// MessagesReply slice of MessageWriteReply
type MessagesReply []MessageWriteReply

// MessageMarkRequest is the struct used when marking messages as read.
type MessageMarkRequest struct {
	ChannelId string `json:"channelid"`
	Sequence  int64  `json:"sequence"`
	Older     bool   `json:"older"`
	Read      bool   `json:"read"`
}

// MessageMarkReply response struct when marking message as read.
type MessageMarkReply struct{}

// MessageDeleteRequest arg used for deleting a message.
type MessageDeleteRequest struct {
	ChannelId string `json:"channelid"`
	Sequence  int64  `json:"sequence"`
}

// MessageDeleteReply response for message delete.
type MessageDeleteReply struct{}

// MessageHead performs a Head request for message.
func (c *Client) MessageHead(ctx context.Context, r MessageHeadRequest) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelId), nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}

// MessageWrite performs a Post request.
func (c *Client) MessageWrite(ctx context.Context, r MessageWriteRequest) (*MessageWriteReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelId), bytes.NewBuffer([]byte(r.Message)))
	if err != nil {
		return nil, err
	}

	res := MessageWriteReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Messages performs a Get request.
func (c *Client) Messages(ctx context.Context, r MessagesRequest) (*MessagesReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("unread", fmt.Sprintf("%t", r.UnRead))
	req.URL.RawQuery = q.Encode()

	res := MessagesReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// MessageMark perorms a Post request for message mark.
func (c *Client) MessageMark(ctx context.Context, r MessageMarkRequest) (*MessageMarkReply, error) {
	payloadStr := fmt.Sprintf("{\"read\":%t}", r.Read)
	channelURL := fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelId)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%v", channelURL, r.Sequence), bytes.NewBuffer([]byte(payloadStr)))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("older", fmt.Sprintf("%t", r.Older))
	req.URL.RawQuery = q.Encode()

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &MessageMarkReply{}, nil
}

// MessageDelete performs a Delete request.
func (c *Client) MessageDelete(ctx context.Context, r MessageDeleteRequest) (*MessageDeleteReply, error) {
	channelURL := fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelId)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%v", channelURL, r.Sequence), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &MessageDeleteReply{}, nil
}

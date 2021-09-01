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

type MessageHeadRequest struct {
	ChannelId string `json:"channelid"`
}

type MessageHeadReply struct {
}

type MessageWriteRequest struct {
	ChannelId string `json:"channelid"`
	Message   string `json:"message"`
}

type MessageWriteReply struct {
	Sequence   int64  `json:"sequence"`
	Received   string `json:"received"`
	ContenType string `json:"content_type"`
	Payload    string `json:"payload"`
}

type MessagesRequest struct {
	ChannelId string `json:"channelid"`
	UnRead    bool   `json:"unread"`
}

type MessagesReply []MessageWriteReply

type MessageMarkRequest struct {
	ChannelId string `json:"channelid"`
	Sequence  int64  `json:"sequence"`
	Older     bool   `json:"older"`
	Read      bool   `json:"read"`
}

type MessageMarkReply struct {
}

type MessageDeleteRequest struct {
	ChannelId string `json:"channelid"`
	Sequence  int64  `json:"sequence"`
}

type MessageDeleteReply struct {
}

func (c *Client) MessageHead(ctx context.Context, r MessageHeadRequest) (*MessageHeadReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &MessageHeadReply{}, nil
}

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

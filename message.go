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

// MessageHeadRequest hold data for HEAD message request
type MessageHeadRequest struct {
	ChannelID string `json:"channelid"`
}

// MessageHeadReply define HEAD message reply
type MessageHeadReply struct {
}

// MessageWriteRequest hold data for write message request
type MessageWriteRequest struct {
	ChannelID string `json:"channelid"`
	Message   string `json:"message"`
}

// MessageWriteReply hold data for write message reply
type MessageWriteReply struct {
	Sequence    int64  `json:"sequence"`
	Received    string `json:"received"`
	ContentType string `json:"content_type"`
	Payload     string `json:"payload"`
}

// MessagesRequest hold data for get messages request
type MessagesRequest struct {
	ChannelID string `json:"channelid"`
	UnRead    bool   `json:"unread"`
}

// MessagesReply hold data for get messages reply
type MessagesReply []MessageWriteReply

// MessageMarkRequest hold data for mark message request
type MessageMarkRequest struct {
	ChannelID string `json:"channelid"`
	Sequence  int64  `json:"sequence"`
	Older     bool   `json:"older"`
	Read      bool   `json:"read"`
}

// MessageMarkReply define mark message reply
type MessageMarkReply struct {
}

// MessageDeleteRequest hold data for delete message request
type MessageDeleteRequest struct {
	ChannelID string `json:"channelid"`
	Sequence  int64  `json:"sequence"`
}

// MessageDeleteReply define delete message reply
type MessageDeleteReply struct {
}

// MessageHead send HEAD message request
func (c *Client) MessageHead(ctx context.Context, r MessageHeadRequest) (*MessageHeadReply, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodHead, fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelID),
		nil,
	)

	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &MessageHeadReply{}, nil
}

// MessageWrite write a message
func (c *Client) MessageWrite(ctx context.Context, r MessageWriteRequest) (*MessageWriteReply, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelID), bytes.NewBuffer([]byte(r.Message)),
	)

	if err != nil {
		return nil, err
	}

	res := MessageWriteReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Messages get messages list
func (c *Client) Messages(ctx context.Context, r MessagesRequest) (*MessagesReply, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelID),
		nil,
	)

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

// MessageMark mark a message
func (c *Client) MessageMark(ctx context.Context, r MessageMarkRequest) (*MessageMarkReply, error) {
	payloadStr := fmt.Sprintf("{\"read\":%t}", r.Read)
	channelURL := fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelID)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%v", channelURL, r.Sequence), bytes.NewBuffer([]byte(payloadStr)),
	)

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

// MessageDelete delete a message
func (c *Client) MessageDelete(ctx context.Context, r MessageDeleteRequest) (*MessageDeleteReply, error) {
	channelURL := fmt.Sprintf("%s/channel/%s", c.getMessageBaseEndpoint(), r.ChannelID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%v", channelURL, r.Sequence), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &MessageDeleteReply{}, nil
}

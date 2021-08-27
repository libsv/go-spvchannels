package spvchannels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) getChanelBaseEndpoint() string {
	return fmt.Sprintf("https://%s/api/%s/account", c.cfg.BaseURL, c.cfg.Version)
}

func (c *Client) getTokenBaseEndpoint(accountid, channelid string) string {
	return fmt.Sprintf("%s/%s/channel/%s/api-token", c.getChanelBaseEndpoint(), accountid, channelid)
}

type ChannelsRequest struct {
	AccountId string `json:"accountid"`
}

type ChannelsReply struct {
	Channels []struct {
		Id          string `json:"id"`
		Href        string `json:"href"`
		PublicRead  bool   `json:"public_read"`
		PublicWrite bool   `json:"public_write"`
		Sequenced   bool   `json:"sequenced"`
		Locked      bool   `json:"locked"`
		Head        int    `json:"head"`
		Retention   struct {
			MinAgeDays int  `json:"min_age_days"`
			MaxAgeDays int  `json:"max_age_days"`
			AutoPrune  bool `json:"auto_prune"`
		} `json:"retention"`
		AccessTokens []struct {
			Id          string `json:"id"`
			Token       string `json:"token"`
			Description string `json:"description"`
			CanRead     bool   `json:"can_read"`
			CanWrite    bool   `json:"can_write"`
		} `json:"access_tokens"`
	} `json:"channels"`
}

type ChannelRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type ChannelReply struct {
	Id          string `json:"id"`
	Href        string `json:"href"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Locked      bool   `json:"locked"`
	Head        int    `json:"head"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
	AccessTokens []struct {
		Id          string `json:"id"`
		Token       string `json:"token"`
		Description string `json:"description"`
		CanRead     bool   `json:"can_read"`
		CanWrite    bool   `json:"can_write"`
	} `json:"access_tokens"`
}

type ChannelUpdateRequest struct {
	AccountId   string `json:"accountid"`
	ChannelId   string `json:"channelid"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Locked      bool   `json:"locked"`
}

type ChannelUpdateReply struct {
	PublicRead  bool `json:"public_read"`
	PublicWrite bool `json:"public_write"`
	Locked      bool `json:"locked"`
}

type ChannelDeleteRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type ChannelDeleteReply struct {
}

type ChannelCreateRequest struct {
	AccountId   string `json:"accountid"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
}

type ChannelCreateReply struct {
	Id          string `json:"id"`
	Href        string `json:"href"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Locked      bool   `json:"locked"`
	Head        int    `json:"head"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
	AccessTokens []struct {
		Id          string `json:"id"`
		Token       string `json:"token"`
		Description string `json:"description"`
		CanRead     bool   `json:"can_read"`
		CanWrite    bool   `json:"can_write"`
	} `json:"access_tokens"`
}

type TokenRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
	TokenId   string `json:"tokenid"`
}

type TokenReply struct {
	Id          string `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

type TokenDeleteRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
	TokenId   string `json:"tokenid"`
}

type TokenDeleteReply struct {
}

type TokensRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type TokensReply []TokenReply

type TokenCreateRequest struct {
	AccountId   string `json:"accountid"`
	ChannelId   string `json:"channelid"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

type TokenCreateReply struct {
	Id          string `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

func (c *Client) Channels(ctx context.Context, r ChannelsRequest) (*ChannelsReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/channel/list", c.getChanelBaseEndpoint(), r.AccountId), nil)
	if err != nil {
		return nil, err
	}

	res := ChannelsReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) Channel(ctx context.Context, r ChannelRequest) (*ChannelReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	res := ChannelReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) ChannelUpdate(ctx context.Context, r ChannelUpdateRequest) (*ChannelUpdateReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountId, r.ChannelId), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	res := ChannelUpdateReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) ChannelDelete(ctx context.Context, r ChannelDeleteRequest) (*ChannelDeleteReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &ChannelDeleteReply{}, nil
}

func (c *Client) ChannelCreate(ctx context.Context, r ChannelCreateRequest) (*ChannelCreateReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s/channel", c.getChanelBaseEndpoint(), r.AccountId), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	res := ChannelCreateReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) Token(ctx context.Context, r TokenRequest) (*TokenReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), r.TokenId), nil)
	if err != nil {
		return nil, err
	}

	res := TokenReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) TokenDelete(ctx context.Context, r TokenDeleteRequest) (*TokenDeleteReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), r.TokenId), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &TokenDeleteReply{}, nil
}

func (c *Client) Tokens(ctx context.Context, r TokensRequest) (*TokensReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	res := TokensReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) TokenCreate(ctx context.Context, r TokenCreateRequest) (*TokenCreateReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	res := TokenCreateReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

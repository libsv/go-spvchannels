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

type GetChannelsRequest struct {
	AccountId string `json:"accountid"`
}

type GetChannelsReply struct {
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

type GetChannelRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type GetChannelReply struct {
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

type UpdateChannelRequest struct {
	AccountId   string `json:"accountid"`
	ChannelId   string `json:"channelid"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Locked      bool   `json:"locked"`
}

type UpdateChannelReply struct {
	PublicRead  bool `json:"public_read"`
	PublicWrite bool `json:"public_write"`
	Locked      bool `json:"locked"`
}

type DeleteChannelRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type DeleteChannelReply struct {
}

type CreateChannelRequest struct {
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

type CreateChannelReply struct {
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

type GetTokenRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
	TokenId   string `json:"tokenid"`
}

type GetTokenReply struct {
	Id          string `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

type DeleteTokenRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
	TokenId   string `json:"tokenid"`
}

type DeleteTokenReply struct {
}

type GetTokensRequest struct {
	AccountId string `json:"accountid"`
	ChannelId string `json:"channelid"`
}

type GetTokensReply []GetTokenReply

type CreateTokenRequest struct {
	AccountId   string `json:"accountid"`
	ChannelId   string `json:"channelid"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

type CreateTokenReply struct {
	Id          string `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

func (c *Client) GetChannels(ctx context.Context, r GetChannelsRequest) (*GetChannelsReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/channel/list", c.getChanelBaseEndpoint(), r.AccountId), nil)
	if err != nil {
		return nil, err
	}

	res := GetChannelsReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetChannel(ctx context.Context, r GetChannelRequest) (*GetChannelReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	res := GetChannelReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateChannel(ctx context.Context, r UpdateChannelRequest) (*UpdateChannelReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountId, r.ChannelId), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	res := UpdateChannelReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteChannel(ctx context.Context, r DeleteChannelRequest) (*DeleteChannelReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &DeleteChannelReply{}, nil
}

func (c *Client) CreateChannel(ctx context.Context, r CreateChannelRequest) (*CreateChannelReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s/channel", c.getChanelBaseEndpoint(), r.AccountId), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	res := CreateChannelReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetToken(ctx context.Context, r GetTokenRequest) (*GetTokenReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), r.TokenId), nil)
	if err != nil {
		return nil, err
	}

	res := GetTokenReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteToken(ctx context.Context, r DeleteTokenRequest) (*DeleteTokenReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), r.TokenId), nil)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &DeleteTokenReply{}, nil
}

func (c *Client) GetTokens(ctx context.Context, r GetTokensRequest) (*GetTokensReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), nil)
	if err != nil {
		return nil, err
	}

	res := GetTokensReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CreateToken(ctx context.Context, r CreateTokenRequest) (*CreateTokenReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getTokenBaseEndpoint(r.AccountId, r.ChannelId), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	res := CreateTokenReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

package spvchannels

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c = &Client{
	cfg: ClientConfig{
		Insecure: true,
		BaseURL:  "somedomain",
	},
	HTTPClient: nil,
}

func TestChannels(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Channels": {
			request: `{
				AccountId: "1",
			}`,
			reply: `{
				"channels": [
					{
						"id": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
						"href": "https://localhost:5010/api/v1/channel/H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
						"public_read": true,
						"public_write": true,
						"sequenced": true,
						"locked": false,
						"head": 0,
						"retention": {
							"min_age_days": 0,
							"max_age_days": 99999,
							"auto_prune": true
						},
						"access_tokens": [
							{
								"id": "1",
								"token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
								"description": "Owner",
								"can_read": true,
								"can_write": true
							}
						]
					}
				]
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelsRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.Channels(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelsReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestChannel(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Channel": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
			}`,
			reply: `{
				"id": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"href": "https://localhost:5010/api/v1/channel/H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"public_read": true,
				"public_write": true,
				"sequenced": true,
				"locked": false,
				"head": 0,
				"retention": {
					"min_age_days": 0,
					"max_age_days": 99999,
					"auto_prune": true
				},
				"access_tokens": [
					{
						"id": "1",
						"token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
						"description": "Owner",
						"can_read": true,
						"can_write": true
					}
				]
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.Channel(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestChannelUpdate(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock ChannelUpdate": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"PublicRead": true,
				"PublicWrite": true,
				"Locked": true
			  }`,
			reply: `{
				"PublicRead": true,
				"PublicWrite": true,
				"Locked": true
			  }`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelUpdateRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.ChannelUpdate(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelUpdateReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestChannelDelete(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock ChannelDelete": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
			  }`,
			reply: "",
			err:   nil,
			code:  http.StatusNoContent,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelDeleteRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.ChannelDelete(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelDeleteReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestChannelCreate(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock ChannelCreate": {
			request: `{
				"AccountId" "1",
				"PublicRead": true,
				"PublicWrite": true,
				"Sequenced": true,
				"Retention": {
				  "MinAgeDays": 0,
				  "MaxAgeDays": 99999,
				  "AutoPrune": true
				}
			  }`,
			reply: `{
				"Id": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"Href": "https://localhost:5010/api/v1/channel/H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"PublicRead": true,
				"PublicWrite": true,
				"Sequenced": true,
				"Locked": false,
				"Head": 0,
				"Retention": {
					"MinAgeDays": 0,
					"MaxAgeDays": 99999,
					"AutoPrune": true
				},
				"AccessTokens": [
					{
						"Id": "1",
						"Token": "OEdvoTD3ozLxDfXrko2J3RKNHI7LrGW-sxyYF1aoLUNJI2mcFH9CMQXv3oRPbkcgx0EM3nEhYT61F6T72sPXEA",
						"Description": "Owner",
						"CanRead": true,
						"CanWrite": true
					}
				]
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelCreateRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.ChannelCreate(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelCreateReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestToken(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Token": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				TokenId : "1",
			  }`,
			reply: `{
				"Id": "1",
				"Token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
				"Description": "Owner",
				"CanRead": true,
				"CanWrite": true
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokenRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.Token(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokenReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestTokenDelete(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock TokenDelete": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				TokenId : "1",
			  }`,
			reply: "",
			err:   nil,
			code:  http.StatusNoContent,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokenDeleteRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.TokenDelete(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokenDeleteReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestTokens(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Tokens": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
			  }`,
			reply: `[
				{
					"Id": "1",
					"Token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
					"Description": "Owner",
					"CanRead": true,
					"CanWrite": true
				}
			]`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokensRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.Tokens(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokensReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestTokenCreate(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock TokenCreate": {
			request: `{
				AccountId: "1",
				ChannelId: "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"Description": "Owner",
				"CanRead": true,
				"CanWrite": true
			  }`,
			reply: `{
					"Id": "1",
					"Token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
					"Description": "Owner",
					"CanRead": true,
					"CanWrite": true
				}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokenCreateRequest
			json.Unmarshal([]byte(test.request), &req)
			resp, err := c.TokenCreate(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokenCreateReply
			json.Unmarshal([]byte(test.reply), &expectedResp)
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

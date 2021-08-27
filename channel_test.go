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

func TestGetChannels(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock GetChannels": {
			request: `{
				"accountid": "1"
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

			var req GetChannelsRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.GetChannels(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp GetChannelsReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestGetChannel(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock GetChannel": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
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

			var req GetChannelRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.GetChannel(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp GetChannelReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUpdateChannel(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock UpdateChannel": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
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

			var req UpdateChannelRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.UpdateChannel(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp UpdateChannelReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestDeleteChannel(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock DeleteChannel": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
			  }`,
			reply: "{}",
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

			var req DeleteChannelRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.DeleteChannel(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp DeleteChannelReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestCreateChannel(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock CreateChannel": {
			request: `{
				"accountid": "1",
				"public_read": true,
				"public_write": true,
				"sequenced": true,
				"retention": {
				  "min_age_days": 0,
				  "max_age_days": 99999,
				  "auto_prune": true
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

			var req CreateChannelRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.CreateChannel(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp CreateChannelReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestGetToken(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock GetToken": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"tokenid": "1"
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

			var req GetTokenRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.GetToken(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp GetTokenReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestDeleteToken(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock DeleteToken": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"tokenid": "1"
			  }`,
			reply: "{}",
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

			var req DeleteTokenRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.DeleteToken(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp DeleteTokenReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestGetTokens(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock GetTokens": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
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

			var req GetTokensRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.GetTokens(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp GetTokensReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestCreateToken(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock CreateToken": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"description": "Owner",
				"can_read": true,
				"can_write": true
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

			var req CreateTokenRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.CreateToken(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp CreateTokenReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

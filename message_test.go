package spvchannels

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockDo func(req *http.Request) (*http.Response, error)

// MockClient mocks the http client.
type MockClient struct {
	MockDo
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestHeadMessage(t *testing.T) {
	tests := map[string]struct {
		channelId string
		exp       bool
		mockDo    func(*http.Request) (*http.Response, error)
		err       error
	}{
		"should return 200 for valid head": {
			channelId: "abc",
			mockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(nil)),
				}, nil
			},
			exp: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := &Client{
				cfg: ClientConfig{
					Insecure: true,
					BaseURL:  "somedomain",
				},
				HTTPClient: &MockClient{
					MockDo: test.mockDo,
				},
			}
			result, err := c.HeadMessage(context.Background(), test.channelId)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.exp, result)
		})
	}
}

func TestGetMessage(t *testing.T) {
	tests := map[string]struct {
		channelId string
		message   *[]*Message
		mockDo    func(*http.Request) (*http.Response, error)
		err       error
	}{
		"should return message array with single message and correct id": {
			channelId: "abc-234",
			mockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("[{\"sequence\":2,\"received\":\"2021-08-24T08:49:46.210Z\",\"content_type\":\"text/plain\",\"payload\":\"this is a test\"}]"))),
				}, nil
			},
			message: &[]*Message{{
				Sequence:    2,
				Received:    time.Date(2021, time.August, 24, 8, 49, 46, 210000000, time.UTC),
				ContentType: "text/plain",
				Payload:     "this is a test",
			}},
		},
		"should return error when not found": {
			channelId: "abc-234",
			message:   &[]*Message{},
			mockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(bytes.NewReader(nil)),
				}, nil
			},
			err: errors.New("unknown error, status code: 404"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := &Client{
				cfg: ClientConfig{
					Insecure: true,
					BaseURL:  "somedomain",
				},
				HTTPClient: &MockClient{
					MockDo: test.mockDo,
				},
			}
			res, err := c.GetMessage(context.Background(), test.channelId)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.message, res)
		})
	}
}

func TestWriteMessage(t *testing.T) {
	tests := map[string]struct {
		channelId string
		mockDo    func(*http.Request) (*http.Response, error)
		message   *Message
		err       error
	}{
		"should return a new message": {
			channelId: "abc-234",
			mockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 201,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{\"sequence\":2,\"received\":\"2021-08-24T08:49:46.210Z\",\"content_type\":\"application/json\",\"payload\":\"this is a test\"}"))),
				}, nil
			},
			message: &Message{
				Sequence:    2,
				Received:    time.Date(2021, time.August, 24, 8, 49, 46, 210000000, time.UTC),
				ContentType: "application/json",
				Payload:     "this is a test",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := &Client{
				cfg: ClientConfig{
					Insecure: true,
					BaseURL:  "somedomain",
				},
				HTTPClient: &MockClient{
					MockDo: test.mockDo,
				},
			}
			res, err := c.WriteMessage(context.Background(), test.channelId)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.message, res)
		})
	}
}

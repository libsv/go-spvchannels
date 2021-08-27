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

func TestMessageHead(t *testing.T) {
	tests := map[string]struct {
		channelId string
		exp       bool
		response  string
		code      int
		err       error
	}{
		"should return 200 for valid head": {
			channelId: "abc",
			response:  "",
			code:      http.StatusOK,
			exp:       true,
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
					MockDo: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: test.code,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(test.response))),
						}, nil
					},
				},
			}
			result, err := c.MessageHead(context.Background(), test.channelId)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.exp, result)
		})
	}
}

func TestMessageGet(t *testing.T) {
	tests := map[string]struct {
		channelId string
		message   *[]*Message
		response  string
		code      int
		err       error
	}{
		"should return message array with single message and correct id": {
			channelId: "abc-234",
			response:  `[{"sequence":2,"received":"2021-08-24T08:49:46.210Z","content_type":"text/plain","payload":"this is a test"}]`,
			code:      http.StatusOK,
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
			code:      http.StatusNotFound,
			err:       errors.New("unknown error, status code: 404"),
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
					MockDo: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: test.code,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(test.response))),
						}, nil
					},
				},
			}
			res, err := c.MessageGet(context.Background(), test.channelId)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.message, res)
		})
	}
}

func TestMessage(t *testing.T) {
	tests := map[string]struct {
		channelId string
		response  string
		code      int
		message   *Message
		err       error
	}{
		"should return a new message": {
			channelId: "abc-234",
			response:  `{"sequence":2,"received":"2021-08-24T08:49:46.210Z","content_type":"application/json","payload":"this is a test"}`,
			code:      http.StatusCreated,
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
					MockDo: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: test.code,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(test.response))),
						}, nil
					},
				},
			}
			res, err := c.Message(context.Background(), test.channelId)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.message, res)
		})
	}
}

func TestMessageSequence(t *testing.T) {
	tests := map[string]struct {
		channelId string
		sequence  int64
		older     bool
		response  string
		code      int
		result    *Sequence
		err       error
	}{
		"should return a read result": {
			channelId: "abc-234",
			sequence:  2,
			code:      http.StatusOK,
			response:  `{"read": true}`,
			result: &Sequence{
				Read: true,
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
					MockDo: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: test.code,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(test.response))),
						}, nil
					},
				},
			}
			res, err := c.MessageSequence(context.Background(), test.channelId, test.sequence, test.older)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Equal(t, test.result, res)
		})
	}
}

func TestMessageSequenceDelete(t *testing.T) {
	tests := map[string]struct {
		err       error
		channelId string
		sequence  int64
		response  string
		code      int
	}{
		"should return OK when delete successful": {
			channelId: "abc",
			sequence:  2,
			response:  "",
			code:      http.StatusOK,
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
					MockDo: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: test.code,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(test.response))),
						}, nil
					},
				},
			}
			err := c.MessageSequenceDelete(context.Background(), test.channelId, test.sequence)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
			assert.Nil(t, err)
		})
	}
}

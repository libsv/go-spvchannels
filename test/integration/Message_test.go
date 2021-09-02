// +build integration

package integration

import (
	"context"
	"testing"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

func getClientConfigWithToken() (spv.ClientConfig, string) {
	cfg := getClientConfig()
	replyCreateChannel, _ := createChannel(cfg)
	cfg.Token = replyCreateChannel.AccessTokens[0].Token
	return cfg, replyCreateChannel.Id
}

func writeMessage(cfg spv.ClientConfig, channelid string) (*spv.MessageWriteReply, error) {
	client := spv.NewClient(cfg)

	r := spv.MessageWriteRequest{
		ChannelId: channelid,
		Message:   "Hello, this is a message",
	}

	reply, err := client.MessageWrite(context.Background(), r)
	return reply, err
}

func TestMessageIntegration(t *testing.T) {

	t.Run("TestMessageHead", func(t *testing.T) {
		cfg, channelid := getClientConfigWithToken()
		client := spv.NewClient(cfg)

		r := spv.MessageHeadRequest{
			ChannelId: channelid,
		}

		_, err := client.MessageHead(context.Background(), r)
		assert.NoError(t, err)
	})

	t.Run("TestMessageWrite", func(t *testing.T) {
		cfg, channelid := getClientConfigWithToken()

		reply, err := writeMessage(cfg, channelid)
		assert.NoError(t, err)
		assert.True(t, reply.Sequence > 0)
		assert.NotEmpty(t, reply.Payload)
	})

	t.Run("TestMessageWrite", func(t *testing.T) {
		cfg, channelid := getClientConfigWithToken()

		reply, err := writeMessage(cfg, channelid)
		assert.NoError(t, err)
		assert.True(t, reply.Sequence > 0)
		assert.NotEmpty(t, reply.Payload)
	})

	t.Run("TestMessages", func(t *testing.T) {
		cfg, channelid := getClientConfigWithToken()
		_, err := writeMessage(cfg, channelid)
		assert.NoError(t, err)

		client := spv.NewClient(cfg)

		r := spv.MessagesRequest{
			ChannelId: channelid,
			UnRead:    false,
		}

		reply, err := client.Messages(context.Background(), r)
		assert.NoError(t, err)
		assert.True(t, len(*reply) > 0)
		assert.True(t, (*reply)[0].Sequence > 0)
		assert.NotEmpty(t, (*reply)[0].Payload)
	})

	t.Run("TestMessageMark", func(t *testing.T) {
		cfg, channelid := getClientConfigWithToken()
		replyWriteMessage, _ := writeMessage(cfg, channelid)

		client := spv.NewClient(cfg)

		r := spv.MessageMarkRequest{
			ChannelId: channelid,
			Sequence:  replyWriteMessage.Sequence,
			Older:     false,
			Read:      false,
		}

		_, err := client.MessageMark(context.Background(), r)
		assert.NoError(t, err)
	})

	t.Run("TestMessageDelete", func(t *testing.T) {
		cfg, channelid := getClientConfigWithToken()
		replyWriteMessage, _ := writeMessage(cfg, channelid)

		client := spv.NewClient(cfg)

		r := spv.MessageDeleteRequest{
			ChannelId: channelid,
			Sequence:  replyWriteMessage.Sequence,
		}

		_, err := client.MessageDelete(context.Background(), r)
		assert.NoError(t, err)
	})
}
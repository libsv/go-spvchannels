// +build integration

package integration

import (
	"context"

	"testing"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

func createChannel(cfg spv.ClientConfig) (*spv.ChannelCreateReply, error) {
	client := spv.NewClient(cfg)

	r := spv.ChannelCreateRequest{
		AccountId:   accountid,
		PublicRead:  true,
		PublicWrite: true,
		Sequenced:   true,
		Retention: struct {
			MinAgeDays int  "json:\"min_age_days\""
			MaxAgeDays int  "json:\"max_age_days\""
			AutoPrune  bool "json:\"auto_prune\""
		}{
			MinAgeDays: 0,
			MaxAgeDays: 99999,
			AutoPrune:  true,
		},
	}

	reply, err := client.ChannelCreate(context.Background(), r)
	return reply, err
}

func getChannel(cfg spv.ClientConfig, accountid string, channelid string) (*spv.ChannelReply, error) {
	client := spv.NewClient(cfg)

	r := spv.ChannelRequest{
		AccountId: accountid,
		ChannelId: channelid,
	}

	reply, err := client.Channel(context.Background(), r)
	return reply, err
}

func getChannels(cfg spv.ClientConfig, accountid string) (*spv.ChannelsReply, error) {
	client := spv.NewClient(cfg)

	r := spv.ChannelsRequest{
		AccountId: accountid,
	}

	reply, err := client.Channels(context.Background(), r)
	return reply, err
}

func TestChannelIntegration(t *testing.T) {

	t.Run("TestChannels", func(t *testing.T) {
		cfg := getClientConfig()
		_, err := createChannel(cfg)
		assert.NoError(t, err)

		reply, err := getChannels(cfg, accountid)
		assert.NoError(t, err)
		assert.NotEmpty(t, reply.Channels)
	})

	t.Run("TestChannel", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		reply, err := getChannel(cfg, accountid, replyCreateChannel.Id)
		assert.NoError(t, err)
		assert.Equal(t, reply.Id, replyCreateChannel.Id)
	})

	t.Run("TestChannelUpdate", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		client := spv.NewClient(cfg)
		r := spv.ChannelUpdateRequest{
			AccountId:   accountid,
			ChannelId:   replyCreateChannel.Id,
			PublicRead:  false,
			PublicWrite: false,
			Locked:      false,
		}

		reply, err := client.ChannelUpdate(context.Background(), r)
		assert.NoError(t, err)
		assert.Equal(t, reply.PublicRead, r.PublicRead)
		assert.Equal(t, reply.PublicWrite, r.PublicWrite)
		assert.Equal(t, reply.Locked, r.Locked)
	})

	t.Run("TestChannelDelete", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		replyGetChannelsBefore, _ := getChannels(cfg, accountid)

		client := spv.NewClient(cfg)
		r := spv.ChannelDeleteRequest{
			AccountId: accountid,
			ChannelId: replyCreateChannel.Id,
		}
		_, err := client.ChannelDelete(context.Background(), r)
		assert.NoError(t, err)
		replyGetChannelsAfter, _ := getChannels(cfg, accountid)

		assert.Equal(t, len(replyGetChannelsBefore.Channels), len(replyGetChannelsAfter.Channels)+1)
	})

	t.Run("TestChannelCreate", func(t *testing.T) {
		cfg := getClientConfig()
		reply, err := createChannel(cfg)
		assert.NotNil(t, reply)
		assert.NoError(t, err)
		assert.Equal(t, len(reply.AccessTokens), 1)
		assert.NotEmpty(t, reply.Id)
		assert.NotEmpty(t, reply.AccessTokens[0].Token)
	})
}

func TestChannelTokenIntegration(t *testing.T) {

	t.Run("TestToken", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		client := spv.NewClient(cfg)
		r := spv.TokenRequest{
			AccountId: accountid,
			ChannelId: replyCreateChannel.Id,
			TokenId:   replyCreateChannel.AccessTokens[0].Id,
		}

		reply, err := client.Token(context.Background(), r)
		assert.NoError(t, err)
		assert.Equal(t, reply.Id, replyCreateChannel.AccessTokens[0].Id)
		assert.Equal(t, reply.Token, replyCreateChannel.AccessTokens[0].Token)
		assert.Equal(t, reply.Description, replyCreateChannel.AccessTokens[0].Description)
		assert.Equal(t, reply.CanRead, replyCreateChannel.AccessTokens[0].CanRead)
		assert.Equal(t, reply.CanWrite, replyCreateChannel.AccessTokens[0].CanWrite)
	})

	t.Run("TestTokenDelete", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		client := spv.NewClient(cfg)
		r := spv.TokenDeleteRequest{
			AccountId: accountid,
			ChannelId: replyCreateChannel.Id,
			TokenId:   replyCreateChannel.AccessTokens[0].Id,
		}

		_, err := client.TokenDelete(context.Background(), r)
		assert.NoError(t, err)

		r2 := spv.TokensRequest{
			AccountId: accountid,
			ChannelId: replyCreateChannel.Id,
		}

		// Token list after deleting the only one is empty
		reply, _ := client.Tokens(context.Background(), r2)
		assert.Equal(t, len(*reply), 0)
	})

	t.Run("TestTokens", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		client := spv.NewClient(cfg)
		r := spv.TokensRequest{
			AccountId: accountid,
			ChannelId: replyCreateChannel.Id,
		}

		reply, err := client.Tokens(context.Background(), r)
		assert.NoError(t, err)
		assert.Equal(t, len(*reply), 1)
		assert.Equal(t, (*reply)[0].Id, replyCreateChannel.AccessTokens[0].Id)
		assert.Equal(t, (*reply)[0].Token, replyCreateChannel.AccessTokens[0].Token)
		assert.Equal(t, (*reply)[0].Description, replyCreateChannel.AccessTokens[0].Description)
		assert.Equal(t, (*reply)[0].CanRead, replyCreateChannel.AccessTokens[0].CanRead)
		assert.Equal(t, (*reply)[0].CanWrite, replyCreateChannel.AccessTokens[0].CanWrite)
	})

	t.Run("TestTokenCreate", func(t *testing.T) {
		cfg := getClientConfig()
		replyCreateChannel, _ := createChannel(cfg)

		client := spv.NewClient(cfg)
		r := spv.TokenCreateRequest{
			AccountId:   accountid,
			ChannelId:   replyCreateChannel.Id,
			Description: "TestTokenCreate",
			CanRead:     true,
			CanWrite:    true,
		}

		reply, err := client.TokenCreate(context.Background(), r)
		assert.NoError(t, err)
		assert.NotEmpty(t, reply.Id)
		assert.NotEmpty(t, reply.Token)
		assert.Equal(t, reply.Description, r.Description)
		assert.Equal(t, reply.CanRead, r.CanRead)
		assert.Equal(t, reply.CanWrite, r.CanWrite)
	})
}
